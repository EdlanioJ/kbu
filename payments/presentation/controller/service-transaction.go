package controller

import (
	"context"
	"errors"

	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/EdlanioJ/kbu/payments/domain/usecase"
	"github.com/EdlanioJ/kbu/payments/presentation/utils/log"
)

var (
	errOnRegisterServiceTransaction  = errors.New("error on register service payment")
	errOnNotFoundServiceTransaction  = errors.New("no service payment was found")
	errOnFindAllTransactionByService = errors.New("error on list payments by service destination")
	errOnFindOneTransactionByService = errors.New("error on get payment by service")
)

type ServiceTransaction struct {
	Transaction usecase.ServiceTransaction
}

func NewServiceTransaction(Transaction usecase.ServiceTransaction) *ServiceTransaction {
	return &ServiceTransaction{
		Transaction: Transaction,
	}
}

func (c *ServiceTransaction) RegisterServiceTransaction(ctx context.Context, fromId string, serviceId string, servicePriceId string, amount float64, currency string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.RegisterServiceTransaction(fromId, serviceId, servicePriceId, amount, currency)

	if err != nil {
		log.Error(ctx, err, errOnRegisterServiceTransaction.Error())
		return nil, errOnRegisterServiceTransaction
	}

	return transaction, nil
}

func (c *ServiceTransaction) FindAllByServiceId(ctx context.Context, serviceId string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	transactions, total, err := c.Transaction.FindAllByServiceId(serviceId, page, limit, sort)
	if err != nil {
		log.Error(ctx, err, errOnFindAllTransactionByService.Error())
		return nil, 0, errOnFindAllTransactionByService
	}

	if len(transactions) == 0 {
		return nil, 0, errOnNotFoundServiceTransaction
	}
	return transactions, total, nil
}

func (c *ServiceTransaction) FindOneByService(ctx context.Context, serviceId string, transactionId string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.FindOneByService(serviceId, transactionId)

	if err != nil {
		log.Error(ctx, err, errOnFindOneTransactionByService.Error())

		return nil, errOnFindOneTransactionByService
	}

	if transaction == nil {
		return nil, errOnNotFoundServiceTransaction
	}

	return transaction, nil
}
