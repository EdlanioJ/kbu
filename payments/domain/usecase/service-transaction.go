package usecase

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type ServiceTransaction interface {
	RegisterServiceTransaction(fromId string, serviceId string, servicePriceId string, amount float64, currency string) (*entity.Transaction, error)
	FindAllByServiceId(serviceId string, page int, limit int, sort string) ([]*entity.Transaction, int, error)
	FindOneByService(serviceId string, transactionId string) (*entity.Transaction, error)
}
