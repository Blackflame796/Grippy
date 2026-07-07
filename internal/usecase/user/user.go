package user_usecase

import (
	entity "Grippy/internal/domain"
	"Grippy/internal/repository"
	s3_storage "Grippy/pkg/s3"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repository.UserRepository
	s3   *s3_storage.S3Client
}

func NewUserUseCase(r *repository.UserRepository, s3 *s3_storage.S3Client) *UserUseCase {
	return &UserUseCase{repo: r, s3: s3}
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

func (uc *UserUseCase) UploadAvatar(ctx context.Context, userID uuid.UUID, file io.Reader, originalFilename string, contentType string) (string, error) {
	ext := filepath.Ext(originalFilename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	avatarURL, err := uc.s3.Upload(ctx, "users/avatars", uniqueFilename, file, contentType)
	if err != nil {
		return "", fmt.Errorf("usecase upload failed: %w", err)
	}

	_, oldAvatarURL, err := uc.repo.UpdateAvatar(ctx, userID, avatarURL)
	if err != nil {
		_ = uc.s3.Delete(ctx, avatarURL)
		return "", fmt.Errorf("failed to update avatar in database: %w", err)
	}

	if oldAvatarURL != "" {
		go func(url string) {
			_ = uc.s3.Delete(context.Background(), url)
		}(oldAvatarURL)
	}

	return avatarURL, nil
}
