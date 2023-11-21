// usecase/attendance_usecase.go

package usecase

import (
	"semi_systems/attendance/domain"
)

type AttendanceInputPort interface {
	GetAllAttendance() error
	UpdateStatus(name string, status bool) error
}

type AttendanceOutputPort interface {
	GetAllAttendance(res []*domain.Attendance) error
	UpdateStatus(success bool) error
}

type AttendanceRepository interface {
	GetAll() ([]*domain.Attendance, error)
	UpdateStatus(name string, status bool) error
}

type attendance struct {
	outputPort     AttendanceOutputPort
	AttendanceRepo AttendanceRepository
}

type AttendanceInputFactory func(outputPort AttendanceOutputPort) AttendanceInputPort

func NewAttendanceInputFactory(ar AttendanceRepository) AttendanceInputFactory {
	return func(o AttendanceOutputPort) AttendanceInputPort {
		return &attendance{
			outputPort:     o,
			AttendanceRepo: ar,
		}
	}
}

func (a *attendance) GetAllAttendance() error {
	attendances, err := a.AttendanceRepo.GetAll()
	if err != nil {
		return err
	}
	return a.outputPort.GetAllAttendance(attendances)
}

func (a *attendance) UpdateStatus(name string, status bool) error {
	err := a.AttendanceRepo.UpdateStatus(name, status)
	if err != nil {
		return err
	}

	return a.outputPort.UpdateStatus(true)
}
