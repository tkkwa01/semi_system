package context

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"semi_systems/packages/errors"
)

type Context interface {
	RequestID() string
	Authenticated() bool
	UID() uint
	UserName() string
	Validate(request interface{}) (invalid bool)
	FieldError(fieldName string, message string)
	IsInValid() bool
	ValidationError() error
	DB() *gorm.DB
	Transaction(fn func(ctx Context) error) error
	Error() *errors.SubError
}

type ctx struct {
	id       string
	verr     *errors.Error
	getDB    func() *gorm.DB
	db       *gorm.DB
	uid      uint
	userName string
	err      *errors.SubError
}

func New(c *gin.Context, getDB func() *gorm.DB) Context {
	requestID := c.GetHeader("X-Request-Id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	var uid uint
	var userName string

	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			uid = id
			log.Printf("New context: userID is %v", uid)
		} else {
			log.Printf("New context: userID is not uint, actual type: %T", userID)
		}
	}
	if name, exists := c.Get("user_name"); exists {
		if name, ok := name.(string); ok {
			userName = name
		}
	}

	return &ctx{
		id:       requestID,
		verr:     errors.NewValidation(),
		getDB:    getDB,
		uid:      uid,
		userName: userName,
		err:      nil,
	}
}

func (c ctx) RequestID() string {
	return c.id
}

func (c ctx) Authenticated() bool {
	return c.uid != 0
}

func (c ctx) UID() uint {
	return c.uid
}

func (c ctx) UserName() string {
	return c.userName
}

func (c ctx) Error() *errors.SubError {
	if c.err == nil {
		c.err = errors.New()
	}
	return c.err
}
