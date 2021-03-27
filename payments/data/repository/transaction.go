package repository

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type TransactionRepository interface {
	Register(transaction *entity.Transaction) error
	Save(transaction *entity.Transaction) error
	Find(id string) (*entity.Transaction, error)
	FindOneByAccount(transactionId string, accountId string) (*entity.Transaction, error)
	FindOneByService(transactionId string, serviceId string) (*entity.Transaction, error)
	FindOneByStore(transactionId string, storeId string) (*entity.Transaction, error)
	FindAll(pagination *entity.Pagination) ([]*entity.Transaction, int, error)
	FindByAccountFromId(accountId string, pagination *entity.Pagination) ([]*entity.Transaction, int, error)
	FindByAccountToId(accountId string, pagination *entity.Pagination) ([]*entity.Transaction, int, error)
	FindByServiceId(serviceId string, pagination *entity.Pagination) ([]*entity.Transaction, int, error)
	FindByStoreId(storeId string, pagination *entity.Pagination) ([]*entity.Transaction, int, error)
}
