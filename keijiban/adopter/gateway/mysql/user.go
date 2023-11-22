package mysql

import (
	"gorm.io/gorm"
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/usecase"
	"semi_systems/packages/context"
	"semi_systems/packages/errors"
)

type user struct{}

func NewUserRepository() usecase.UserRepository {
	return &user{}
}

func dbError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errors.NotFound()
	default:
		return errors.NewUnexpected(err)
	}
}

func (u user) CreateUser(ctx context.Context, user *domain.User) (uint, error) {
	db := ctx.DB()

	if err := db.Create(user).Error; err != nil {
		return 0, dbError(err)
	}
	return user.ID, nil
}

func (u user) GetAllUser(ctx context.Context) ([]*domain.User, error) {
	db := ctx.DB()

	var users []*domain.User
	if err := db.Find(&users).Error; err != nil {
		return nil, dbError(err)
	}
	return users, nil
}

func (u user) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	db := ctx.DB()

	var user domain.User
	err := db.Where(&domain.User{ID: id}).First(&user).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &user, nil
}

func (u user) UpdateUser(ctx context.Context, user *domain.User) error {
	db := ctx.DB()

	if err := db.Model(&user).Updates(user).Error; err != nil {
		return dbError(err)
	}
	return nil
}

func (u user) DeleteUser(ctx context.Context, id uint) error {
	db := ctx.DB()

	var user domain.User
	res := db.Where("id = ?", id).Delete(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
