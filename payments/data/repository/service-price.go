package repository

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type ServicePriceRepository interface {
	Find(id string) (*entity.ServicePrice, error)
}
