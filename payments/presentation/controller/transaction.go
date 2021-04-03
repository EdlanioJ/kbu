package controller

import (
	"context"
	"errors"

	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/EdlanioJ/kbu/payments/domain/usecase"
	"github.com/EdlanioJ/kbu/payments/presentation/validator"
	log "github.com/sirupsen/logrus"
)

var (
	errOnRegister            = errors.New("an error on register payment")
	errOnFindAllTransaction  = errors.New("an error on list payments")
	errOnNotFoundTransaction = errors.New("no payment was found")
	errOnListByType          = errors.New("an error on list payments by type")
	errOnListByExternalID    = errors.New("an error on list payments by reference")
	errOnListByAccountFrom   = errors.New("an error on list payments by account from")
	errOnListByAccountTo     = errors.New("an error on list payments by account destination")
	errOnCompeteTransaction  = errors.New("error on complete payment")
	errOnCancelTransaction   = errors.New("error on cancel payment")
)

type Transaction struct {
	Transaction usecase.Transaction
	logger      *log.Logger
}

func NewTransaction(transaction usecase.Transaction) *Transaction {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})

	return &Transaction{
		Transaction: transaction,
		logger:      logger,
	}
}

func (c *Transaction) Register(ctx context.Context, accountFrom, accountTo, externalID, transactionType, currency string, amount float64) (*entity.Transaction, error) {
	err := validator.RegisterParams(accountFrom, accountTo, externalID, transactionType, currency, amount)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, err
	}
	transaction, err := c.Transaction.Register(accountFrom, accountTo, externalID, transactionType, currency, amount)

	if err != nil {
		c.logger.
			WithFields(log.Fields{
				"from_account_id": accountFrom,
				"to_account_id":   accountTo,
				"reference_id":    externalID,
				"type":            transactionType,
				"currency":        currency,
				"amount":          amount,
			}).WithContext(ctx).
			WithError(err).
			Error(errOnRegister)

		return nil, errOnRegister
	}

	return transaction, nil
}

func (c *Transaction) Get(ctx context.Context, id string) (*entity.Transaction, error) {
	err := validator.GetParams(id)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, err
	}

	transaction, err := c.Transaction.Find(id)

	if err != nil {
		c.logger.
			WithField("transaction_id", id).
			WithContext(ctx).
			WithError(err).
			Error(errOnNotFoundTransaction)
		return nil, errOnNotFoundTransaction
	}

	return transaction, nil
}

func (c *Transaction) List(ctx context.Context, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	err := validator.GetAllParams(page, limit, sort)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}
	transactions, total, err := c.Transaction.FindAll(page, limit, sort)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"page":  page,
					"limit": limit,
					"sort":  sort,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnFindAllTransaction)
		return nil, 0, errOnFindAllTransaction
	}

	if len(transactions) == 0 {
		c.logger.WithContext(ctx).WithFields(
			log.Fields{
				"page":  page,
				"limit": limit,
				"sort":  sort,
			},
		).Info("no payment was found")
		return nil, 0, nil
	}

	return transactions, total, nil
}

func (c *Transaction) GetByType(ctx context.Context, transactionID, transactionType string) (*entity.Transaction, error) {
	err := validator.GetByTypeParams(transactionID, transactionType)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, err
	}

	transaction, err := c.Transaction.FindByType(transactionType, transactionID)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"type":           transactionType,
					"transaction_id": transactionID,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnNotFoundTransaction)
		return nil, errOnNotFoundTransaction
	}

	return transaction, nil
}

func (c *Transaction) ListByType(ctx context.Context, transactionType string, page, limit int, sort string) ([]*entity.Transaction, int, error) {
	err := validator.ListByTypeParams(transactionType, page, limit, sort)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}

	transactions, total, err := c.Transaction.FindAllByType(transactionType, page, limit, sort)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"type":  transactionType,
					"page":  page,
					"limit": limit,
					"sort":  sort,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnListByType)
		return nil, 0, errOnListByType
	}

	if len(transactions) == 0 {
		c.logger.WithContext(ctx).WithFields(
			log.Fields{
				"type":  transactionType,
				"page":  page,
				"limit": limit,
				"sort":  sort,
			},
		).Warn("no payment was found")
		return nil, 0, nil
	}

	return transactions, total, nil
}

func (c *Transaction) GetByExternalID(ctx context.Context, transactionID, externalID string) (*entity.Transaction, error) {
	err := validator.GetByExternalIDParams(transactionID, externalID)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, err
	}

	transaction, err := c.Transaction.FindByExternalID(externalID, transactionID)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"reference_id":   externalID,
					"transaction_id": transactionID,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnNotFoundTransaction)
		return nil, errOnNotFoundTransaction
	}

	return transaction, nil
}

