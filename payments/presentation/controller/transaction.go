package controller

import (
	"context"
	"errors"

	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/EdlanioJ/kbu/payments/domain/usecase"
	"github.com/EdlanioJ/kbu/payments/presentation/utils/log"
)

var (
	errOnFindTransaction     = errors.New("an error on get payment")
	errOnFindAllTransaction  = errors.New("an error on list payments")
	errOnNotFoundTransaction = errors.New("no payment was found")
	errOnCompeteTransaction  = errors.New("error on complete payment")
	errOnCancelTransaction   = errors.New("error on cancel payment")
)

type Transaction struct {
	Transaction usecase.Transaction
}

func NewTransaction(transaction usecase.Transaction) *Transaction {
	return &Transaction{
		Transaction: transaction,
	}
}

func (c *Transaction) Find(ctx context.Context, id string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.Find(id)

	if err != nil {
		log.Error(ctx, err, errOnFindTransaction.Error())

		return nil, errOnFindTransaction
	}

	if transaction == nil {
		log.Info(ctx, errOnNotFoundTransaction.Error())

		return nil, errOnNotFoundTransaction
	}

	return transaction, nil
}

func (c *Transaction) FindAll(ctx context.Context, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	transactions, total, err := c.Transaction.FindAll(page, limit, sort)

	if err != nil {
		log.Error(ctx, err, errOnFindAllTransaction.Error())

		return nil, 0, errOnFindAllTransaction
	}

	if len(transactions) == 0 {
		log.Info(ctx, errOnNotFoundTransaction.Error())

		return nil, 0, errOnNotFoundTransaction
	}

	return transactions, total, nil
}

func (c *Transaction) Complete(ctx context.Context, transactionId string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.Complete(transactionId)

	if err != nil {
		log.Error(ctx, err, errOnCompeteTransaction.Error())

		return nil, errOnCompeteTransaction
	}

	return transaction, nil
}

func (c *Transaction) Error(ctx context.Context, transactionId string) (*entity.Transaction, error) {
	transaction, err := c.Transaction.Error(transactionId)

	if err != nil {
		log.Error(ctx, err, errOnCancelTransaction.Error())

		return nil, errOnCancelTransaction
	}

	return transaction, nil
}
