package user_usecase

import (
	entity "Grippy/internal/domain"
	"Grippy/internal/repository"
	"context"
	// "github.com/google/uuid"
)

type UserUseCase struct {
	repo *repository.UserRepository
}

func NewUserUseCase(r *repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) Update(ctx context.Context, input UpdateUserRequest, claims *entity.Claims) (*UserResponse, error) {
	inputUserInfo := &entity.User{
		ID:       claims.ID,
		Email:    input.Email,
		Username: input.Username,
	}

	user, err := uc.repo.UpdateUserInfo(ctx, inputUserInfo)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
