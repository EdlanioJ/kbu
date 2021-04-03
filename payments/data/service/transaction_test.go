package service_test

import (
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu/payments/data/service"
	"github.com/EdlanioJ/kbu/payments/data/service/mock"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	uuid "github.com/satori/go.uuid"
	tMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find account from", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		is := require.New(t)

		fromID := uuid.NewV4().String()

		mockAccountRepo.On("Find", fromID).Return(nil, errors.New("invalid user"))
		transactionService := service.NewTransaction(nil, mockAccountRepo)

		result, err := transactionService.Register(fromID, "", "", "", "", 0)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail if find from return nil transaction", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		is := require.New(t)

		fromID := uuid.NewV4().String()
		mockAccountRepo.On("Find", fromID).Return(nil, nil)
		transactionService := service.NewTransaction(nil, mockAccountRepo)

		result, err := transactionService.Register(fromID, "", "", "", "", 0)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "no account from was found")
	})

	t.Run("should fail on find account destination", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000.93)

		toID := uuid.NewV4().String()
		mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockAccountRepo.On("Find", toID).Return(nil, errors.New("invalid param"))

		transactionService := service.NewTransaction(nil, mockAccountRepo)
		result, err := transactionService.Register(accountFrom.ID, toID, "", "", "", 0)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail if find destination return nil transaction", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000.93)

		toID := uuid.NewV4().String()
		mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockAccountRepo.On("Find", toID).Return(nil, nil)

		transactionService := service.NewTransaction(nil, mockAccountRepo)
		result, err := transactionService.Register(accountFrom.ID, toID, "", "", "", 0)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "no account destination was found")
	})

	t.Run("should fail on withdrow from account", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(30)
		accountTo, _ := entity.NewAccount(200)

		mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)

		transactionService := service.NewTransaction(nil, mockAccountRepo)
		result, err := transactionService.Register(accountFrom.ID, accountTo.ID, "", "", "", 40)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "account does not have balance")
	})

	t.Run("should fail on new transaction", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(300.00)
		accountTo, _ := entity.NewAccount(200)

		mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)

		transactionService := service.NewTransaction(nil, mockAccountRepo)
		result, err := transactionService.Register(accountFrom.ID, accountTo.ID, "", "", "", 0)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should fail on register transaction", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)

		mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		mockTransactionRepo.On("Register", tMock.Anything).Return(errors.New("register error"))

		transactionService := service.NewTransaction(mockTransactionRepo, mockAccountRepo)
		result, err := transactionService.Register(accountFrom.ID, accountTo.ID, externalID, transactionType, currency, 40)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "register error")
	})

	t.Run("should fail on save account from", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00

		mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)
		mockTransactionRepo.On("Register", tMock.Anything).Return(nil)
		_ = accountFrom.Withdow(amount)
		mockAccountRepo.On("Save", accountFrom).Return(errors.New("error on save"))

		transactionService := service.NewTransaction(mockTransactionRepo, mockAccountRepo)
		result, err := transactionService.Register(accountFrom.ID, accountTo.ID, externalID, transactionType, currency, 40)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
		is.EqualError(err, "error on save")
	})

	t.Run("should succeed", func(t *testing.T) {
		mockAccountRepo := mock.NewMockAccountRepository()
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00

		mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)
		mockTransactionRepo.On("Register", tMock.Anything).Return(nil)
		_ = accountFrom.Withdow(amount)
		mockAccountRepo.On("Save", accountFrom).Return(nil)

		transactionService := service.NewTransaction(mockTransactionRepo, mockAccountRepo)
		result, err := transactionService.Register(accountFrom.ID, accountTo.ID, externalID, transactionType, currency, amount)

		is.Nil(err)
		is.Equal(result.AccountFromID, accountFrom.ID)
		is.Equal(result.AccountToID, accountTo.ID)
		is.Equal(result.Amount, amount)
		is.Equal(result.Currency, currency)
		is.Equal(result.ExternalID, externalID)
		is.Equal(result.Type, transactionType)
	})
}

func TestFindByType(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find transaction by type", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		transactionID := uuid.NewV4().String()
		transactionType := entity.TransactionToUser

		mockTransactionRepo.On("FindByType", transactionID, transactionType).Return(nil, errors.New("error on find"))
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, err := serviceTransaction.FindByType(transactionType, transactionID)
		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		transactionType := entity.TransactionToUser
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

		mockTransactionRepo.On("FindByType", transaction.ID, transactionType).Return(transaction, nil)
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, err := serviceTransaction.FindByType(transactionType, transaction.ID)
		mockTransactionRepo.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}

