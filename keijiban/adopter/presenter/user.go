package presenter

import (
	"github.com/gin-gonic/gin"
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/usecase"
)

type user struct {
	c *gin.Context
}

type UserOutputFactory func(c *gin.Context) usecase.UserOutputPort

func NewUserOutputFactory() UserOutputFactory {
	return func(c *gin.Context) usecase.UserOutputPort {
		return &user{c: c}
	}
}

func (u user) CreateUser(id uint) error {
	//TODO implement me
	panic("implement me")
}

func (u user) GetAllUser(res []*domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u user) GetUserByID(res domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u user) UpdateUser(res domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u user) DeleteUser() error {
	//TODO implement me
	panic("implement me")
}
