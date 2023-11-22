package usecase

import (
	stderrors "errors"
	"github.com/dgrijalva/jwt-go"
	"semi_systems/config"
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/resource/request"
	"semi_systems/keijiban/resource/response"
	"semi_systems/packages/context"
	"semi_systems/packages/errors"
	"strconv"
	"time"
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
	GetUserByName(ctx context.Context, name string) (*domain.User, error)
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
	newUser, err := domain.NewUser(ctx, req)
	if err != nil {
		return err
	}

	if ctx.IsInValid() {
		return ctx.ValidationError()
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

	if req.Name != "" {
		user.Name = req.Name
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
	user, err := u.UserRepo.GetUserByName(ctx, req.Name)
	if err != nil {
		return err
	}

	if user.Password.IsValid(req.Password) {
		token, refreshToken, err := issueJWTToken(strconv.Itoa(int(user.ID)), config.Env.App.Secret)
		if err != nil {
			return errors.NewUnexpected(err)
		}

		var res response.UserLogin
		res.Token = token
		res.RefreshToken = refreshToken

		return u.outputPort.Login(req.Session, &res)
	}
	return u.outputPort.Login(req.Session, nil)
}

func issueJWTToken(userID string, secretKey string) (string, string, error) {
	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	// リフレッシュトークンの生成
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func (u User) RefreshToken(req *request.UserRefreshToken) error {
	var res response.UserLogin

	// リフレッシュトークンの検証
	claims, err := verifyToken(req.RefreshToken, config.Env.App.Secret)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	if claims == nil {
		return nil // トークンが無効な場合は何もしない
	}

	// 新しいトークンとリフレッシュトークンを生成
	newToken, newRefreshToken, err := issueJWTToken(claims["uid"].(string), config.Env.App.Secret)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	res.Token = newToken
	res.RefreshToken = newRefreshToken

	return u.outputPort.RefreshToken(req.Session, &res)
}

func verifyToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, stderrors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, stderrors.New("invalid token")
}
