package domain

import (
	"semi_systems/keijiban/domain/vobj"
	"semi_systems/keijiban/resource/request"
	"semi_systems/packages/context"
)

type User struct {
	ID            uint                `json:"id"`
	Email         string              `json:"email"`
	Password      vobj.Password       `json:"-"`
	RecoveryToken *vobj.RecoveryToken `json:"-"`
}

func NewUser(ctx context.Context, dto *request.UserCreate) (*User, error) {
	var user = User{
		Email:         dto.Email,
		RecoveryToken: vobj.NewRecoveryToken(""),
	}

	if ctx.IsInValid() {
		return nil, ctx.ValidationError()
	}

	password, err := vobj.NewPassword(ctx, dto.Password, dto.PasswordConfirm)
	if err != nil {
		return nil, err
	}

	user.Password = *password

	return &user, nil
}
