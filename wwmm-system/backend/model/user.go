package model

import "time"

type User struct {
	UserID       int       `json:"userId"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Salt         string    `json:"-"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	RealName     string    `json:"realName"`
	Sex          int       `json:"sex"`
	Avatar       string    `json:"avatar"`
	Bio          string    `json:"bio"`
	Role         int       `json:"role"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

const (
	RoleVoter        = 0
	RolePhotographer = 1
	RoleAdmin        = 2
)

const (
	StatusDisabled = 0
	StatusEnabled  = 1
)
