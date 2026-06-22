package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type ProbationService interface {
	CreateProbation(row models.Probation) (int, dto.Error)
	GetProbationAdmin(row []models.Probation, page dto.Pagination, fliter dto.ProbationFilter) ([]models.Probation, int, int64, dto.Error)
	GetProbation(row []models.Probation, page dto.Pagination, fliter dto.ProbationFilter) ([]models.Probation, int, int64, dto.Error)
	GetProbationById(row models.Probation, filter dto.GetByID) (models.Probation, int, dto.Error)
	DeleteProbation(row models.Probation, id uuid.UUID) (int, dto.Error)
	PatchProbation(row models.Probation, id uuid.UUID) (models.Probation, int, dto.Error)
}

type Probationservice struct {
	repo repository.ProbationRepository
}

func NewProbationService(repo repository.ProbationRepository) ProbationService {

	return &Probationservice{
		repo: repo,
	}

}

func (s *Probationservice) CreateProbation(row models.Probation) (int, dto.Error) {

	return s.repo.CreateProbation(row)

}
func (s *Probationservice) GetProbationAdmin(row []models.Probation, page dto.Pagination, fliter dto.ProbationFilter) ([]models.Probation, int, int64, dto.Error) {

	return s.repo.GetProbationAdmin(row, page, fliter)

}

func (s *Probationservice) GetProbation(row []models.Probation, page dto.Pagination, fliter dto.ProbationFilter) ([]models.Probation, int, int64, dto.Error) {

	return s.repo.GetProbation(row, page, fliter)

}

func (s *Probationservice) GetProbationById(row models.Probation, filter dto.GetByID) (models.Probation, int, dto.Error) {

	return s.repo.GetProbationById(row, filter)

}

func (s *Probationservice) PatchProbation(row models.Probation, id uuid.UUID) (models.Probation, int, dto.Error) {

	return s.repo.PatchProbation(row, id)

}

func (s *Probationservice) DeleteProbation(row models.Probation, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeleteProbation(row, id)

}
