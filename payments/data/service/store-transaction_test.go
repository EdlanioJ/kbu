package service_test

import (
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu/payments/data/service"
	"github.com/EdlanioJ/kbu/payments/data/service/mock"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	uuid "github.com/satori/go.uuid"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterStoreTransaction(t *testing.T) {
	t.Parallel()
	t.Run("should fail on find account", func(t *testing.T) {
		mockAccountRepository := mock.NewMockAccountRepository()
		is := require.New(t)

		accountFromId := uuid.NewV4().String()
		storeId := uuid.NewV4().String()
		amount := 30.00
		currency := "AKZ"

		mockAccountRepository.On("Find", accountFromId).Return(&entity.Account{}, errors.New("account not found"))
		storeTransaction := service.NewStoreTransaction(mockAccountRepository, nil, nil)

		result, err := storeTransaction.RegisterStoreTransaction(accountFromId, storeId, amount, currency)

		mockAccountRepository.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "account not found")
	})
	t.Run("should fail on find store", func(t *testing.T) {
		mockAccountRepository := mock.NewMockAccountRepository()
		mockStoreRepository := mock.NewMockStoreRepository()
		is := require.New(t)

		accountFromId := uuid.NewV4().String()
		storeId := uuid.NewV4().String()
		amount := 30.00
		currency := "AKZ"

		mockAccountRepository.On("Find", accountFromId).Return(&entity.Account{}, nil)
		mockStoreRepository.On("FindStoreByIdAndStatus", storeId, entity.StoreActive).Return(&entity.Store{}, errors.New("store not found"))

		storeTransaction := service.NewStoreTransaction(mockAccountRepository, mockStoreRepository, nil)

		result, err := storeTransaction.RegisterStoreTransaction(accountFromId, storeId, amount, currency)

		mockAccountRepository.AssertExpectations(t)
		mockStoreRepository.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "store not found")
	})

	t.Run("should fail on withdrow account", func(t *testing.T) {
		mockAccountRepository := mock.NewMockAccountRepository()
		mockStoreRepository := mock.NewMockStoreRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(20)

		storeId := uuid.NewV4().String()
		amount := 30.00
		currency := "AKZ"

		mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockStoreRepository.On("FindStoreByIdAndStatus", storeId, entity.StoreActive).Return(&entity.Store{}, nil)

		storeTransaction := service.NewStoreTransaction(mockAccountRepository, mockStoreRepository, nil)

		result, err := storeTransaction.RegisterStoreTransaction(accountFrom.ID, storeId, amount, currency)

		mockAccountRepository.AssertExpectations(t)
		mockStoreRepository.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "account does not have balance")
	})
	t.Run("should fail on create new transaction", func(t *testing.T) {
		mockAccountRepository := mock.NewMockAccountRepository()
		mockStoreRepository := mock.NewMockStoreRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store 1", "store description 1", uuid.NewV4().String(), uuid.NewV4().String())

		amount := 0.00
		currency := "AKZ"

		mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockStoreRepository.On("FindStoreByIdAndStatus", store.ID, entity.StoreActive).Return(store, nil)

		storeTransaction := service.NewStoreTransaction(mockAccountRepository, mockStoreRepository, nil)

		result, err := storeTransaction.RegisterStoreTransaction(accountFrom.ID, store.ID, amount, currency)

		mockAccountRepository.AssertExpectations(t)
		mockStoreRepository.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "amount: Missing required field")
	})
	t.Run("should fail on register store transaction", func(t *testing.T) {
		mockAccountRepository := mock.NewMockAccountRepository()
		mockStoreRepository := mock.NewMockStoreRepository()
		mockTransactionRepository := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store 1", "store description 1", uuid.NewV4().String(), uuid.NewV4().String())

		amount := 40.00
		currency := "AKZ"

		mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockStoreRepository.On("FindStoreByIdAndStatus", store.ID, entity.StoreActive).Return(store, nil)
		mockTransactionRepository.On("Register", testifyMock.Anything).Return(errors.New("fail on register transaction"))

		storeTransaction := service.NewStoreTransaction(mockAccountRepository, mockStoreRepository, mockTransactionRepository)

		result, err := storeTransaction.RegisterStoreTransaction(accountFrom.ID, store.ID, amount, currency)

		mockAccountRepository.AssertExpectations(t)
		mockStoreRepository.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "fail on register transaction")
	})
	t.Run("should fail on save account", func(t *testing.T) {
		mockAccountRepository := mock.NewMockAccountRepository()
		mockStoreRepository := mock.NewMockStoreRepository()
		mockTransactionRepository := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store 1", "store description 1", uuid.NewV4().String(), uuid.NewV4().String())

		amount := 40.00
		currency := "AKZ"

		mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockStoreRepository.On("FindStoreByIdAndStatus", store.ID, entity.StoreActive).Return(store, nil)
		mockTransactionRepository.On("Register", testifyMock.Anything).Return(nil)
		_ = accountFrom.Withdow(amount)

		mockAccountRepository.On("Save", accountFrom).Return(errors.New("fail on save account"))

		storeTransaction := service.NewStoreTransaction(mockAccountRepository, mockStoreRepository, mockTransactionRepository)

		result, err := storeTransaction.RegisterStoreTransaction(accountFrom.ID, store.ID, amount, currency)

		mockAccountRepository.AssertExpectations(t)
		mockStoreRepository.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "fail on save account")
	})
	t.Run("should succeed on register store transaction", func(t *testing.T) {
		mockAccountRepository := mock.NewMockAccountRepository()
		mockStoreRepository := mock.NewMockStoreRepository()
		mockTransactionRepository := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store 1", "store description 1", uuid.NewV4().String(), uuid.NewV4().String())

		amount := 40.00
		currency := "AKZ"

		mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
		mockStoreRepository.On("FindStoreByIdAndStatus", store.ID, entity.StoreActive).Return(store, nil)
		mockTransactionRepository.On("Register", testifyMock.Anything).Return(nil)
		_ = accountFrom.Withdow(amount)

		mockAccountRepository.On("Save", accountFrom).Return(nil)

		storeTransaction := service.NewStoreTransaction(mockAccountRepository, mockStoreRepository, mockTransactionRepository)

		result, err := storeTransaction.RegisterStoreTransaction(accountFrom.ID, store.ID, amount, currency)

		mockAccountRepository.AssertExpectations(t)
		mockStoreRepository.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result.Amount, amount)
		is.Equal(result.AccountFromID, accountFrom.ID)
		is.Equal(result.StoreID, store.ID)
		is.Equal(result.Status, entity.TransactionPending)
	})
}

