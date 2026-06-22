package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type UserService interface {
	GetUser(row []models.User, page dto.Pagination, fliter dto.UserFilter) ([]models.User, int, int64, dto.Error)
	GetUserById(row models.User, id uuid.UUID) (models.User, int, dto.Error)
	DeleteUser(row models.User, id uuid.UUID) (int, dto.Error)
	PatchUser(row models.User, id uuid.UUID) (models.User, int, dto.Error)
}

type userservice struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {

	return &userservice{
		repo: repo,
	}

}

func (s *userservice) GetUser(row []models.User, page dto.Pagination, filter dto.UserFilter) ([]models.User, int, int64, dto.Error) {

	return s.repo.GetUser(row, page, filter)

}

func (s *userservice) GetUserById(row models.User, id uuid.UUID) (models.User, int, dto.Error) {

	return s.repo.GetUserById(row, id)

}

func (s *userservice) PatchUser(row models.User, id uuid.UUID) (models.User, int, dto.Error) {

	return s.repo.PatchUser(row, id)

}

func (s *userservice) DeleteUser(row models.User, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeleteUser(row, id)

}
