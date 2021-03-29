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

func TestRegisterServiceTransaction(t *testing.T) {
	t.Parallel()

	t.Run("should fail on register service transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockServiceTransactionUseCase()

		accountFrom := uuid.NewV4().String()
		serviceTo := uuid.NewV4().String()
		amount := 30.00
		currency := "AOA"
		servicePrice := uuid.NewV4().String()

		transactionUseCase.On("RegisterServiceTransaction", accountFrom, serviceTo, servicePrice, amount, currency).Return(nil, errors.New("internal error"))

		c := controller.NewServiceTransaction(transactionUseCase)

		result, err := c.RegisterServiceTransaction(context.TODO(), accountFrom, serviceTo, servicePrice, amount, currency)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "error on register service payment")
	})

	t.Run("should fail on validate params register service transaction", func(t *testing.T) {
		is := require.New(t)

		accountFrom := "invalid id"
		serviceTo := uuid.NewV4().String()
		amount := 30.00
		currency := "AOA"
		servicePrice := "invalid id"

		c := controller.NewServiceTransaction(nil)

		result, err := c.RegisterServiceTransaction(context.TODO(), accountFrom, serviceTo, servicePrice, amount, currency)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed on register service transaction", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockServiceTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		service, _ := entity.NewService("service", "service description", uuid.NewV4().String(), uuid.NewV4().String())

		transaction, _ := entity.NewTransaction(accountFrom, nil, service, nil, 19, "AOA")

		servicePrice := uuid.NewV4().String()
		transactionUseCase.On("RegisterServiceTransaction", accountFrom.ID, service.ID, servicePrice, transaction.Amount, transaction.Currency).Return(transaction, nil)

		c := controller.NewServiceTransaction(transactionUseCase)

		result, err := c.RegisterServiceTransaction(context.TODO(), accountFrom.ID, service.ID, servicePrice, transaction.Amount, transaction.Currency)

		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}

func TestFindAllByServiceId(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate list transaction by service destination", func(t *testing.T) {
		is := require.New(t)

		serviceId := uuid.NewV1().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		c := controller.NewServiceTransaction(nil)

		result, total, err := c.FindAllByServiceId(context.TODO(), serviceId, page, limit, sort)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.Equal(0, total)
	})

	t.Run("should fail on list transaction by service destination", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockServiceTransactionUseCase()

		serviceId := uuid.NewV4().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		transactionUseCase.On("FindAllByServiceId", serviceId, page, limit, sort).Return(nil, 0, errors.New("internal error"))

		c := controller.NewServiceTransaction(transactionUseCase)

		result, total, err := c.FindAllByServiceId(context.TODO(), serviceId, page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "error on list payments by service destination")
		is.Equal(0, total)
	})

	t.Run("should fail on list transaction by service destination if list is empty", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockServiceTransactionUseCase()

		serviceId := uuid.NewV4().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		transactionUseCase.On("FindAllByServiceId", serviceId, page, limit, sort).Return(nil, 0, nil)

		c := controller.NewServiceTransaction(transactionUseCase)

		result, total, err := c.FindAllByServiceId(context.TODO(), serviceId, page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(result)
		is.EqualError(err, "no service payment was found")
		is.Equal(0, total)
	})

	t.Run("should succeed on list transaction by service destination", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockServiceTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		service, _ := entity.NewService("service", "service description", uuid.NewV4().String(), uuid.NewV4().String())

		transaction, _ := entity.NewTransaction(accountFrom, nil, service, nil, 19, "AOA")
		page := 1
		limit := 10
		sort := "created_at DESC"

		transactionUseCase.On("FindAllByServiceId", service.ID, page, limit, sort).Return([]*entity.Transaction{transaction}, 1, nil)

		c := controller.NewServiceTransaction(transactionUseCase)

		result, total, err := c.FindAllByServiceId(context.TODO(), service.ID, page, limit, sort)

		transactionUseCase.AssertExpectations(t)

		is.Nil(err)
		is.NotEmpty(result)
		is.Equal(len(result), 1)
		is.Equal(1, total)
		is.Equal(result[0], transaction)
	})
}

func TestFindOneByService(t *testing.T) {
	t.Parallel()

	t.Run("should fail on validate get transaction by service", func(t *testing.T) {
		is := require.New(t)

		transactionId := uuid.NewV1().String()
		serviceId := uuid.NewV1().String()

		c := controller.NewServiceTransaction(nil)

		result, err := c.FindOneByService(context.TODO(), serviceId, transactionId)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on get transaction by service", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockServiceTransactionUseCase()

		transactionId := uuid.NewV4().String()
		serviceId := uuid.NewV4().String()

		transactionUseCase.On("FindOneByService", serviceId, transactionId).Return(nil, errors.New("internal error"))

		c := controller.NewServiceTransaction(transactionUseCase)

		result, err := c.FindOneByService(context.TODO(), serviceId, transactionId)

		is.Nil(result)
		is.EqualError(err, "error on get payment by service")
	})

	t.Run("should fail on get transaction by service if returns nil", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockServiceTransactionUseCase()

		transactionId := uuid.NewV4().String()
		serviceId := uuid.NewV4().String()

		transactionUseCase.On("FindOneByService", serviceId, transactionId).Return(nil, nil)

		c := controller.NewServiceTransaction(transactionUseCase)

		result, err := c.FindOneByService(context.TODO(), serviceId, transactionId)

		is.Nil(result)
		is.EqualError(err, "no service payment was found")
	})

	t.Run("should succeed on get transaction by service", func(t *testing.T) {
		is := require.New(t)
		transactionUseCase := mock.NewMockServiceTransactionUseCase()

		accountFrom, _ := entity.NewAccount(200)
		service, _ := entity.NewService("service", "service description", uuid.NewV4().String(), uuid.NewV4().String())

		transaction, _ := entity.NewTransaction(accountFrom, nil, service, nil, 19, "AOA")

		transactionUseCase.On("FindOneByService", service.ID, transaction.ID).Return(transaction, nil)

		c := controller.NewServiceTransaction(transactionUseCase)

		result, err := c.FindOneByService(context.TODO(), service.ID, transaction.ID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result, transaction)
	})
}
