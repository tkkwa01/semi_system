// usecase/attendance_usecase.go

package usecase

import (
	"semi_systems/attendance/domain"
	"semi_systems/packages/context"
)

type AttendanceInputPort interface {
	GetAllAttendance(ctx context.Context) error
	UpdateStatus(ctx context.Context, name string, status bool) error
}

type AttendanceOutputPort interface {
	GetAllAttendance(res []*domain.Attendance) error
	UpdateStatus(success bool) error
}

type AttendanceRepository interface {
	GetAll(ctx context.Context) ([]*domain.Attendance, error)
	UpdateStatus(ctx context.Context, name string, status bool) error
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

func (a *attendance) GetAllAttendance(ctx context.Context) error {
	attendances, err := a.AttendanceRepo.GetAll(ctx)
	if err != nil {
		return err
	}
	return a.outputPort.GetAllAttendance(attendances)
}

func (a *attendance) UpdateStatus(ctx context.Context, name string, status bool) error {
	err := a.AttendanceRepo.UpdateStatus(ctx, name, status)
	if err != nil {
		return err
	}

	return a.outputPort.UpdateStatus(true)
}
