package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"semi_systems/config"
	"semi_systems/keijiban/adopter/presenter"
	"semi_systems/keijiban/resource/request"
	"semi_systems/keijiban/usecase"
	"semi_systems/packages/context"
	"semi_systems/packages/http/middleware"
	"semi_systems/packages/http/router"
	"strconv"
)

type user struct {
	inputFactory  usecase.UserInputFactory
	outputFactory func(c *gin.Context) usecase.UserOutputPort
	UserRepo      usecase.UserRepository
}

func NewUser(r *router.Router, inputFactory usecase.UserInputFactory, outputFactory presenter.UserOutputFactory) {
	handler := user{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
	}

	r.Group("user", nil, func(r *router.Router) {
		r.Post("", handler.CreateUser)
		r.Post("login", handler.Login)
		r.Post("refresh-token", handler.RefreshToken)
	})

	r.Group("", []gin.HandlerFunc{middleware.Auth(true, config.UserRealm, true)}, func(r *router.Router) {
		r.Group("user", nil, func(r *router.Router) {
			r.Get("me", handler.GetMe)
		})
	})
}

func (u user) CreateUser(ctx context.Context, c *gin.Context) error {
	var req request.UserCreate

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.CreateUser(ctx, &req)
}

func (u user) Login(ctx context.Context, c *gin.Context) error {
	var req request.UserLogin

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.Login(ctx, &req)
}

func (u user) RefreshToken(_ context.Context, c *gin.Context) error {
	var req request.UserRefreshToken

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.RefreshToken(&req)
}

func (u user) GetMe(ctx context.Context, c *gin.Context) error {
	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.GetUserByID(ctx, ctx.UID())
}

//ここより下のメソッドは今回は使わないことにした

func (u user) GetAll(ctx context.Context, c *gin.Context) error {
	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.GetAllUser(ctx)
}

func (u user) GetUserByID(ctx context.Context, c *gin.Context) error {
	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	return inputPort.GetUserByID(ctx, uint(id))
}

func (u user) UpdateUser(ctx context.Context, c *gin.Context) error {
	var req request.UserUpdate

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return err
	}
	req.ID = uint(id)

	// リクエストボディをバインド
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return err
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.UpdateUser(ctx, &req)
}

func (u user) DeleteUser(ctx context.Context, c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return err
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.DeleteUser(ctx, uint(id))
}
