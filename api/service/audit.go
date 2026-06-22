package service

import (
	"hrms/api/repository"
	"hrms/pkg/models"
)

type AuditService interface {
	CreateAudit(log models.AuditTable)
}

type auditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) AuditService {

	return &auditService{
		repo: repo,
	}
}

func (s *auditService) CreateAudit(log models.AuditTable) {
	s.repo.CreateAudit(log)
}
