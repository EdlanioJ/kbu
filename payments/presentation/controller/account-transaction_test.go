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

func TestRegisterAccountTransaction(t *testing.T) {
	t.Parallel()

	t.Run("should fail on register account transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockAccountTransactionUseCase()

		accountFrom := uuid.NewV4().String()
		accountTo := uuid.NewV4().String()
		amount := 30.00
		currency := "AOA"

		transactionUseCase.On("RegisterAccountTransaction", accountFrom, accountTo, amount, currency).Return(nil, errors.New("internal error"))

		c := controller.NewAccountTransaction(transactionUseCase)

		result, err := c.RegisterAccountTransaction(context.TODO(), accountFrom, accountTo, amount, currency)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "error on register account payment")
	})

	t.Run("should fail on validate params register account transaction", func(t *testing.T) {
		is := require.New(t)

		accountFrom := uuid.NewV4().String()
		accountTo := "invalid id"
		amount := 30.00
		currency := "AOA"

		c := controller.NewAccountTransaction(nil)

		result, err := c.RegisterAccountTransaction(context.TODO(), accountFrom, accountTo, amount, currency)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed on register account transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockAccountTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		accountTo, _ := entity.NewAccount(10)

		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 19, "AOA")

		transactionUseCase.On("RegisterAccountTransaction", accountFrom.ID, accountTo.ID, transaction.Amount, transaction.Currency).Return(transaction, nil)

		c := controller.NewAccountTransaction(transactionUseCase)

		result, err := c.RegisterAccountTransaction(context.TODO(), transaction.AccountFromID, transaction.AccountToID, transaction.Amount, transaction.Currency)

		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}

func TestFindAllByAccountTo(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate list transaction by account destination", func(t *testing.T) {
		is := require.New(t)

		accountTo := "invalid id"
		page := -4
		limit := 10
		sort := "created_at DESC"

		c := controller.NewAccountTransaction(nil)

		result, total, err := c.FindAllByAccountTo(context.TODO(), accountTo, page, limit, sort)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.Equal(0, total)
	})

	t.Run("should fail on list transaction by account destination", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockAccountTransactionUseCase()

		accountTo := uuid.NewV4().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		transactionUseCase.On("FindAllByAccountTo", accountTo, page, limit, sort).Return(nil, 0, errors.New("internal error"))

		c := controller.NewAccountTransaction(transactionUseCase)

		result, total, err := c.FindAllByAccountTo(context.TODO(), accountTo, page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "error on list payments by account destination")
		is.Equal(0, total)
	})

	t.Run("should fail on list transaction by account destination if returns nil", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockAccountTransactionUseCase()

		accountTo := uuid.NewV4().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		transactionUseCase.On("FindAllByAccountTo", accountTo, page, limit, sort).Return(nil, 0, nil)

		c := controller.NewAccountTransaction(transactionUseCase)

		result, total, err := c.FindAllByAccountTo(context.TODO(), accountTo, page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "no account payment was found")
		is.Equal(0, total)
	})

	t.Run("should succeed on list transaction by account destination", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockAccountTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		accountTo, _ := entity.NewAccount(10)

		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 19, "AOA")

		page := 1
		limit := 10
		sort := "created_at DESC"

		transactionUseCase.On("FindAllByAccountTo", accountTo.ID, page, limit, sort).Return([]*entity.Transaction{transaction}, 1, nil)

		c := controller.NewAccountTransaction(transactionUseCase)

		result, total, err := c.FindAllByAccountTo(context.TODO(), accountTo.ID, page, limit, sort)

		is.Nil(err)
		is.NotEmpty(result)
		is.Equal(len(result), 1)
		is.Equal(1, total)
		is.Equal(result[0], transaction)
	})
}

func TestFindOneByAccount(t *testing.T) {
	t.Parallel()

	t.Run("should fail on get transaction by account from", func(t *testing.T) {
		is := require.New(t)

		transactionId := "invalid id"
		accountId := uuid.NewV1().String()

		c := controller.NewAccountTransaction(nil)

		result, err := c.FindOneByAccount(context.TODO(), accountId, transactionId)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on get transaction by account from", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockAccountTransactionUseCase()

		transactionId := uuid.NewV4().String()
		accountId := uuid.NewV4().String()

		transactionUseCase.On("FindOneByAccount", accountId, transactionId).Return(nil, errors.New("internal error"))

		c := controller.NewAccountTransaction(transactionUseCase)

		result, err := c.FindOneByAccount(context.TODO(), accountId, transactionId)

		is.Nil(result)
		is.EqualError(err, "error on get payment by account from")
	})

	t.Run("should fail on get transaction by account from if returns nil", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockAccountTransactionUseCase()

		transactionId := uuid.NewV4().String()
		accountId := uuid.NewV4().String()

		transactionUseCase.On("FindOneByAccount", accountId, transactionId).Return(nil, nil)

		c := controller.NewAccountTransaction(transactionUseCase)

		result, err := c.FindOneByAccount(context.TODO(), accountId, transactionId)

		is.Nil(result)
		is.EqualError(err, "no account payment was found")
	})

	t.Run("should succeed on get transaction by account from", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockAccountTransactionUseCase()
		accountFrom, _ := entity.NewAccount(200)
		accountTo, _ := entity.NewAccount(10)

		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 19, "AOA")

		transactionUseCase.On("FindOneByAccount", transaction.AccountFromID, transaction.ID).Return(transaction, nil)

		c := controller.NewAccountTransaction(transactionUseCase)

		result, err := c.FindOneByAccount(context.TODO(), transaction.AccountFromID, transaction.ID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result, transaction)
	})
}