func TestFindAllByType(t *testing.T) {
	t.Parallel()
	t.Run("should fail on find all by type", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		page := 1
		limit := 10
		sort := "created_at DESC"
		transactionType := entity.TransactionToService
		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		mockTransactionRepo.On("FindAllByType", transactionType, pagination).Return(nil, 0, errors.New("error on find"))
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, total, err := serviceTransaction.FindAllByType(transactionType, page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
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

		transactionType := entity.TransactionToService
		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)
		transactions := []*entity.Transaction{transaction}
		totalResult := len(transactions)
		mockTransactionRepo.On("FindAllByType", transactionType, pagination).Return(transactions, totalResult, nil)
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, total, err := serviceTransaction.FindAllByType(transactionType, page, limit, sort)

		is.Nil(err)
		is.Equal(totalResult, total)
		is.Equal(result[0], transaction)

	})
}

func TestFindByExternalID(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find transaction by external id", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		transactionID := uuid.NewV4().String()
		externalID := uuid.NewV4().String()

		mockTransactionRepo.On("FindByExternalID", transactionID, externalID).Return(nil, errors.New("error on find"))
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, err := serviceTransaction.FindByExternalID(externalID, transactionID)
		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		transactionID := uuid.NewV4().String()
		externalID := uuid.NewV4().String()

		mockTransactionRepo.On("FindByExternalID", transactionID, externalID).Return(nil, nil)
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, err := serviceTransaction.FindByExternalID(externalID, transactionID)
		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.Nil(err)
	})
}

func TestFindAllByExternalID(t *testing.T) {
	t.Parallel()
	t.Run("should fail on find all by external id", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		page := 1
		limit := 10
		sort := "created_at DESC"
		externalID := uuid.NewV4().String()
		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		mockTransactionRepo.On("FindAllByExternalID", externalID, pagination).Return(nil, 0, errors.New("error on find"))
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, total, err := serviceTransaction.FindAllByExternalID(externalID, page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		page := 1
		limit := 10
		sort := "created_at DESC"
		transactionType := entity.TransactionToService
		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)
		transactions := []*entity.Transaction{transaction}
		totalResult := len(transactions)
		mockTransactionRepo.On("FindAllByExternalID", transaction.ExternalID, pagination).Return(transactions, totalResult, nil)
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, total, err := serviceTransaction.FindAllByExternalID(transaction.ExternalID, page, limit, sort)

		is.Nil(err)
		is.Equal(totalResult, total)
		is.Equal(result[0], transaction)

	})
}

func TestFindByFromAccountID(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find transaction by account from id", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		transactionID := uuid.NewV4().String()
		accountID := uuid.NewV4().String()

		mockTransactionRepo.On("FindByFromAccountID", transactionID, accountID).Return(nil, errors.New("error on find"))
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, err := serviceTransaction.FindByFromAccountID(accountID, transactionID)
		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		transactionID := uuid.NewV4().String()
		accountID := uuid.NewV4().String()

		mockTransactionRepo.On("FindByFromAccountID", transactionID, accountID).Return(nil, nil)
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, err := serviceTransaction.FindByFromAccountID(accountID, transactionID)
		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.Nil(err)
	})
}

func TestFindAllByFromAccountID(t *testing.T) {
	t.Parallel()
	t.Run("should fail on find all by external id", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		page := 1
		limit := 10
		sort := "created_at DESC"
		accountID := uuid.NewV4().String()
		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		mockTransactionRepo.On("FindAllByFromAccountID", accountID, pagination).Return(nil, 0, errors.New("error on find"))
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, total, err := serviceTransaction.FindAllByFromAccountID(accountID, page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		page := 1
		limit := 10
		sort := "created_at DESC"
		transactionType := entity.TransactionToService
		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)
		transactions := []*entity.Transaction{transaction}
		totalResult := len(transactions)
		mockTransactionRepo.On("FindAllByFromAccountID", transaction.AccountFromID, pagination).Return(transactions, totalResult, nil)
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, total, err := serviceTransaction.FindAllByFromAccountID(transaction.AccountFromID, page, limit, sort)

		is.Nil(err)
		is.Equal(totalResult, total)
		is.Equal(result[0], transaction)

	})
}

