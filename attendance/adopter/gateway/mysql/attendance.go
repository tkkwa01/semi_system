package mysql

import (
	"gorm.io/gorm"
	"semi_systems/attendance/domain"
	"semi_systems/attendance/usecase"
	"semi_systems/packages/context"
	"semi_systems/packages/errors"
)

type attendance struct{}

func NewAttendanceRepository() usecase.AttendanceRepository {
	return &attendance{}
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

func (a attendance) GetAll(ctx context.Context) ([]*domain.Attendance, error) {
	db := ctx.DB()

	var attendances []*domain.Attendance
	if err := db.Find(&attendances).Error; err != nil {
		return nil, dbError(err)
	}
	return attendances, nil
}

func (a *attendance) UpdateStatus(ctx context.Context, name string, status bool) error {
	db := ctx.DB()
	return db.Model(&domain.Attendance{}).Where("name = ?", name).Update("status", status).Error
}
