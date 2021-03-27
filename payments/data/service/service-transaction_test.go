package service_test

import (
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu/payments/data/service"
	repoMock "github.com/EdlanioJ/kbu/payments/data/service/mock"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_RegisterServiceTransaction_Fail_On_MissingParams(t *testing.T) {
	serviceTransaction := service.NewServiceTransaction(nil, nil, nil, nil)

	accountFromId := uuid.NewV4().String()
	serviceId := uuid.NewV4().String()
	result, err := serviceTransaction.RegisterServiceTransaction(accountFromId, serviceId, "", 0, "AKZ")

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "missing amount or servicePriceId")
}

func Test_RegisterServiceTransaction_Fail_On_FindAccount(t *testing.T) {
	mockAccountRepository := new(repoMock.MockAccountRepository)

	accountFromId := uuid.NewV4().String()
	serviceId := uuid.NewV4().String()
	amount := 30.00

	mockAccountRepository.On("Find", accountFromId).Return(&entity.Account{}, errors.New("account not found"))
	serviceTransaction := service.NewServiceTransaction(mockAccountRepository, nil, nil, nil)

	result, err := serviceTransaction.RegisterServiceTransaction(accountFromId, serviceId, "", amount, "AKZ")

	mockAccountRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "account not found")
}

func Test_RegisterServiceTransaction_Fail_On_FindService(t *testing.T) {
	mockAccountRepository := new(repoMock.MockAccountRepository)
	mockServiceRepository := new(repoMock.MockServiceRepository)

	accountFromId := uuid.NewV4().String()
	serviceId := uuid.NewV4().String()
	amount := 30.00

	mockAccountRepository.On("Find", accountFromId).Return(&entity.Account{}, nil)
	mockServiceRepository.On("FindServiceByIdAndStatus", serviceId, entity.ServiceActive).Return(&entity.Service{}, errors.New("service not found"))
	serviceTransaction := service.NewServiceTransaction(mockAccountRepository, nil, mockServiceRepository, nil)

	result, err := serviceTransaction.RegisterServiceTransaction(accountFromId, serviceId, "", amount, "AKZ")

	mockAccountRepository.AssertExpectations(t)
	mockServiceRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "service not found")
}

func Test_RegisterServiceTransaction_Fail_On_FindServicePrice(t *testing.T) {
	mockAccountRepository := new(repoMock.MockAccountRepository)
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockServicePriceRepository := new(repoMock.MockServicePriceRepository)

	accountFromId := uuid.NewV4().String()
	serviceId := uuid.NewV4().String()
	servicePriceId := uuid.NewV4().String()

	mockAccountRepository.On("Find", accountFromId).Return(&entity.Account{}, nil)
	mockServiceRepository.On("FindServiceByIdAndStatus", serviceId, entity.ServiceActive).Return(&entity.Service{}, nil)
	mockServicePriceRepository.On("Find", servicePriceId).Return(&entity.ServicePrice{}, errors.New("service price not found"))

	serviceTransaction := service.NewServiceTransaction(mockAccountRepository, mockServicePriceRepository, mockServiceRepository, nil)

	result, err := serviceTransaction.RegisterServiceTransaction(accountFromId, serviceId, servicePriceId, 0, "AKZ")

	mockAccountRepository.AssertExpectations(t)
	mockServiceRepository.AssertExpectations(t)
	mockServicePriceRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "service price not found")
}

func Test_RegisterServiceTransaction_Fail_On_WithdrowAccount(t *testing.T) {
	mockAccountRepository := new(repoMock.MockAccountRepository)
	mockServiceRepository := new(repoMock.MockServiceRepository)

	accountFrom, _ := entity.NewAccount(20)
	serviceId := uuid.NewV4().String()
	amount := 30.00

	mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockServiceRepository.On("FindServiceByIdAndStatus", serviceId, entity.ServiceActive).Return(&entity.Service{}, nil)

	serviceTransaction := service.NewServiceTransaction(mockAccountRepository, nil, mockServiceRepository, nil)

	result, err := serviceTransaction.RegisterServiceTransaction(accountFrom.ID, serviceId, "", amount, "AKZ")

	mockAccountRepository.AssertExpectations(t)
	mockServiceRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "account does not have balance")
}

func Test_RegisterServiceTransaction_Fail_On_NewTransaction(t *testing.T) {
	mockAccountRepository := new(repoMock.MockAccountRepository)
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockServicePriceRepository := new(repoMock.MockServicePriceRepository)

	accountFrom, _ := entity.NewAccount(200)
	serviceEntity, _ := entity.NewService("service 1", "service 1 descriptions", uuid.NewV4().String(), uuid.NewV4().String())
	servicePriceId := uuid.NewV4().String()

	mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockServiceRepository.On("FindServiceByIdAndStatus", serviceEntity.ID, entity.ServiceActive).Return(serviceEntity, nil)
	mockServicePriceRepository.On("Find", servicePriceId).Return(&entity.ServicePrice{Amount: 0}, nil)

	serviceTransaction := service.NewServiceTransaction(mockAccountRepository, mockServicePriceRepository, mockServiceRepository, nil)

	result, err := serviceTransaction.RegisterServiceTransaction(accountFrom.ID, serviceEntity.ID, servicePriceId, 0, "AKZ")

	mockAccountRepository.AssertExpectations(t)
	mockServiceRepository.AssertExpectations(t)
	mockServicePriceRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "amount: Missing required field")
}

