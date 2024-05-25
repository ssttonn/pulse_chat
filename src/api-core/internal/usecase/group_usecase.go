package usecase

import (
	"context"
	"pulse/src/api-core/internal/domain"

	"github.com/google/uuid"
)

type GroupUseCase interface {
	AddNewMember(ctx context.Context, groupID string, userID string, role domain.GroupMemberRole) error
	CreateNewGroup(ctx context.Context, creatorID string, name string, groupType domain.GroupType) error
}

type groupUseCase struct {
	groupRepo domain.GroupRepository
	userRepo  domain.UserRepository
}

func NewGroupUseCase(groupRepo domain.GroupRepository, userRepo domain.UserRepository) GroupUseCase {
	return &groupUseCase{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// AddNewMember implements [GroupUseCase].
func (g *groupUseCase) AddNewMember(ctx context.Context, groupID string, userID string, role domain.GroupMemberRole) error {
	existingGroup, err := g.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}

	err = g.groupRepo.AddMember(ctx, existingGroup.ID, userID, role)
	return err
}

// CreateNewGroup implements [GroupUseCase].
func (g *groupUseCase) CreateNewGroup(ctx context.Context, creatorID string, name string, groupType domain.GroupType) error {
	user, err := g.userRepo.GetByID(ctx, creatorID)
	if err != nil {
		return err
	}

	newGroup := domain.Group{
		ID:   uuid.NewString(),
		Name: name,
		Type: groupType,
	}
	err = g.groupRepo.Create(ctx, &newGroup)

	if err != nil {
		return err
	}

	err = g.groupRepo.AddMember(ctx, newGroup.ID, user.ID, domain.GroupRoleAdmin)

	return err
}
