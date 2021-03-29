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

func TestRegisterStoreTransaction(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate register store transaction", func(t *testing.T) {
		is := require.New(t)

		accountFrom := uuid.NewV4().String()
		storeId := uuid.NewV4().String()
		amount := 30.00
		currency := "AKZ"

		c := controller.NewStoreTransaction(nil)

		result, err := c.RegisterStoreTransaction(context.TODO(), accountFrom, storeId, amount, currency)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on register store transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockStoreTransactionUseCase()

		accountFrom := uuid.NewV4().String()
		storeId := uuid.NewV4().String()
		amount := 30.00
		currency := "AOA"

		transactionUseCase.On("RegisterStoreTransaction", accountFrom, storeId, amount, currency).Return(nil, errors.New("internal error"))

		c := controller.NewStoreTransaction(transactionUseCase)

		result, err := c.RegisterStoreTransaction(context.TODO(), accountFrom, storeId, amount, currency)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "error on register store payment")
	})

	t.Run("should succeed on register store transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockStoreTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store", "store description", uuid.NewV4().String(), uuid.NewV4().String())

		transaction, _ := entity.NewTransaction(accountFrom, nil, nil, store, 19, "AOA")

		transactionUseCase.On("RegisterStoreTransaction", accountFrom.ID, store.ID, transaction.Amount, transaction.Currency).Return(transaction, nil)

		c := controller.NewStoreTransaction(transactionUseCase)

		result, err := c.RegisterStoreTransaction(context.TODO(), transaction.AccountFromID, transaction.StoreID, transaction.Amount, transaction.Currency)

		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}

func TestFindAllByStoreId(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate list transaction by store destination", func(t *testing.T) {
		is := require.New(t)

		storeId := uuid.NewV1().String()
		page := 1
		limit := -10
		sort := "created_at DESC"

		c := controller.NewStoreTransaction(nil)

		result, total, err := c.FindAllByStoreId(context.TODO(), storeId, page, limit, sort)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.Equal(0, total)
	})

	t.Run("should fail on list transaction by store destination", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockStoreTransactionUseCase()

		storeId := uuid.NewV4().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		transactionUseCase.On("FindAllByStoreId", storeId, page, limit, sort).Return(nil, 0, errors.New("internal error"))

		c := controller.NewStoreTransaction(transactionUseCase)

		result, total, err := c.FindAllByStoreId(context.TODO(), storeId, page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "error on list payments by store destination")
		is.Equal(0, total)
	})

	t.Run("should fail on list transaction by store destination", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockStoreTransactionUseCase()

		storeId := uuid.NewV4().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		transactionUseCase.On("FindAllByStoreId", storeId, page, limit, sort).Return(nil, 0, nil)

		c := controller.NewStoreTransaction(transactionUseCase)

		result, total, err := c.FindAllByStoreId(context.TODO(), storeId, page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "no store payment was found")
		is.Equal(0, total)
	})

	t.Run("should succeed on list transaction by store destination", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockStoreTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store", "store description", uuid.NewV4().String(), uuid.NewV4().String())

		transaction, _ := entity.NewTransaction(accountFrom, nil, nil, store, 19, "AKZ")

		page := 1
		limit := 10
		sort := "created_at DESC"

		transactions := []*entity.Transaction{transaction}

		transactionUseCase.On("FindAllByStoreId", store.ID, page, limit, sort).Return(transactions, len(transactions), nil)

		c := controller.NewStoreTransaction(transactionUseCase)

		result, total, err := c.FindAllByStoreId(context.TODO(), store.ID, page, limit, sort)

		transactionUseCase.AssertExpectations(t)
		is.Nil(err)
		is.NotEmpty(result)
		is.Equal(len(result), total)
		is.Equal(result[0], transaction)
	})
}

func TestFindOneByStore(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate get transaction by store", func(t *testing.T) {
		is := require.New(t)

		transactionId := uuid.NewV1().String()
		storeId := uuid.NewV1().String()

		c := controller.NewStoreTransaction(nil)

		result, err := c.FindOneByStore(context.TODO(), storeId, transactionId)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on get transaction by store", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockStoreTransactionUseCase()

		transactionId := uuid.NewV4().String()
		storeId := uuid.NewV4().String()

		transactionUseCase.On("FindOneByStore", storeId, transactionId).Return(nil, errors.New("internal error"))

		c := controller.NewStoreTransaction(transactionUseCase)

		result, err := c.FindOneByStore(context.TODO(), storeId, transactionId)

		is.Nil(result)
		is.EqualError(err, "error on get payment by store")
	})

	t.Run("should fail on get transaction by store", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockStoreTransactionUseCase()

		transactionId := uuid.NewV4().String()
		storeId := uuid.NewV4().String()

		transactionUseCase.On("FindOneByStore", storeId, transactionId).Return(nil, nil)

		c := controller.NewStoreTransaction(transactionUseCase)

		result, err := c.FindOneByStore(context.TODO(), storeId, transactionId)

		is.Nil(result)
		is.EqualError(err, "no store payment was found")
	})

	t.Run("should succeed on get transaction by store", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockStoreTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store", "store description", uuid.NewV4().String(), uuid.NewV4().String())
		transaction, _ := entity.NewTransaction(accountFrom, nil, nil, store, 19, "AKZ")

		transactionUseCase.On("FindOneByStore", store.ID, transaction.ID).Return(transaction, nil)

		c := controller.NewStoreTransaction(transactionUseCase)

		result, err := c.FindOneByStore(context.TODO(), store.ID, transaction.ID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result, transaction)
	})
}
