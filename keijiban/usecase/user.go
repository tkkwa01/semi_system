package usecase

import (
	jwt "github.com/ken109/gin-jwt"
	"semi_systems/config"
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/resource/request"
	"semi_systems/keijiban/resource/response"
	"semi_systems/packages/context"
	"semi_systems/packages/errors"
)

type UserInputPort interface {
	CreateUser(ctx context.Context, req *request.UserCreate) error
	GetAllUser(ctx context.Context) error
	GetUserByID(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, req *request.UserUpdate) error
	DeleteUser(ctx context.Context, id uint) error
	Login(ctx context.Context, req *request.UserLogin) error
	RefreshToken(req *request.UserRefreshToken) error
}

type UserOutputPort interface {
	CreateUser(id uint) error
	GetAllUser(res []*domain.User) error
	GetUserByID(res domain.User) error
	UpdateUser(res domain.User) error
	DeleteUser() error
	Login(isSession bool, res *response.UserLogin) error
	RefreshToken(isSession bool, res *response.UserLogin) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (uint, error)
	GetAllUser(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	EmailExist(ctx context.Context, email string) (bool, error)
	GetByRecoveryToken(ctx context.Context, recoverToken string) (*domain.User, error)
}

type User struct {
	UserRepo   UserRepository
	outputPort UserOutputPort
}

type UserInputFactory func(outputPort UserOutputPort) UserInputPort

func NewUserInputFactory(ur UserRepository) UserInputFactory {
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

func (u User) Login(ctx context.Context, req *request.UserLogin) error {
	user, err := u.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if user.Password.IsValid(req.Password) {
		var res response.UserLogin
		res.Token, res.RefreshToken, err = jwt.IssueToken(config.UserRealm, jwt.Claims{
			"uid": user.ID,
		})
		if err != nil {
			return errors.NewUnexpected(err)
		}
		return u.outputPort.Login(req.Session, &res)
	}
	return u.outputPort.Login(req.Session, nil)
}

func (u User) RefreshToken(req *request.UserRefreshToken) error {
	var (
		res response.UserLogin
		ok  bool
		err error
	)

	ok, res.Token, res.RefreshToken, err = jwt.RefreshToken(config.UserRealm, req.RefreshToken)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	if !ok {
		return nil
	}
	return u.outputPort.RefreshToken(req.Session, &res)
}
