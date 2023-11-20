// usecase/attendance_usecase.go

package usecase

import (
	"semi_systems/kobasemi_attendance_API/domain"
	"semi_systems/kobasemi_attendance_API/repository"
)

type AttendanceUsecase struct {
	attendanceRepo repository.AttendanceRepository
}

func NewAttendanceUsecase(repo repository.AttendanceRepository) *AttendanceUsecase {
	return &AttendanceUsecase{
		attendanceRepo: repo,
	}
}

func (uc *AttendanceUsecase) GetAllAttendances() ([]domain.Attendance, error) {
	return uc.attendanceRepo.GetAll()
}

func (uc *AttendanceUsecase) RegisterAttendance(attendance *domain.Attendance) error {
	return uc.attendanceRepo.Register(attendance)
}

func (uc *AttendanceUsecase) UpdateAttendanceStatus(name string, status bool) (bool, error) {
	exists, err := uc.attendanceRepo.Exists(name)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	if err := uc.attendanceRepo.UpdateStatus(name, status); err != nil {
		return false, err
	}

	return true, nil
}