func TestFindByToAccountID(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find transaction by account from id", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		transactionID := uuid.NewV4().String()
		accountID := uuid.NewV4().String()

		mockTransactionRepo.On("FindByToAccountID", transactionID, accountID).Return(nil, errors.New("error on find"))
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, err := serviceTransaction.FindByToAccountID(accountID, transactionID)
		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		transactionID := uuid.NewV4().String()
		accountID := uuid.NewV4().String()

		mockTransactionRepo.On("FindByToAccountID", transactionID, accountID).Return(nil, nil)
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, err := serviceTransaction.FindByToAccountID(accountID, transactionID)
		mockTransactionRepo.AssertExpectations(t)

		is.Nil(result)
		is.Nil(err)
	})
}

func TestFindAllByToAccountID(t *testing.T) {
	t.Parallel()
	t.Run("should fail on find all by external id", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		page := 1
		limit := 10
		sort := "created_at DESC"
		accountID := uuid.NewV4().String()
		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		mockTransactionRepo.On("FindAllByToAccountID", accountID, pagination).Return(nil, 0, errors.New("error on find"))
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, total, err := serviceTransaction.FindAllByToAccountID(accountID, page, limit, sort)

		is.Nil(result)
		is.Equal(0, total)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)
		page := 1
		limit := 10
		sort := "created_at DESC"
		transactionType := entity.TransactionToService
		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		accountFrom, _ := entity.NewAccount(3000)
		accountTo, _ := entity.NewAccount(200)
		externalID := uuid.NewV4().String()
		currency := "AOA"
		amount := 30.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)
		transactions := []*entity.Transaction{transaction}
		totalResult := len(transactions)
		mockTransactionRepo.On("FindAllByToAccountID", transaction.AccountToID, pagination).Return(transactions, totalResult, nil)
		serviceTransaction := service.NewTransaction(mockTransactionRepo, nil)

		result, total, err := serviceTransaction.FindAllByToAccountID(transaction.AccountToID, page, limit, sort)

		is.Nil(err)
		is.Equal(totalResult, total)
		is.Equal(result[0], transaction)
	})
}

func TestFind(t *testing.T) {
	t.Parallel()

	t.Run("should fail on find", func(t *testing.T) {
		mockTransactionRepo := mock.NewMockTransactionRepository()
		is := require.New(t)

		id := uuid.NewV4().String()
		mockTransactionRepo.On("Find", id).Return(&entity.Transaction{}, errors.New("transaction not found"))
		transactionService := service.NewTransaction(mockTransactionRepo, nil)

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
		externalID := uuid.NewV4().String()
		transactionType := entity.TransactionToUser
		amount := 200.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, "AOA", amount)

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transactionService := service.NewTransaction(mockTransactionRepo, nil)

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
		transactionService := service.NewTransaction(mockTransactionRepo, nil)

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
		externalID := uuid.NewV4().String()
		transactionType := entity.TransactionToUser
		amount := 200.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, "AOA", amount)
		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}
		mockTransactionRepo.On("FindAll", pagination).Return([]*entity.Transaction{transaction}, 1, nil)
		transactionService := service.NewTransaction(mockTransactionRepo, nil)

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
		transactionService := service.NewTransaction(mockTransactionRepo, nil)

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
		externalID := uuid.NewV4().String()
		transactionType := entity.TransactionToUser
		amount := 200.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, "AOA", amount)

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transaction.Status = entity.TransactionCompleted
		mockTransactionRepo.On("Save", transaction).Return(errors.New("failure on save"))

		transactionService := service.NewTransaction(mockTransactionRepo, nil)
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
		externalID := uuid.NewV4().String()
		transactionType := entity.TransactionToUser
		amount := 200.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, "AOA", amount)

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transaction.Status = entity.TransactionCompleted
		mockTransactionRepo.On("Save", transaction).Return(nil)

		transactionService := service.NewTransaction(mockTransactionRepo, nil)
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
		transactionService := service.NewTransaction(mockTransactionRepo, nil)

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
		externalID := uuid.NewV4().String()
		transactionType := entity.TransactionToUser
		amount := 200.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, "AOA", amount)

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transaction.Status = entity.TransactionCanceled
		mockTransactionRepo.On("Save", transaction).Return(errors.New("failure on save"))

		transactionService := service.NewTransaction(mockTransactionRepo, nil)
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
		externalID := uuid.NewV4().String()
		transactionType := entity.TransactionToUser
		amount := 200.00
		transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, "AOA", amount)

		mockTransactionRepo.On("Find", transaction.ID).Return(transaction, nil)
		transaction.Status = entity.TransactionCanceled
		mockTransactionRepo.On("Save", transaction).Return(nil)

		transactionService := service.NewTransaction(mockTransactionRepo, nil)
		result, err := transactionService.Error(transaction.ID)

		mockTransactionRepo.AssertExpectations(t)

		is.Nil(err)
		is.Equal(result, transaction)
	})
}
