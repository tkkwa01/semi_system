package presenter

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"semi_systems/attendance/domain"
	"semi_systems/attendance/usecase"
)

type attendance struct {
	c *gin.Context
}

type AttendanceOutputFactory func(c *gin.Context) usecase.AttendanceOutputPort

func NewAttendanceOutputFactory() AttendanceOutputFactory {
	return func(c *gin.Context) usecase.AttendanceOutputPort {
		return &attendance{c: c}
	}
}

func (a attendance) GetAllAttendance(res []*domain.Attendance) error {
	a.c.JSON(http.StatusOK, res)
	return nil
}

func (a attendance) UpdateStatus(success bool) error {
	a.c.JSON(http.StatusOK, gin.H{"success": success})
	return nil
}
