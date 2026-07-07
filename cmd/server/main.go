package main

import (
	handlers "ToDoApp/internal/core/transport/http/handlers"
	middlewares "ToDoApp/internal/core/transport/http/middlewares"
	"ToDoApp/pkg/database"
	"ToDoApp/pkg/logger"
	"ToDoApp/pkg/repository"
	"net/http"
	"runtime"
)

func main() {
	logger.InitLogger("development")
	defer logger.Log.Sync()
	runtime.GOMAXPROCS(runtime.NumCPU())
	connectLink := "postgres://Blackflame:pavel08180919@localhost:5432/Grippy?sslmode=disable"
	dbPool := database.InitDB(connectLink)
	defer database.CloseDB()
	logger.Log.Info("Starting server...")
	appRouter := handlers.NewRouter()
	appRouter.Use(middlewares.ZapLogger)
	todoRepo := repository.NewToDoRepository(dbPool)
	todoHandler := handlers.NewToDoHandler(todoRepo)
	todoHandler.RegisterRoutes(appRouter)
	err := http.ListenAndServe(":8080", appRouter.Init())
	if err != nil {
		logger.Log.Fatalf("Error of starting server: %v", err)
	}
	logger.Log.Info("Server successfully started")
}
