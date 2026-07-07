package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	redis_repo "Grippy/internal/infrastructure/redis/repository"
	"Grippy/internal/repository"
	"Grippy/internal/transport/http/handlers"
	"Grippy/internal/transport/http/middlewares"

	"Grippy/internal/transport/http/router"
	auth_usecase "Grippy/internal/usecase/auth"
	user_usecase "Grippy/internal/usecase/user"

	"Grippy/pkg/database"
	"Grippy/pkg/logger"
	cache "Grippy/pkg/redis"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	httpServer  *http.Server
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
}

func Init(address string, port int) (*Server, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.InitLogger("development")

	redisClient, err := cache.NewRedisClient(fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	connectLink := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), dbPort, os.Getenv("POSTGRES_DB"))
	dbPool := database.InitDB(connectLink)

	logger.Log.Info("Database and Redis connected")

	mainRouter := router.NewMainRouter()
	mainRouter.Use(middlewares.LoggingMiddleware)
	mainRouter.Use(middlewares.RecoveryMiddleware)

	accessTTL := 15 * time.Minute
	refreshTTL := 30 * 24 * time.Hour
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		logger.Log.Fatalln("SECRET_KEY not found in environment")
	}
	userRepo := repository.NewUserRepository(dbPool)
	sessionRepo := redis_repo.NewSessionRepository(redisClient)

	authUC := auth_usecase.NewAuthUseCase(
		sessionRepo,
		userRepo,
		secretKey,
		accessTTL,
		refreshTTL,
	)
	authHandler := handlers.NewAuthHandler(authUC)
	authRouter := router.New("/auth", mainRouter)
	authHandler.RegisterRoutes(authRouter)

	authMiddleware := middlewares.NewAuthMiddleware(authUC)

	userUC := user_usecase.NewUserUseCase(userRepo)
	userRouter := router.New("/user", mainRouter)
	userRouter.Use(authMiddleware)
	userHandler := handlers.NewUserHandler(userUC)
	userHandler.RegisterRoutes(userRouter)

	todoRepo := repository.NewToDoRepository(dbPool)
	todoRouter := router.New("/todos", mainRouter)
	todoRouter.Use(authMiddleware)
	todoHandler := handlers.NewToDoHandler(todoRepo)
	todoHandler.RegisterRoutes(todoRouter)

	srv := &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", address, port),
			Handler:      mainRouter.ServeMux(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		dbPool:      dbPool,
		redisClient: redisClient,
	}

	return srv, nil
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	logger.Log.Info("Server listening on ", ln.Addr().String())

	go func() {
		if err := s.httpServer.Serve(ln); err != http.ErrServerClosed {
			logger.Log.Errorf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Log.Errorf("Shutdown error: %v", err)
	}

	if s.dbPool != nil {
		s.dbPool.Close()
		logger.Log.Info("PostgreSQL pool closed")
	}

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			logger.Log.Errorf("Error closing Redis: %v", err)
		} else {
			logger.Log.Info("Redis client closed")
		}
	}

	logger.Log.Info("Resources closed, server stopped")
	return nil
}
