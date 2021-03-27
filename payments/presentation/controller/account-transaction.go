package controller

import (
	"context"
	"errors"

	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/EdlanioJ/kbu/payments/domain/usecase"
	"github.com/EdlanioJ/kbu/payments/presentation/utils/log"
)

var (
	errOnRegisterAccountTransaction  = errors.New("error on register account payment")
	errOnNotFoundAccountTransaction  = errors.New("no account payment was found")
	errOnFindAllTransactionByAccount = errors.New("error on list payments by account destination")
	errOnFindOneTransactionByAccount = errors.New("error on get payment by account from")
)

type AccountTransaction struct {
	Transaction usecase.AccountTransaction
}

func NewAccountTransaction(Transaction usecase.AccountTransaction) *AccountTransaction {
	return &AccountTransaction{
		Transaction: Transaction,
	}
}

func (c *AccountTransaction) RegisterAccountTransaction(ctx context.Context, fromId string, toId string, amount float64, currency string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.RegisterAccountTransaction(fromId, toId, amount, currency)

	if err != nil {
		log.Error(ctx, err, errOnRegisterAccountTransaction.Error())

		return nil, errOnRegisterAccountTransaction
	}
	return transaction, nil
}

func (c *AccountTransaction) FindAllByAccountTo(ctx context.Context, accountId string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	transactions, total, err := c.Transaction.FindAllByAccountTo(accountId, page, limit, sort)

	if err != nil {
		log.Error(ctx, err, errOnFindAllTransactionByAccount.Error())

		return nil, 0, errOnFindAllTransactionByAccount
	}

	if len(transactions) == 0 {
		return nil, 0, errOnNotFoundAccountTransaction
	}

	return transactions, total, nil
}

func (c *AccountTransaction) FindOneByAccount(ctx context.Context, fromId string, transactionId string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.FindOneByAccount(fromId, transactionId)

	if err != nil {
		log.Error(ctx, err, errOnFindOneTransactionByAccount.Error())

		return nil, errOnFindOneTransactionByAccount
	}

	if transaction == nil {
		log.Error(ctx, err, errOnNotFoundAccountTransaction.Error())

		return nil, errOnNotFoundAccountTransaction
	}
	return transaction, nil
}