func (c *Transaction) ListByExternalID(ctx context.Context, externalID string, page, limit int, sort string) ([]*entity.Transaction, int, error) {
	err := validator.ListByExternalIDParams(externalID, page, limit, sort)
	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}

	transactions, total, err := c.Transaction.FindAllByExternalID(externalID, page, limit, sort)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"reference_id": externalID,
					"page":         page,
					"limit":        limit,
					"sort":         sort,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnListByExternalID)
		return nil, 0, errOnListByExternalID
	}

	if len(transactions) == 0 {
		c.logger.
			WithFields(
				log.Fields{
					"reference_id": externalID,
					"page":         page,
					"limit":        limit,
					"sort":         sort,
				},
			).WithContext(ctx).
			Warn("no payment was found")
		return nil, 0, nil
	}

	return transactions, total, nil
}

func (c *Transaction) GetByAccountFrom(ctx context.Context, transactionID, accountID string) (*entity.Transaction, error) {
	err := validator.GetByAccountFromParams(transactionID, accountID)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, err
	}

	transaction, err := c.Transaction.FindByFromAccountID(accountID, transactionID)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"from_account_id": accountID,
					"transaction_id":  transactionID,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnNotFoundTransaction)
		return nil, errOnNotFoundTransaction
	}

	return transaction, nil
}

func (c *Transaction) ListByAccountFrom(ctx context.Context, accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	err := validator.ListByAccountFromParams(accountID, page, limit, sort)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}

	transactions, total, err := c.Transaction.FindAllByFromAccountID(accountID, page, limit, sort)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"from_account_id": accountID,
					"page":            page,
					"limit":           limit,
					"sort":            sort,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnListByAccountFrom)
		return nil, 0, errOnListByAccountFrom
	}

	if len(transactions) == 0 {
		c.logger.
			WithFields(
				log.Fields{
					"from_account_id": accountID,
					"page":            page,
					"limit":           limit,
					"sort":            sort,
				},
			).WithContext(ctx).
			Warn("no payment was found")
		return nil, 0, nil
	}

	return transactions, total, nil
}

func (c *Transaction) GetByAccoutTo(ctx context.Context, transactionID, accountID string) (*entity.Transaction, error) {
	err := validator.GetByAccoutToParams(transactionID, accountID)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, err
	}

	transaction, err := c.Transaction.FindByToAccountID(accountID, transactionID)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"to_account_id":  accountID,
					"transaction_id": transactionID,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnNotFoundTransaction)
		return nil, errOnNotFoundTransaction
	}

	return transaction, nil
}

func (c *Transaction) ListByAccountTo(ctx context.Context, accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	err := validator.ListByAccoutToParams(accountID, page, limit, sort)
	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}

	transactions, total, err := c.Transaction.FindAllByToAccountID(accountID, page, limit, sort)

	if err != nil {
		c.logger.
			WithFields(
				log.Fields{
					"destination_id": accountID,
					"page":           page,
					"limit":          limit,
					"sort":           sort,
				},
			).WithContext(ctx).
			WithError(err).
			Error(errOnListByAccountTo)
		return nil, 0, errOnListByAccountTo
	}

	if len(transactions) == 0 {
		c.logger.
			WithFields(
				log.Fields{
					"destination_id": accountID,
					"page":           page,
					"limit":          limit,
					"sort":           sort,
				},
			).WithContext(ctx).
			Warn("no payment was found")
		return nil, 0, nil
	}

	return transactions, total, nil
}

func (c *Transaction) Complete(ctx context.Context, transactionId string) (*entity.Transaction, error) {
	err := validator.CompleteParams(transactionId)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, err
	}

	transaction, err := c.Transaction.Complete(transactionId)

	if err != nil {
		c.logger.
			WithField("transaction_id", transactionId).
			WithContext(ctx).WithError(err).
			Error(errOnCompeteTransaction)
		return nil, errOnCompeteTransaction
	}

	return transaction, nil
}

func (c *Transaction) Error(ctx context.Context, transactionId string) (*entity.Transaction, error) {
	err := validator.ErrorParams(transactionId)

	if err != nil {
		c.logger.WithContext(ctx).Error(err)
		return nil, err
	}

	transaction, err := c.Transaction.Error(transactionId)

	if err != nil {
		c.logger.
			WithContext(ctx).
			WithField("transaction_id", transactionId).
			WithError(err).
			Error(errOnCancelTransaction)
		return nil, errOnCancelTransaction
	}

	return transaction, nil
}
