package db

import (
	"context"
	"database/sql"
	"pulse/src/api-core/internal/domain"
)

type postgresGroupRepo struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) domain.GroupRepository {
	return &postgresGroupRepo{
		db: db,
	}
}

func (p *postgresGroupRepo) AddMember(ctx context.Context, groupID string, userID string, role domain.GroupMemberRole) error {
	insertQuery := `
		INSERT INTO group_members (group_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (group_id, user_id) DO NOTHING
	`

	_, err := p.db.ExecContext(ctx, insertQuery, groupID, userID, role)

	return err
}

// Create implements [domain.GroupRepository].
func (p *postgresGroupRepo) Create(ctx context.Context, group *domain.Group) error {
	insertQuery := `
		INSERT INTO groups (id, name, type)
		VALUES ($1, $2, $3)
	`

	_, err := p.db.ExecContext(ctx, insertQuery, group.ID, group.Name, group.Type)

	return err
}

func (p *postgresGroupRepo) GetByID(ctx context.Context, id string) (*domain.Group, error) {
	query := `
		SELECT id, name, type, created_at FROM groups WHERE id = $1
	`

	var group domain.Group
	err := p.db.QueryRowContext(ctx, query, id).Scan(&group.ID, &group.Name, &group.Type, &group.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrGroupNotFound
		}
		return nil, err
	}

	return &group, nil
}
