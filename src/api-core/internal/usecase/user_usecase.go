package usecase

import (
	"context"
	"pulse/src/api-core/internal/domain"

	"github.com/google/uuid"
)

type UserUseCase interface {
	RegisterUser(ctx context.Context, username string, avatarURL string) (*domain.User, error)
}

type userUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (uuc *userUseCase) RegisterUser(ctx context.Context, username string, avataURL string) (*domain.User, error) {
	newUser := &domain.User{
		ID:        uuid.NewString(),
		AvatarURL: avataURL,
		Username:  username,
	}

	err := uuc.repo.Create(ctx, newUser)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}
