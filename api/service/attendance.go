package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type AttendanceService interface {
	CreateAttendance(row models.Attendance) (uuid.UUID, int, dto.Error)
	GetAttendanceAdmin(row []models.Attendance, page dto.Pagination, fliter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error)
	GetAttendance(row []models.Attendance, page dto.Pagination, fliter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error)
	GetAttendanceById(row models.Attendance, filter dto.GetByID) (models.Attendance, int, dto.Error)
	DeleteAttendance(row models.Attendance, id uuid.UUID) (int, dto.Error)
	PatchAttendance(row models.Attendance, id uuid.UUID) (models.Attendance, int, dto.Error)
}

type attendanceservice struct {
	repo repository.AttendanceRepository
}

func NewAttendanceService(repo repository.AttendanceRepository) AttendanceService {

	return &attendanceservice{
		repo: repo,
	}

}

func (s *attendanceservice) CreateAttendance(row models.Attendance) (uuid.UUID, int, dto.Error) {

	return s.repo.CreateAttendance(row)

}
func (s *attendanceservice) GetAttendanceAdmin(row []models.Attendance, page dto.Pagination, fliter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error) {

	return s.repo.GetAttendanceAdmin(row, page, fliter)

}

func (s *attendanceservice) GetAttendance(row []models.Attendance, page dto.Pagination, fliter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error) {

	return s.repo.GetAttendance(row, page, fliter)

}

func (s *attendanceservice) GetAttendanceById(row models.Attendance, filter dto.GetByID) (models.Attendance, int, dto.Error) {

	return s.repo.GetAttendanceById(row, filter)

}

func (s *attendanceservice) PatchAttendance(row models.Attendance, id uuid.UUID) (models.Attendance, int, dto.Error) {

	return s.repo.PatchAttendance(row, id)

}

func (s *attendanceservice) DeleteAttendance(row models.Attendance, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeleteAttendance(row, id)

}
