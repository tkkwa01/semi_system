package repository

import (
	"gorm.io/gorm"
	"semi_systems/kobasemi_attendance_API/domain"
)

type AttendanceRepository interface {
	GetAll() ([]domain.Attendance, error)
	Register(attendance *domain.Attendance) error
	UpdateStatus(name string, status bool) error
	Exists(name string) (bool, error)
}

type AttendanceGormRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &AttendanceGormRepository{db: db}
}

func (repo *AttendanceGormRepository) GetAll() ([]domain.Attendance, error) {
	var attendances []domain.Attendance
	if err := repo.db.Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (repo *AttendanceGormRepository) Register(attendance *domain.Attendance) error {
	return repo.db.Create(attendance).Error
}

func (repo *AttendanceGormRepository) UpdateStatus(name string, status bool) error {
	return repo.db.Model(&domain.Attendance{}).Where("name = ?", name).Update("status", status).Error
}

func (repo *AttendanceGormRepository) Exists(name string) (bool, error) {
	var count int64
	err := repo.db.Model(&domain.Attendance{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
