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

func TestFind(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate", func(t *testing.T) {
		is := require.New(t)

		id := uuid.NewV1().String()

		c := controller.NewTransaction(nil)

		result, err := c.Find(context.TODO(), id)

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

		result, err := c.Find(context.TODO(), id)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "an error on get payment")
	})

	t.Run("should fail on find if returns nil usecase transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		id := uuid.NewV4().String()
		transactionUseCase.On("Find", id).Return(nil, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Find(context.TODO(), id)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "no payment was found")
	})

	t.Run("should succeed on find usecase transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		accountTo, _ := entity.NewAccount(10)

		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 19, "AKZ")

		transactionUseCase.On("Find", transaction.ID).Return(transaction, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Find(context.TODO(), transaction.ID)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}

func TestFindAll(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find all on validate", func(t *testing.T) {
		is := require.New(t)

		page := -11
		limit := 10
		sort := "created_at DESC"

		c := controller.NewTransaction(nil)

		result, total, err := c.FindAll(context.TODO(), page, limit, sort)

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

		result, total, err := c.FindAll(context.TODO(), page, limit, sort)

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

		result, total, err := c.FindAll(context.TODO(), page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.NotNil(err)
		is.EqualError(err, "no payment was found")
	})

	t.Run("should succeed on find all usecase transactions", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockTransactionUseCase()

		page := 1
		limit := 10
		sort := "created_at DESC"

		accountFrom, _ := entity.NewAccount(200)
		accountTo, _ := entity.NewAccount(10)

		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 19, "AKZ")

		transactionUseCase.On("FindAll", page, limit, sort).Return([]*entity.Transaction{transaction}, 1, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, total, err := c.FindAll(context.TODO(), page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result[0], transaction)
		is.Equal(1, total)
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

		accountFrom, _ := entity.NewAccount(200)
		accountTo, _ := entity.NewAccount(10)

		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 19, "AKZ")

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

		accountFrom, _ := entity.NewAccount(200)
		accountTo, _ := entity.NewAccount(10)

		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 19, "AKZ")

		transaction.Status = entity.TransactionError

		transactionUseCase.On("Error", transaction.ID).Return(transaction, nil)

		c := controller.NewTransaction(transactionUseCase)

		result, err := c.Error(context.TODO(), transaction.ID)
		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}
