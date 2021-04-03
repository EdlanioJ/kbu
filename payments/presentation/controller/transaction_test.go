package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/EdlanioJ/kbu/payments/presentation/controller"
	"github.com/EdlanioJ/kbu/payments/presentation/controller/mock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)
		accountFrom := uuid.NewV1().String()
		accountTo := uuid.NewV4().String()
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AKZ"
		amount := 30.00

		c := controller.NewTransaction(nil)

		result, err := c.Register(context.TODO(), accountFrom, accountTo, externalID, transactionType, currency, amount)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on register", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom := uuid.NewV4().String()
		accountTo := uuid.NewV4().String()
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00

		transactionUseCase.On("Register", accountFrom, accountTo, externalID, transactionType, currency, amount).Return(nil, errors.New("register error"))
		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Register(context.TODO(), accountFrom, accountTo, externalID, transactionType, currency, amount)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "an error on register payment")
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactionUseCase.On("Register", accountFrom.ID, accountTo.ID, externalID, transactionType, currency, amount).Return(transaction, nil)
		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Register(context.TODO(), accountFrom.ID, accountTo.ID, externalID, transactionType, currency, amount)

		is.NotNil(result)
		is.Equal(result, transaction)
		is.Nil(err)
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate", func(t *testing.T) {
		is := require.New(t)

		id := uuid.NewV1().String()

		c := controller.NewTransaction(nil)

		result, err := c.Get(context.TODO(), id)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find usecase transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		id := uuid.NewV4().String()
		transactionUseCase.On("Find", id).Return(nil, errors.New("usecase error"))

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Get(context.TODO(), id)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "no payment was found")
	})

	t.Run("should succeed on find usecase transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactionUseCase.On("Find", transaction.ID).Return(transaction, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Get(context.TODO(), transaction.ID)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find all on validate", func(t *testing.T) {
		is := require.New(t)

		page := -11
		limit := 10
		sort := "created_at DESC"

		c := controller.NewTransaction(nil)

		result, total, err := c.List(context.TODO(), page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find all usecase transactions", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		page := 1
		limit := 10
		sort := "created_at DESC"
		transactionUseCase.On("FindAll", page, limit, sort).Return(nil, 0, errors.New("usecase error"))

		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.List(context.TODO(), page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.NotNil(err)
		is.EqualError(err, "an error on list payments")
	})

	t.Run("should fail on find all if returns nil usecase transactions", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		page := 1
		limit := 10
		sort := "created_at DESC"
		transactionUseCase.On("FindAll", page, limit, sort).Return(nil, 0, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.List(context.TODO(), page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.Nil(err)
	})

	t.Run("should succeed on find all usecase transactions", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		page := 1
		limit := 10
		sort := "created_at DESC"

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactionUseCase.On("FindAll", page, limit, sort).Return([]*entity.Transaction{transaction}, 1, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.List(context.TODO(), page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result[0], transaction)
		is.Equal(1, total)
	})
}

func TestGetByType(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)

		transactionID := uuid.NewV1().String()
		transactionType := "invalid type"

		c := controller.NewTransaction(nil)

		result, err := c.GetByType(context.TODO(), transactionID, transactionType)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find by type", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		transactionID := uuid.NewV4().String()
		transactionType := entity.TransactionToService

		transactionUseCase.On("FindByType", transactionType, transactionID).Return(nil, errors.New("error on get"))
		c := controller.NewTransaction(transactionUseCase)

		result, err := c.GetByType(context.TODO(), transactionID, transactionType)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "no payment was found")
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactionUseCase.On("FindByType", transactionType, transaction.ID).Return(transaction, nil)
		c := controller.NewTransaction(transactionUseCase)

		result, err := c.GetByType(context.TODO(), transaction.ID, transactionType)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result, transaction)
	})
}

