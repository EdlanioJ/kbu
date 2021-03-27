package usecase

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type StoreTransaction interface {
	RegisterStoreTransaction(fromAccountId string, storeId string, amount float64, currency string) (*entity.Transaction, error)
	FindAllByStoreId(storeId string, page int, limit int, sort string) ([]*entity.Transaction, int, error)
	FindOneByStore(storeId string, transactionId string) (*entity.Transaction, error)
}
