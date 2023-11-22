package presenter

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	u.c.JSON(http.StatusCreated, gin.H{"id": id})
	return nil
}

func (u user) GetAllUser(res []*domain.User) error {
	u.c.JSON(http.StatusOK, res)
	return nil
}

func (u user) GetUserByID(res domain.User) error {
	u.c.JSON(http.StatusOK, res)
	return nil
}

func (u user) UpdateUser(res domain.User) error {
	u.c.JSON(http.StatusOK, res)
	return nil
}

func (u user) DeleteUser() error {
	u.c.JSON(http.StatusOK, "")
	return nil
}