func TestListByType(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)

		transactionType := "invalid type"
		page := 1
		limit := -10
		sort := "id Desc"
		c := controller.NewTransaction(nil)

		result, total, err := c.ListByType(context.TODO(), transactionType, page, limit, sort)

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		transactionType := entity.TransactionToService
		page := 1
		limit := 10
		sort := "created_at Desc"

		transactionUseCase.On("FindAllByType", transactionType, page, limit, sort).Return(nil, 0, errors.New("internal error"))
		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.ListByType(context.TODO(), transactionType, page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "an error on list payments by type")
	})

	t.Run("should return nil", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		transactionType := entity.TransactionToService
		page := 1
		limit := 10
		sort := "created_at Desc"

		transactionUseCase.On("FindAllByType", transactionType, page, limit, sort).Return(nil, 0, nil)
		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.ListByType(context.TODO(), transactionType, page, limit, sort)

		is.Nil(result)
		is.Equal(total, 0)
		is.Nil(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactions := []*entity.Transaction{transaction}
		page := 1
		limit := 10
		sort := "created_at Desc"

		transactionUseCase.On("FindAllByType", transactionType, page, limit, sort).Return(transactions, len(transactions), nil)
		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.ListByType(context.TODO(), transactionType, page, limit, sort)

		is.Equal(len(transactions), total)
		is.Nil(err)
		is.NotNil(result)
		is.Equal(result[0], transaction)

	})
}

func TestGetByExternalID(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)
		transactionID := uuid.NewV1().String()
		externalID := uuid.NewV1().String()

		c := controller.NewTransaction(nil)

		result, err := c.GetByExternalID(context.TODO(), transactionID, externalID)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		transactionID := uuid.NewV4().String()
		externalID := uuid.NewV4().String()
		transactionUseCase.On("FindByExternalID", externalID, transactionID).Return(nil, errors.New("internal error"))

		c := controller.NewTransaction(transactionUseCase)

		_, err := c.GetByExternalID(context.TODO(), transactionID, externalID)

		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "no payment was found")
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactionUseCase.On("FindByExternalID", externalID, transaction.ID).Return(transaction, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.GetByExternalID(context.TODO(), transaction.ID, externalID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result, transaction)

	})
}

func TestListByExternalID(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)

		externalID := uuid.NewV1().String()
		page := -2
		limit := -20
		sort := "id desc"

		c := controller.NewTransaction(nil)

		result, total, err := c.ListByExternalID(context.TODO(), externalID, page, limit, sort)

		is.Equal(total, 0)
		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		externalID := uuid.NewV4().String()
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByExternalID", externalID, page, limit, sort).Return(nil, 0, errors.New("internal error"))

		c := controller.NewTransaction(transactionUseCase)

		result, _, err := c.ListByExternalID(context.TODO(), externalID, page, limit, sort)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "an error on list payments by reference")
	})

	t.Run("should return nil", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		externalID := uuid.NewV4().String()
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByExternalID", externalID, page, limit, sort).Return(nil, 0, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, _, err := c.ListByExternalID(context.TODO(), externalID, page, limit, sort)

		is.Nil(result)
		is.Nil(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()
		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactions := []*entity.Transaction{transaction}
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByExternalID", externalID, page, limit, sort).Return(transactions, len(transactions), nil)

		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.ListByExternalID(context.TODO(), externalID, page, limit, sort)

		is.Nil(err)
		is.Equal(len(transactions), total)
		is.NotNil(result)
		is.Equal(result[0], transaction)
	})
}

func TestGetByAccountFrom(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)
		transactionID := uuid.NewV1().String()
		accountID := uuid.NewV1().String()

		c := controller.NewTransaction(nil)

		result, err := c.GetByAccountFrom(context.TODO(), transactionID, accountID)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		transactionID := uuid.NewV4().String()
		accountID := uuid.NewV4().String()
		transactionUseCase.On("FindByFromAccountID", accountID, transactionID).Return(nil, errors.New("internal error"))

		c := controller.NewTransaction(transactionUseCase)

		_, err := c.GetByAccountFrom(context.TODO(), transactionID, accountID)

		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "no payment was found")
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactionUseCase.On("FindByFromAccountID", accountFrom.ID, transaction.ID).Return(transaction, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.GetByAccountFrom(context.TODO(), transaction.ID, accountFrom.ID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result, transaction)

	})
}