func Test_RegisterServiceTransaction_Fail_On_RegisterTransaction(t *testing.T) {
	mockAccountRepository := new(repoMock.MockAccountRepository)
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockTransactionRepository := new(repoMock.MockTransactionRepository)
	accountFrom, _ := entity.NewAccount(200)
	serviceEntity, _ := entity.NewService("service 1", "service 1 descriptions", uuid.NewV4().String(), uuid.NewV4().String())
	amount := 29.00

	mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockServiceRepository.On("FindServiceByIdAndStatus", serviceEntity.ID, entity.ServiceActive).Return(serviceEntity, nil)

	_ = accountFrom.Withdow(amount)
	mockTransactionRepository.On("Register", mock.Anything).Return(errors.New("fail on register transaction"))

	serviceTransaction := service.NewServiceTransaction(mockAccountRepository, nil, mockServiceRepository, mockTransactionRepository)

	result, err := serviceTransaction.RegisterServiceTransaction(accountFrom.ID, serviceEntity.ID, "", amount, "AKZ")

	mockAccountRepository.AssertExpectations(t)
	mockServiceRepository.AssertExpectations(t)
	mockTransactionRepository.AssertExpectations((t))

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "fail on register transaction")
}

func Test_RegisterServiceTransaction_Fail_On_SaveAccount(t *testing.T) {
	mockAccountRepository := new(repoMock.MockAccountRepository)
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockTransactionRepository := new(repoMock.MockTransactionRepository)
	accountFrom, _ := entity.NewAccount(200)
	serviceEntity, _ := entity.NewService("service 1", "service 1 descriptions", uuid.NewV4().String(), uuid.NewV4().String())
	amount := 29.00

	mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockServiceRepository.On("FindServiceByIdAndStatus", serviceEntity.ID, entity.ServiceActive).Return(serviceEntity, nil)

	_ = accountFrom.Withdow(amount)
	mockTransactionRepository.On("Register", mock.Anything).Return(nil)
	mockAccountRepository.On("Save", accountFrom).Return(errors.New("fail on save account"))

	serviceTransaction := service.NewServiceTransaction(mockAccountRepository, nil, mockServiceRepository, mockTransactionRepository)

	result, err := serviceTransaction.RegisterServiceTransaction(accountFrom.ID, serviceEntity.ID, "", amount, "AKZ")

	mockAccountRepository.AssertExpectations(t)
	mockServiceRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "fail on save account")
}

func Test_RegisterServiceTransaction_Success(t *testing.T) {
	mockAccountRepository := new(repoMock.MockAccountRepository)
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockTransactionRepository := new(repoMock.MockTransactionRepository)
	accountFrom, _ := entity.NewAccount(200)
	serviceEntity, _ := entity.NewService("service 1", "service 1 descriptions", uuid.NewV4().String(), uuid.NewV4().String())
	amount := 29.00

	mockAccountRepository.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockServiceRepository.On("FindServiceByIdAndStatus", serviceEntity.ID, entity.ServiceActive).Return(serviceEntity, nil)

	_ = accountFrom.Withdow(amount)
	mockAccountRepository.On("Save", accountFrom).Return(nil)
	mockTransactionRepository.On("Register", mock.Anything).Return(nil)

	serviceTransaction := service.NewServiceTransaction(mockAccountRepository, nil, mockServiceRepository, mockTransactionRepository)

	result, err := serviceTransaction.RegisterServiceTransaction(accountFrom.ID, serviceEntity.ID, "", amount, "AKZ")

	mockAccountRepository.AssertExpectations(t)
	mockServiceRepository.AssertExpectations(t)
	mockTransactionRepository.AssertExpectations((t))

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, amount, result.Amount)
	assert.Equal(t, entity.ServicePending, result.Status)
}

func Test_FindAllByServiceId_Fail_On_FindService(t *testing.T) {
	mockServiceRepository := new(repoMock.MockServiceRepository)

	id := uuid.NewV4().String()
	page := 1
	limit := 10
	sort := "created_at DESC"

	mockServiceRepository.On("Find", id).Return(&entity.Service{}, errors.New("service not found"))

	serviceTransaction := service.NewServiceTransaction(nil, nil, mockServiceRepository, nil)

	result, total, err := serviceTransaction.FindAllByServiceId(id, page, limit, sort)

	mockServiceRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, total, 0)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "service not found")
}

