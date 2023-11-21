package http

import (
	"github.com/gin-gonic/gin"
	"semi_systems/attendance/adopter/presenter"
	"semi_systems/attendance/domain"
	"semi_systems/attendance/usecase"
	"semi_systems/packages/context"
	"semi_systems/packages/http/router"
)

type attendance struct {
	inputFactory   usecase.AttendanceInputFactory
	outputFactory  func(c *gin.Context) usecase.AttendanceOutputPort
	AttendanceRepo usecase.AttendanceRepository
}

func NewAttendance(r *router.Router, inputFactory usecase.AttendanceInputFactory, outputFactory presenter.AttendanceOutputFactory) {
	handler := attendance{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
	}

	r.Group("attendances", nil, func(r *router.Router) {
		r.Get("/watch", handler.GetAll)
		r.Post("/register", handler.UpdateStatus)
	})
}

func (a attendance) GetAll(ctx context.Context, c *gin.Context) error {
	outputPort := a.outputFactory(c)
	inputPort := a.inputFactory(outputPort)

	return inputPort.GetAllAttendance(ctx)
}

func (a attendance) UpdateStatus(ctx context.Context, c *gin.Context) error {
	var update domain.Attendance
	if err := c.BindJSON(&update); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return err
	}

	outputPort := a.outputFactory(c)
	inputPort := a.inputFactory(outputPort)

	return inputPort.UpdateStatus(ctx, update.Name, update.Status)
}
