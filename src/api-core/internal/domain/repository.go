package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type GroupRepository interface {
	Create(ctx context.Context, group *Group) error
	AddMember(ctx context.Context, groupID string, userID string, role GroupMemberRole) error
	GetByID(ctx context.Context, id string) (*Group, error)
}
