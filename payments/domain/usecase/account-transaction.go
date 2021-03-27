package usecase

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type AccountTransaction interface {
	RegisterAccountTransaction(fromAccountId string, toAccountId string, amount float64, currency string) (*entity.Transaction, error)
	FindAllByAccountTo(accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error)
	FindOneByAccount(accountFromId string, transactionId string) (*entity.Transaction, error)
}
