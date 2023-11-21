package domain

import (
	"semi_systems/keijiban/resource/request"
	"semi_systems/packages/context"
)

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Introduction string `json:"introduction"`
}

func NewUser(ctx context.Context, dto *request.UserCreate) (*User, error) {
	var user = User{
		Email:        dto.Email,
		Introduction: dto.Introduction,
	}

	if ctx.IsInValid() {
		return nil, ctx.ValidationError()
	}

	return &user, nil
}
