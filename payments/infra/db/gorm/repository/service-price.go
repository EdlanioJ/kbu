package repository

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/jinzhu/gorm"
)

type ServicePriceRepositoryGORM struct {
	DB *gorm.DB
}

func NewServicePriceRepository(db *gorm.DB) *ServicePriceRepositoryGORM {
	return &ServicePriceRepositoryGORM{
		DB: db,
	}
}

func (s *ServicePriceRepositoryGORM) Find(id string) (*entity.ServicePrice, error) {
	servicePrice := &entity.ServicePrice{}

	err := s.DB.First(servicePrice, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return servicePrice, nil
}
