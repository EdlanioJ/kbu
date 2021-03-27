package service_test

import (
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu/payments/data/service"
	"github.com/EdlanioJ/kbu/payments/data/service/mock"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		id := uuid.NewV4().String()
		mockTransactionRepo.On("Find", id).Return(&entity.Transaction{}, errors.New("transaction not found"))
		transactionService := service.NewTransaction(mockTransactionRepo)

		result, err := transactionService.Find(id)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "transaction not found")
	})

	t.Run("should succeed on find", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000.93)
		accountTo, _ := entity.NewAccount(200)
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 200, "USD")

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transactionService := service.NewTransaction(mockTransactionRepo)

		result, err := transactionService.Find(transaction.ID)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result.ID, transaction.ID)
		is.Equal(transaction, result)
	})
}

func TestFindAll(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find all", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}
		mockTransactionRepo.On("FindAll", pagination).Return([]*entity.Transaction{}, 0, errors.New("empty list"))
		transactionService := service.NewTransaction(mockTransactionRepo)

		result, total, err := transactionService.FindAll(page, limit, sort)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
		is.EqualError(err, "empty list")
	})
	t.Run("should succeed on find all", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		accountFrom, _ := entity.NewAccount(3000.93)
		accountTo, _ := entity.NewAccount(200)
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 200, "USD")
		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}
		mockTransactionRepo.On("FindAll", pagination).Return([]*entity.Transaction{transaction}, 1, nil)
		transactionService := service.NewTransaction(mockTransactionRepo)

		result, total, err := transactionService.FindAll(page, limit, sort)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0], transaction)
	})
}

func TestComplete(t *testing.T) {
	t.Parallel()
	t.Run("should fail on complete", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		id := uuid.NewV4().String()

		mockTransactionRepo.On("Find", id).Return(&entity.Transaction{}, errors.New("transaction not found"))
		transactionService := service.NewTransaction(mockTransactionRepo)

		result, err := transactionService.Complete(id)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "transaction not found")
	})

	t.Run("should fail on save complete", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000.93)
		accountTo, _ := entity.NewAccount(200)
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 200, "USD")

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transaction.Status = entity.TransactionCompleted
		mockTransactionRepo.On("Save", transaction).Return(errors.New("failure on save"))

		transactionService := service.NewTransaction(mockTransactionRepo)
		result, err := transactionService.Complete(transaction.ID)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "failure on save")
	})

	t.Run("should succeed on complete", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000.93)
		accountTo, _ := entity.NewAccount(200)
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 200, "USD")

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transaction.Status = entity.TransactionCompleted
		mockTransactionRepo.On("Save", transaction).Return(nil)

		transactionService := service.NewTransaction(mockTransactionRepo)
		result, err := transactionService.Complete(transaction.ID)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result, transaction)
	})
}

func TestError(t *testing.T) {
	t.Parallel()
	t.Run("should fail on find transaction", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		id := uuid.NewV4().String()

		mockTransactionRepo.On("Find", id).Return(&entity.Transaction{}, errors.New("transaction not found"))
		transactionService := service.NewTransaction(mockTransactionRepo)

		result, err := transactionService.Error(id)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "transaction not found")
	})
	t.Run("should fail on save transaction", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000.93)
		accountTo, _ := entity.NewAccount(200)
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 200, "USD")

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transaction.Status = entity.TransactionError
		mockTransactionRepo.On("Save", transaction).Return(errors.New("failure on save"))

		transactionService := service.NewTransaction(mockTransactionRepo)
		result, err := transactionService.Error(transaction.ID)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "failure on save")
	})

	t.Run("should succeed error", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000.93)
		accountTo, _ := entity.NewAccount(200)
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 200, "USD")

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transaction.Status = entity.TransactionError
		mockTransactionRepo.On("Save", transaction).Return(nil)

		transactionService := service.NewTransaction(mockTransactionRepo)
		result, err := transactionService.Error(transaction.ID)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)

	})
}