func TestFindAllByStoreId(t *testing.T) {
	t.Parallel()
	t.Run("should fail on find store", func(t *testing.T) {
		mockStoreRepository := mock.NewMockStoreRepository()
		is := require.New(t)

		storeId := uuid.NewV4().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		mockStoreRepository.On("Find", storeId).Return(&entity.Store{}, errors.New("store not found"))

		storeTransaction := service.NewStoreTransaction(nil, mockStoreRepository, nil)

		result, total, err := storeTransaction.FindAllByStoreId(storeId, page, limit, sort)

		mockStoreRepository.AssertExpectations(t)

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
		is.EqualError(err, "store not found")
	})

	t.Run("should fail on find transaction by store", func(t *testing.T) {
		mockStoreRepository := mock.NewMockStoreRepository()
		mockTransactionRepository := mock.NewMockTransactionRepository()
		is := require.New(t)

		storeId := uuid.NewV4().String()
		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		mockStoreRepository.On("Find", storeId).Return(&entity.Store{}, nil)
		mockTransactionRepository.On("FindByStoreId", storeId, pagination).Return([]*entity.Transaction{}, 0, errors.New("empty transactions"))

		storeTransaction := service.NewStoreTransaction(nil, mockStoreRepository, mockTransactionRepository)

		result, total, err := storeTransaction.FindAllByStoreId(storeId, page, limit, sort)

		mockStoreRepository.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
		is.EqualError(err, "empty transactions")
	})
	t.Run("should succeed on find transaction by store", func(t *testing.T) {
		mockStoreRepository := mock.NewMockStoreRepository()
		mockTransactionRepository := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store 1", "store description 1", uuid.NewV4().String(), uuid.NewV4().String())

		amount := 40.00
		currency := "AKZ"

		transaction, _ := entity.NewTransaction(accountFrom, nil, nil, store, amount, currency)
		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		mockStoreRepository.On("Find", store.ID).Return(store, nil)
		mockTransactionRepository.On("FindByStoreId", store.ID, pagination).Return([]*entity.Transaction{transaction}, 1, nil)

		storeTransaction := service.NewStoreTransaction(nil, mockStoreRepository, mockTransactionRepository)

		result, total, err := storeTransaction.FindAllByStoreId(store.ID, page, limit, sort)

		mockStoreRepository.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)

		is.Nil(err)
		is.Equal(total, 1)
		is.NotNil(result)
		is.Equal(result[0].Amount, amount)
		is.Equal(result[0].AccountFromID, accountFrom.ID)
		is.Equal(result[0].StoreID, store.ID)
		is.Equal(result[0].Status, entity.TransactionPending)

	})
}

