package usecase

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type Transaction interface {
	Register(fromAccount, toAccount, externalID, typeTransaction, currency string, amount float64) (*entity.Transaction, error)
	Find(id string) (*entity.Transaction, error)
	FindAll(page int, limit int, sort string) ([]*entity.Transaction, int, error)
	FindByType(typeTransaction, transactionID string) (*entity.Transaction, error)
	FindAllByType(typeTransaction string, page int, limit int, sort string) ([]*entity.Transaction, int, error)
	FindByExternalID(externalID, transactionID string) (*entity.Transaction, error)
	FindAllByExternalID(externalID string, page int, limit int, sort string) ([]*entity.Transaction, int, error)
	FindAllByFromAccountID(accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error)
	FindByFromAccountID(accountID, transactionID string) (*entity.Transaction, error)
	FindAllByToAccountID(accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error)
	FindByToAccountID(accountID, transactionID string) (*entity.Transaction, error)
	Complete(transactionID string) (*entity.Transaction, error)
	Error(transactionID string) (*entity.Transaction, error)
}
