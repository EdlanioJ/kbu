package controller

import (
	"context"
	"errors"

	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/EdlanioJ/kbu/payments/domain/usecase"
	"github.com/EdlanioJ/kbu/payments/presentation/utils/log"
)

var (
	errOnRegisterStoreTransaction  = errors.New("error on register store payment")
	errOnNotFoundStoreTransaction  = errors.New("no store payment was found")
	errOnFindAllTransactionByStore = errors.New("error on list payments by store destination")
	errOnFindOneTransactionByStore = errors.New("error on get payment by store")
)

type StoreTransaction struct {
	Transaction usecase.StoreTransaction
}

func NewStoreTransaction(Transaction usecase.StoreTransaction) *StoreTransaction {
	return &StoreTransaction{
		Transaction: Transaction,
	}
}

func (c *StoreTransaction) RegisterStoreTransaction(ctx context.Context, fromId string, storeId string, amount float64, currency string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.RegisterStoreTransaction(fromId, storeId, amount, currency)

	if err != nil {
		log.Error(ctx, err, errOnRegisterStoreTransaction.Error())

		return nil, errOnRegisterStoreTransaction
	}
	return transaction, nil
}

func (c *StoreTransaction) FindAllByStoreId(ctx context.Context, storeId string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	transactions, total, err := c.Transaction.FindAllByStoreId(storeId, page, limit, sort)

	if err != nil {
		log.Error(ctx, err, errOnFindAllTransactionByStore.Error())

		return nil, 0, errOnFindAllTransactionByStore
	}

	if len(transactions) == 0 {
		return nil, 0, errOnNotFoundStoreTransaction
	}

	return transactions, total, nil
}
func (c *StoreTransaction) FindOneByStore(ctx context.Context, storeId string, transactionId string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.FindOneByStore(storeId, transactionId)

	if err != nil {
		log.Error(ctx, err, errOnFindOneTransactionByStore.Error())

		return nil, errOnFindOneTransactionByStore
	}

	if transaction == nil {

		return nil, errOnNotFoundStoreTransaction
	}

	return transaction, nil
}
