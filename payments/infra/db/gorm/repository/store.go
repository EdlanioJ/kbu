package repository

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/jinzhu/gorm"
)

type StoreRepositoryGORM struct {
	DB *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepositoryGORM {
	return &StoreRepositoryGORM{
		DB: db,
	}
}

func (s *StoreRepositoryGORM) Find(id string) (*entity.Store, error) {
	store := &entity.Store{}

	err := s.DB.First(store, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return store, nil
}
func (s *StoreRepositoryGORM) FindStoreByIdAndStatus(id string, status string) (*entity.Store, error) {
	store := &entity.Store{}

	err := s.DB.First(store, "id = ? AND status = ?", id, status).Error

	if err != nil {
		return nil, err
	}

	return store, nil
}
