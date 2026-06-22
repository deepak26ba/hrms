package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type PositionService interface {
	CreatePosition(row models.Position) (int, dto.Error)
	GetPosition(row []models.Position, page dto.Pagination, filter dto.PositionFilter) ([]models.Position, int, int64, dto.Error)
	GetPositionById(row models.Position, id uuid.UUID) (models.Position, int, dto.Error)
	DeletePosition(row models.Position, id uuid.UUID) (int, dto.Error)
	PatchPosition(row models.Position, id uuid.UUID) (models.Position, int, dto.Error)
}

type positionservice struct {
	repo repository.PositionRepository
}

func NewPositionService(repo repository.PositionRepository) PositionService {

	return &positionservice{
		repo: repo,
	}

}

func (s *positionservice) CreatePosition(row models.Position) (int, dto.Error) {

	return s.repo.CreatePosition(row)

}

func (s *positionservice) GetPosition(row []models.Position, page dto.Pagination, filter dto.PositionFilter) ([]models.Position, int, int64, dto.Error) {

	return s.repo.GetPosition(row, page, filter)

}

func (s *positionservice) GetPositionById(row models.Position, id uuid.UUID) (models.Position, int, dto.Error) {

	return s.repo.GetPositionById(row, id)

}

func (s *positionservice) PatchPosition(row models.Position, id uuid.UUID) (models.Position, int, dto.Error) {

	return s.repo.PatchPosition(row, id)

}

func (s *positionservice) DeletePosition(row models.Position, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeletePosition(row, id)

}
