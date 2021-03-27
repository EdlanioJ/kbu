package usecase

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type Transaction interface {
	Find(id string) (*entity.Transaction, error)
	FindAll(page int, limit int, sort string) ([]*entity.Transaction, int, error)
	Complete(transactionId string) (*entity.Transaction, error)
	Error(transactionId string) (*entity.Transaction, error)
}