func TestListByAccountFrom(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)

		accountID := uuid.NewV1().String()
		page := -2
		limit := -20
		sort := "id desc"

		c := controller.NewTransaction(nil)

		result, total, err := c.ListByAccountFrom(context.TODO(), accountID, page, limit, sort)

		is.Equal(total, 0)
		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountID := uuid.NewV4().String()
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByFromAccountID", accountID, page, limit, sort).Return(nil, 0, errors.New("internal error"))

		c := controller.NewTransaction(transactionUseCase)

		result, _, err := c.ListByAccountFrom(context.TODO(), accountID, page, limit, sort)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "an error on list payments by account from")
	})

	t.Run("should return nil", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountID := uuid.NewV4().String()
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByFromAccountID", accountID, page, limit, sort).Return(nil, 0, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, _, err := c.ListByAccountFrom(context.TODO(), accountID, page, limit, sort)

		is.Nil(result)
		is.Nil(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()
		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactions := []*entity.Transaction{transaction}
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByFromAccountID", accountFrom.ID, page, limit, sort).Return(transactions, len(transactions), nil)

		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.ListByAccountFrom(context.TODO(), accountFrom.ID, page, limit, sort)

		is.Nil(err)
		is.Equal(len(transactions), total)
		is.NotNil(result)
		is.Equal(result[0], transaction)
	})
}

func TestGetByAccoutTo(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)
		transactionID := uuid.NewV1().String()
		accountID := uuid.NewV1().String()

		c := controller.NewTransaction(nil)

		result, err := c.GetByAccoutTo(context.TODO(), transactionID, accountID)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		transactionID := uuid.NewV4().String()
		accountID := uuid.NewV4().String()
		transactionUseCase.On("FindByToAccountID", accountID, transactionID).Return(nil, errors.New("internal error"))

		c := controller.NewTransaction(transactionUseCase)

		_, err := c.GetByAccoutTo(context.TODO(), transactionID, accountID)

		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "no payment was found")
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactionUseCase.On("FindByToAccountID", accountTo.ID, transaction.ID).Return(transaction, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.GetByAccoutTo(context.TODO(), transaction.ID, accountTo.ID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result, transaction)

	})
}

func TestListByAccountTo(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validation", func(t *testing.T) {
		is := require.New(t)

		accountID := uuid.NewV1().String()
		page := -2
		limit := -20
		sort := "id desc"

		c := controller.NewTransaction(nil)

		result, total, err := c.ListByAccountTo(context.TODO(), accountID, page, limit, sort)

		is.Equal(total, 0)
		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on find", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountID := uuid.NewV4().String()
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByToAccountID", accountID, page, limit, sort).Return(nil, 0, errors.New("internal error"))

		c := controller.NewTransaction(transactionUseCase)

		result, _, err := c.ListByAccountTo(context.TODO(), accountID, page, limit, sort)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "an error on list payments by account destination")
	})

	t.Run("should return nil", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountID := uuid.NewV4().String()
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByToAccountID", accountID, page, limit, sort).Return(nil, 0, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, _, err := c.ListByAccountTo(context.TODO(), accountID, page, limit, sort)

		is.Nil(result)
		is.Nil(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()
		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transactions := []*entity.Transaction{transaction}
		page := 1
		limit := 20
		sort := "created_at desc"

		transactionUseCase.On("FindAllByToAccountID", accountTo.ID, page, limit, sort).Return(transactions, len(transactions), nil)

		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.ListByAccountTo(context.TODO(), accountTo.ID, page, limit, sort)

		is.Nil(err)
		is.Equal(len(transactions), total)
		is.NotNil(result)
		is.Equal(result[0], transaction)
	})
}

func TestComplete(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate", func(t *testing.T) {
		is := require.New(t)

		id := uuid.NewV1().String()

		c := controller.NewTransaction(nil)

		result, err := c.Complete(context.TODO(), id)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on complete usecase", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		id := uuid.NewV4().String()
		transactionUseCase.On("Complete", id).Return(nil, errors.New("usecase error"))

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Complete(context.TODO(), id)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "error on complete payment")
	})

	t.Run("should succeed on complete usecase", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transaction.Status = entity.TransactionCompleted

		transactionUseCase.On("Complete", transaction.ID).Return(transaction, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Complete(context.TODO(), transaction.ID)
		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}
func TestError(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate", func(t *testing.T) {
		is := require.New(t)

		id := uuid.NewV1().String()

		c := controller.NewTransaction(nil)

		result, err := c.Error(context.TODO(), id)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on cancel", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		id := uuid.NewV4().String()
		transactionUseCase.On("Error", id).Return(nil, errors.New("usecase error"))

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Error(context.TODO(), id)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "error on cancel payment")
	})

	t.Run("should succeed on cancel", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		transaction.Status = entity.TransactionCanceled

		transactionUseCase.On("Error", transaction.ID).Return(transaction, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Error(context.TODO(), transaction.ID)
		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}
