package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type AttendanceAuditService interface {
	CreateAttendanceAudit(row models.AttendanceAudit) (int, dto.Error)
	GetAttendanceAuditAdmin(row []models.AttendanceAudit, page dto.Pagination, fliter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error)
	GetAttendanceAudit(row []models.AttendanceAudit, page dto.Pagination, fliter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error)
	GetAttendanceAuditById(row models.AttendanceAudit, filter dto.GetByID) (models.AttendanceAudit, int, dto.Error)
	PatchAttendanceAudit(row models.AttendanceAudit, userID uuid.UUID) (models.AttendanceAudit, int, dto.Error)
}

type attendanceAuditservice struct {
	repo repository.AttendanceAuditRepository
}

func NewAttendanceAuditService(repo repository.AttendanceAuditRepository) AttendanceAuditService {

	return &attendanceAuditservice{
		repo: repo,
	}

}

func (s *attendanceAuditservice) CreateAttendanceAudit(row models.AttendanceAudit) (int, dto.Error) {

	return s.repo.CreateAttendanceAudit(row)

}
func (s *attendanceAuditservice) GetAttendanceAuditAdmin(row []models.AttendanceAudit, page dto.Pagination, fliter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error) {

	return s.repo.GetAttendanceAuditAdmin(row, page, fliter)

}

func (s *attendanceAuditservice) GetAttendanceAudit(row []models.AttendanceAudit, page dto.Pagination, fliter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error) {

	return s.repo.GetAttendanceAudit(row, page, fliter)

}

func (s *attendanceAuditservice) GetAttendanceAuditById(row models.AttendanceAudit, filter dto.GetByID) (models.AttendanceAudit, int, dto.Error) {

	return s.repo.GetAttendanceAuditById(row, filter)

}

func (s *attendanceAuditservice) PatchAttendanceAudit(row models.AttendanceAudit, userID uuid.UUID) (models.AttendanceAudit, int, dto.Error) {

	return s.repo.PatchAttendanceAudit(row, userID)

}
