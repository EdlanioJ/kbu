package repository

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type StoreRepository interface {
	Find(id string) (*entity.Store, error)
	FindStoreByIdAndStatus(id string, status string) (*entity.Store, error)
}
