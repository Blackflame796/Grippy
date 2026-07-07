package auth_usecase

import (
	entity "Grippy/internal/domain"
	redis_repository "Grippy/internal/infrastructure/redis/repository"
	"Grippy/internal/repository"
	user_usecase "Grippy/internal/usecase/user"
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repo       *redis_repository.SessionRepository
	userRepo   *repository.UserRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type jwtClaims struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthUseCase(repo *redis_repository.SessionRepository, userRepo *repository.UserRepository, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthUseCase {
	return &AuthUseCase{
		repo:       repo,
		userRepo:   userRepo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (u *AuthUseCase) GenerateTokenPair(ctx context.Context, id uuid.UUID, username string) (*entity.TokenPair, error) {
	now := time.Now()

	claims := jwtClaims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(u.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessTokenObj.SignedString(u.jwtSecret)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	refreshTokenStr := hex.EncodeToString(b)

	session := &entity.RefreshSession{
		UserID:       id,
		RefreshToken: refreshTokenStr,
		CreatedAt:    now,
		ExpiresAt:    now.Add(u.refreshTTL),
	}

	if err := u.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (u *AuthUseCase) Refresh(ctx context.Context, oldRefreshToken string, userRole string) (*entity.TokenPair, error) {

	session, err := u.repo.GetSessionByToken(ctx, oldRefreshToken)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.IsExpired() {
		_ = u.repo.DeleteSession(ctx, oldRefreshToken)
		return nil, ErrInvalidToken
	}

	if err := u.repo.DeleteSession(ctx, oldRefreshToken); err != nil {
		return nil, err
	}

	return u.GenerateTokenPair(ctx, session.UserID, userRole)
}

func (u *AuthUseCase) Register(ctx context.Context, req SignUpRequest) (*user_usecase.UserResponse, error) {
	_, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	_, err = u.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, ErrUsernameAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userInfo := &entity.UserInfo{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	user, err := u.userRepo.CreateUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	return &user_usecase.UserResponse{
		ID:       user.ID,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (u *AuthUseCase) Login(ctx context.Context, req SignInRequest) (*entity.TokenPair, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return u.GenerateTokenPair(ctx, user.ID, user.Username)
}

func (u *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return ErrInvalidToken
	}

	return u.repo.DeleteSession(ctx, refreshToken)
}

func (u *AuthUseCase) ParseAccessToken(tokenStr string) (*entity.Claims, error) {
	var claims jwtClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return u.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return &entity.Claims{
		ID:        claims.ID,
		Username:  claims.Username,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

func (u *AuthUseCase) ValidateRefreshToken(ctx context.Context, tokenStr string) error {
	session, err := u.repo.GetSessionByToken(ctx, tokenStr)
	if err != nil {
		return ErrSessionNotFound
	}

	if session.IsExpired() {
		_ = u.repo.DeleteSession(ctx, tokenStr)
		return ErrInvalidToken
	}

	return nil
}
