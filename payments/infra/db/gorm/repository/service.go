package repository

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/jinzhu/gorm"
)

type ServiceRepositoryGORM struct {
	DB *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepositoryGORM {
	return &ServiceRepositoryGORM{
		DB: db,
	}
}

func (s *ServiceRepositoryGORM) Find(id string) (*entity.Service, error) {
	service := &entity.Service{}

	err := s.DB.First(service, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return service, nil
}
func (s *ServiceRepositoryGORM) FindServiceByIdAndStatus(id string, status string) (*entity.Service, error) {
	service := &entity.Service{}

	err := s.DB.First(service, "id = ? AND status = ?", id, status).Error

	if err != nil {
		return nil, err
	}

	return service, nil
}
