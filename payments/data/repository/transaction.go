package repository

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type TransactionRepository interface {
	Register(transaction *entity.Transaction) error
	Save(transaction *entity.Transaction) error
	Find(id string) (*entity.Transaction, error)
	FindAll(pagination *entity.Pagination) ([]*entity.Transaction, int, error)
	FindByType(transactionID, transactionType string) (*entity.Transaction, error)
	FindAllByType(transactionType string, pagination *entity.Pagination) ([]*entity.Transaction, int, error)
	FindByExternalID(transactionID, ExternalID string) (*entity.Transaction, error)
	FindAllByExternalID(externalID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error)
	FindByFromAccountID(transactionID, accountID string) (*entity.Transaction, error)
	FindAllByFromAccountID(accountID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error)
	FindByToAccountID(transactionID, accountID string) (*entity.Transaction, error)
	FindAllByToAccountID(accountID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error)
}
