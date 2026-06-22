package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type RolesService interface {
	CreateRoles(row models.Roles) (int, dto.Error)
	GetRoles(row []models.Roles, page dto.Pagination, filter dto.RoleFilter) ([]models.Roles, int, int64, dto.Error)
	GetRolesById(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error)
	DeleteRoles(row models.Roles, id uuid.UUID) (int, dto.Error)
	PatchRoles(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error)
}

type roleservice struct {
	repo repository.RolesRepository
}

func NewRoleService(repo repository.RolesRepository) RolesService {

	return &roleservice{
		repo: repo,
	}

}

func (s *roleservice) CreateRoles(row models.Roles) (int, dto.Error) {

	return s.repo.CreateRoles(row)

}
func (s *roleservice) GetRoles(row []models.Roles, page dto.Pagination, filter dto.RoleFilter) ([]models.Roles, int, int64, dto.Error) {

	return s.repo.GetRoles(row, page, filter)

}

func (s *roleservice) GetRolesById(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error) {

	return s.repo.GetRolesById(row, id)

}

func (s *roleservice) PatchRoles(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error) {

	return s.repo.PatchRoles(row, id)

}

func (s *roleservice) DeleteRoles(row models.Roles, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeleteRoles(row, id)

}