func TestFindOneByStore(t *testing.T) {
	t.Parallel()
	t.Run("should fail on find store", func(t *testing.T) {
		mockStoreRepository := mock.NewMockStoreRepository()
		is := require.New(t)

		storeId := uuid.NewV4().String()
		transactionId := uuid.NewV4().String()

		mockStoreRepository.On("Find", storeId).Return(&entity.Store{}, errors.New("store not found"))

		storeTransaction := service.NewStoreTransaction(nil, mockStoreRepository, nil)

		result, err := storeTransaction.FindOneByStore(storeId, transactionId)

		mockStoreRepository.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "store not found")
	})

	t.Run("should fail on find one transaction by store", func(t *testing.T) {
		mockStoreRepository := mock.NewMockStoreRepository()
		mockTransactionRepository := mock.NewMockTransactionRepository()
		is := require.New(t)

		storeId := uuid.NewV4().String()
		transactionId := uuid.NewV4().String()

		mockStoreRepository.On("Find", storeId).Return(&entity.Store{}, nil)
		mockTransactionRepository.On("FindOneByStore", transactionId, storeId).Return(&entity.Transaction{}, errors.New("transaction not found"))

		storeTransaction := service.NewStoreTransaction(nil, mockStoreRepository, mockTransactionRepository)

		result, err := storeTransaction.FindOneByStore(storeId, transactionId)

		mockStoreRepository.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)

		is.Nil(result)
		is.NotNil(err)
		is.EqualError(err, "transaction not found")
	})
	t.Run("should succeed on find one by store", func(t *testing.T) {
		mockStoreRepository := mock.NewMockStoreRepository()
		mockTransactionRepository := mock.NewMockTransactionRepository()
		is := require.New(t)

		accountFrom, _ := entity.NewAccount(200)
		store, _ := entity.NewStore("store 1", "store description 1", uuid.NewV4().String(), uuid.NewV4().String())

		amount := 40.00
		currency := "AKZ"

		transaction, _ := entity.NewTransaction(accountFrom, nil, nil, store, amount, currency)

		mockStoreRepository.On("Find", store.ID).Return(&entity.Store{}, nil)
		mockTransactionRepository.On("FindOneByStore", transaction.ID, store.ID).Return(transaction, nil)

		storeTransaction := service.NewStoreTransaction(nil, mockStoreRepository, mockTransactionRepository)

		result, err := storeTransaction.FindOneByStore(store.ID, transaction.ID)

		mockStoreRepository.AssertExpectations(t)
		mockTransactionRepository.AssertExpectations(t)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result.Amount, amount)
		is.Equal(result.AccountFromID, accountFrom.ID)
		is.Equal(result.StoreID, store.ID)
		is.Equal(result.Status, entity.TransactionPending)
	})
}
