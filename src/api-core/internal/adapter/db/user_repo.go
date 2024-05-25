package db

import (
	"context"
	"database/sql"
	"pulse/src/api-core/internal/domain"
)

type postgresUserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &postgresUserRepo{
		db: db,
	}
}

func (up *postgresUserRepo) Create(ctx context.Context, user *domain.User) error {
	createQuery := "INSERT INTO users (id, username, avatar_url, created_at) VALUES ($1, $2, $3, $4)"

	_, err := up.db.ExecContext(ctx, createQuery, user.ID, user.Username, user.AvatarURL, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (up *postgresUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	err := up.db.QueryRowContext(ctx, "SELECT id, username, avatar_url, created_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.AvatarURL, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (up *postgresUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := up.db.QueryRowContext(ctx, "SELECT id, username, avatar_url, created_at FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.AvatarURL, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
