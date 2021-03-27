package repository

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type ServiceRepository interface {
	Find(id string) (*entity.Service, error)
	FindServiceByIdAndStatus(id string, status string) (*entity.Service, error)
}
