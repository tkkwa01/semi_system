// api/handler/attendance_handler.go

package handler

import (
	"github.com/gin-gonic/gin"
	"semi_systems/kobasemi_attendance_API/domain"
	"semi_systems/kobasemi_attendance_API/usecase"
)

type AttendanceHandler struct {
	attendanceUsecase *usecase.AttendanceUsecase
}

func NewAttendanceHandler(uc *usecase.AttendanceUsecase) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceUsecase: uc,
	}
}

func (h *AttendanceHandler) RegisterAttendance(c *gin.Context) {
	var attendance domain.Attendance
	if err := c.BindJSON(&attendance); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}
	if err := h.attendanceUsecase.RegisterAttendance(&attendance); err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, gin.H{"message": "Registered"})
}

func (h *AttendanceHandler) GetAllAttendances(c *gin.Context) {
	attendances, err := h.attendanceUsecase.GetAllAttendances()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, attendances)
}

func (h *AttendanceHandler) UpdateAttendance(c *gin.Context) {
	var update domain.Attendance
	if err := c.BindJSON(&update); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	name := update.Name
	status := update.Status

	// 成功した場合はtrue、該当するnameがない場合falseを返す
	if success, err := h.attendanceUsecase.UpdateAttendanceStatus(name, status); err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	} else if !success {
		c.JSON(404, gin.H{"error": "Name not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Status updated"})
}
