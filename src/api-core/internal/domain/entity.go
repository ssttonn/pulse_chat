package domain

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupType string

const (
	GroupTypeGroup  GroupType = "group"
	GroupTypeDirect GroupType = "direct"
)

type Group struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      GroupType `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupMemberRole string

const (
	GroupRoleAdmin  GroupMemberRole = "admin"
	GroupRoleMember GroupMemberRole = "member"
)

type GroupMember struct {
	UserId   string          `json:"user_id"`
	GroupId  string          `json:"group_id"`
	Role     GroupMemberRole `json:"role"`
	JoinedAt time.Time       `json:"joined_at"`
}

var ErrUserNotFound error = errors.New("user not found")
var ErrGroupNotFound error = errors.New("group not found")
