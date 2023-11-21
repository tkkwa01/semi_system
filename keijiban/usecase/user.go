package usecase

import (
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/resource/request"
	"semi_systems/packages/context"
)

type UserInputPort interface {
	CreateUser(ctx context.Context, req *request.UserCreate) error
	GetAllUser(ctx context.Context) error
	GetUserByID(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, req *request.UserUpdate) error
	DeleteUser(ctx context.Context, id uint) error
}

type UserOutputPort interface {
	CreateUser(id uint) error
	GetAllUser(res []*domain.User) error
	GetUserByID(res domain.User) error
	UpdateUser(res domain.User) error
	DeleteUser() error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (uint, error)
	GetAllUser(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type User struct {
	UserRepo   UserRepository
	outputPort UserOutputPort
}

type UsereInputFactory func(outputPort UserOutputPort) UserInputPort

func NewUserInputFactory(ur UserRepository) UsereInputFactory {
	return func(outputPort UserOutputPort) UserInputPort {
		return &User{
			UserRepo:   ur,
			outputPort: outputPort,
		}
	}
}

func (u User) CreateUser(ctx context.Context, req *request.UserCreate) error {
	newUser := &domain.User{
		Email:        req.Email,
		Introduction: req.Introduction,
	}

	id, err := u.UserRepo.CreateUser(ctx, newUser)
	if err != nil {
		return err
	}

	return u.outputPort.CreateUser(id)
}

func (u User) GetAllUser(ctx context.Context) error {
	users, err := u.UserRepo.GetAllUser(ctx)
	if err != nil {
		return err
	}

	return u.outputPort.GetAllUser(users)
}

func (u User) GetUserByID(ctx context.Context, id uint) error {
	user, err := u.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	return u.outputPort.GetUserByID(*user)
}

func (u User) UpdateUser(ctx context.Context, req *request.UserUpdate) error {
	user, err := u.UserRepo.GetUserByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Introduction != "" {
		user.Introduction = req.Introduction
	}

	err = u.UserRepo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return u.outputPort.UpdateUser(*user)
}

func (u User) DeleteUser(ctx context.Context, id uint) error {
	err := u.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return u.outputPort.DeleteUser()
}