func Test_FindAllByServiceId_Fail_On_FindTransactions(t *testing.T) {
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockTransactionRepository := new(repoMock.MockTransactionRepository)
	serviceEntity, _ := entity.NewService("service 1", "service 1 descriptions", uuid.NewV4().String(), uuid.NewV4().String())
	page := 1
	limit := 10
	sort := "created_at DESC"

	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	mockServiceRepository.On("Find", serviceEntity.ID).Return(serviceEntity, nil)
	mockTransactionRepository.On("FindByServiceId", serviceEntity.ID, pagination).Return([]*entity.Transaction{}, 0, errors.New("empty transaction"))
	serviceTransaction := service.NewServiceTransaction(nil, nil, mockServiceRepository, mockTransactionRepository)

	result, total, err := serviceTransaction.FindAllByServiceId(serviceEntity.ID, page, limit, sort)

	mockServiceRepository.AssertExpectations(t)
	mockTransactionRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, total, 0)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "empty transaction")
}

func Test_FindAllByServiceId_Success(t *testing.T) {
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockTransactionRepository := new(repoMock.MockTransactionRepository)
	account, _ := entity.NewAccount(200)
	serviceEntity, _ := entity.NewService("service 1", "service 1 descriptions", uuid.NewV4().String(), uuid.NewV4().String())
	transaction, _ := entity.NewTransaction(account, nil, serviceEntity, nil, 10, "AKZ")
	page := 1
	limit := 10
	sort := "created_at DESC"

	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	mockServiceRepository.On("Find", serviceEntity.ID).Return(serviceEntity, nil)
	mockTransactionRepository.On("FindByServiceId", serviceEntity.ID, pagination).Return([]*entity.Transaction{transaction}, 1, nil)
	serviceTransaction := service.NewServiceTransaction(nil, nil, mockServiceRepository, mockTransactionRepository)

	result, total, err := serviceTransaction.FindAllByServiceId(serviceEntity.ID, page, limit, sort)

	mockServiceRepository.AssertExpectations(t)
	mockTransactionRepository.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, total, 1)
	assert.Equal(t, result[0].ID, transaction.ID)
	assert.Equal(t, result[0].AccountFromID, transaction.AccountFromID)
	assert.Equal(t, result[0].Amount, transaction.Amount)
	assert.Equal(t, result[0].ServiceID, transaction.ServiceID)
}

func Test_FindOneByService_Fail_On_FindService(t *testing.T) {
	mockServiceRepository := new(repoMock.MockServiceRepository)

	serviceId := uuid.NewV4().String()
	transactionId := uuid.NewV4().String()

	mockServiceRepository.On("Find", serviceId).Return(&entity.Service{}, errors.New("service not found"))

	serviceTransaction := service.NewServiceTransaction(nil, nil, mockServiceRepository, nil)

	result, err := serviceTransaction.FindOneByService(serviceId, transactionId)

	mockServiceRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "service not found")
}

func Test_FindOneByService_Fail_On_FindTransaction(t *testing.T) {
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockTransactionRepository := new(repoMock.MockTransactionRepository)
	serviceEntity, _ := entity.NewService("service 1", "service 1 descriptions", uuid.NewV4().String(), uuid.NewV4().String())
	transactionId := uuid.NewV4().String()

	mockServiceRepository.On("Find", serviceEntity.ID).Return(serviceEntity, nil)
	mockTransactionRepository.On("FindOneByService", transactionId, serviceEntity.ID).Return(&entity.Transaction{}, errors.New("transaction not found"))
	serviceTransaction := service.NewServiceTransaction(nil, nil, mockServiceRepository, mockTransactionRepository)

	result, err := serviceTransaction.FindOneByService(serviceEntity.ID, transactionId)

	mockServiceRepository.AssertExpectations(t)
	mockTransactionRepository.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "transaction not found")
}

func Test_FindOneByService_Success(t *testing.T) {
	mockServiceRepository := new(repoMock.MockServiceRepository)
	mockTransactionRepository := new(repoMock.MockTransactionRepository)
	account, _ := entity.NewAccount(200)
	serviceEntity, _ := entity.NewService("service 1", "service 1 descriptions", uuid.NewV4().String(), uuid.NewV4().String())
	transaction, _ := entity.NewTransaction(account, nil, serviceEntity, nil, 10, "AKZ")

	mockServiceRepository.On("Find", serviceEntity.ID).Return(serviceEntity, nil)
	mockTransactionRepository.On("FindOneByService", transaction.ID, serviceEntity.ID).Return(transaction, nil)
	serviceTransaction := service.NewServiceTransaction(nil, nil, mockServiceRepository, mockTransactionRepository)

	result, err := serviceTransaction.FindOneByService(serviceEntity.ID, transaction.ID)

	mockServiceRepository.AssertExpectations(t)
	mockTransactionRepository.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, result.ID, transaction.ID)
	assert.Equal(t, result.AccountFromID, transaction.AccountFromID)
	assert.Equal(t, result.Amount, transaction.Amount)
	assert.Equal(t, result.ServiceID, transaction.ServiceID)
}
