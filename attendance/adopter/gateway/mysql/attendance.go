package mysql

import (
	"gorm.io/gorm"
	"semi_systems/attendance/domain"
	"semi_systems/attendance/usecase"
	"semi_systems/packages/errors"
)

type attendance struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) usecase.AttendanceRepository {
	return &attendance{db: db}
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

func (a attendance) GetAll() ([]*domain.Attendance, error) {
	var attendances []*domain.Attendance
	if err := a.db.Find(&attendances).Error; err != nil {
		return nil, dbError(err)
	}
	return attendances, nil
}

func (a *attendance) UpdateStatus(name string, status bool) error {
	return a.db.Model(&domain.Attendance{}).Where("name = ?", name).Update("status", status).Error
}
