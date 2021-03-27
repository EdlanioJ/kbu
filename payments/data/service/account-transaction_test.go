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

func Test_RegisterAccountTransaction_Success(t *testing.T) {
	mockTransactionRepo := new(repoMock.MockTransactionRepository)
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountFrom, _ := entity.NewAccount(3000.93)
	accountTo, _ := entity.NewAccount(200)
	amount := 200.00
	currency := "AKZ"

	mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)
	mockTransactionRepo.On("Register", mock.Anything).Return(nil)

	_ = accountFrom.Withdow(amount)
	mockAccountRepo.On("Save", accountFrom).Return(nil)
	accountTransaction := service.NewAccountTransaction(mockTransactionRepo, mockAccountRepo)

	result, err := accountTransaction.RegisterAccountTransaction(accountFrom.ID, accountTo.ID, amount, currency)

	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, amount, result.Amount)
}

func Test_RegisterAccountTransaction_Failure_On_FindAccountFrom(t *testing.T) {
	mockAccountRepo := new(repoMock.MockAccountRepository)

	id := uuid.NewV4().String()
	amount := 200.00
	currency := "AKZ"
	mockAccountRepo.On("Find", id).Return(&entity.Account{}, errors.New("Account from not found"))

	accountTransaction := service.NewAccountTransaction(nil, mockAccountRepo)

	result, err := accountTransaction.RegisterAccountTransaction(id, uuid.NewV1().String(), amount, currency)

	mockAccountRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Account from not found")
}

func Test_RegisterAccountTransaction_Failure_On_FindAccountTo(t *testing.T) {
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountFrom, _ := entity.NewAccount(3000.93)
	id := uuid.NewV4().String()
	amount := 200.00
	currency := "AKZ"

	mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockAccountRepo.On("Find", id).Return(&entity.Account{}, errors.New("Account To not found"))

	accountTransaction := service.NewAccountTransaction(nil, mockAccountRepo)

	result, err := accountTransaction.RegisterAccountTransaction(accountFrom.ID, id, amount, currency)

	mockAccountRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Account To not found")
}

func Test_RegisterAccountTransaction_Failure_On_WithdrowAccount(t *testing.T) {
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountFrom, _ := entity.NewAccount(3.00)
	accountTo, _ := entity.NewAccount(200)
	amount := 200.00
	currency := "AKZ"

	mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)

	accountTransaction := service.NewAccountTransaction(nil, mockAccountRepo)

	result, err := accountTransaction.RegisterAccountTransaction(accountFrom.ID, accountTo.ID, amount, currency)

	mockAccountRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "account does not have balance")
}

func Test_RegisterAccountTransaction_Failure_On_NewTransaction(t *testing.T) {
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountFrom, _ := entity.NewAccount(3000.93)
	accountTo, _ := entity.NewAccount(200)
	amount := 0.00
	currency := "AKZ"

	mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)

	accountTransaction := service.NewAccountTransaction(nil, mockAccountRepo)

	result, err := accountTransaction.RegisterAccountTransaction(accountFrom.ID, accountTo.ID, amount, currency)

	mockAccountRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "amount: Missing required field")

}

func Test_RegisterAccountTransaction_Failure_On_RegisterTransaction(t *testing.T) {
	mockTransactionRepo := new(repoMock.MockTransactionRepository)
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountFrom, _ := entity.NewAccount(3000.93)
	accountTo, _ := entity.NewAccount(200)
	amount := 200.00
	currency := "AKZ"

	mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)
	mockTransactionRepo.On("Register", mock.Anything).Return(errors.New("Fail on register"))

	accountTransaction := service.NewAccountTransaction(mockTransactionRepo, mockAccountRepo)

	result, err := accountTransaction.RegisterAccountTransaction(accountFrom.ID, accountTo.ID, amount, currency)

	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Fail on register")
}

func Test_RegisterAccountTransaction_Failure_On_SaveAccountFrom(t *testing.T) {
	mockTransactionRepo := new(repoMock.MockTransactionRepository)
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountFrom, _ := entity.NewAccount(3000.93)
	accountTo, _ := entity.NewAccount(200)
	amount := 200.00
	currency := "AKZ"

	mockAccountRepo.On("Find", accountFrom.ID).Return(accountFrom, nil)
	mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)
	mockTransactionRepo.On("Register", mock.Anything).Return(nil)

	_ = accountFrom.Withdow(amount)
	mockAccountRepo.On("Save", accountFrom).Return(errors.New("fail on save account from"))
	accountTransaction := service.NewAccountTransaction(mockTransactionRepo, mockAccountRepo)

	result, err := accountTransaction.RegisterAccountTransaction(accountFrom.ID, accountTo.ID, amount, currency)

	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "fail on save account from")
}

func Test_FindAllByAccountTo_Success(t *testing.T) {
	mockTransactionRepo := new(repoMock.MockTransactionRepository)
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountFrom, _ := entity.NewAccount(3000.93)
	accountTo, _ := entity.NewAccount(200)
	amount := 200.00
	currency := "AKZ"
	transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, amount, currency)
	page := 1
	limit := 10
	sort := "created_at DESC"

	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	mockAccountRepo.On("Find", accountTo.ID).Return(accountFrom, nil)
	mockTransactionRepo.On("FindByAccountToId", accountTo.ID, pagination).Return([]*entity.Transaction{transaction}, 1, nil)

	accountTransaction := service.NewAccountTransaction(mockTransactionRepo, mockAccountRepo)

	result, total, err := accountTransaction.FindAllByAccountTo(accountTo.ID, page, limit, sort)
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, total, 1)
	assert.Equal(t, result[0].ID, transaction.ID)
	assert.Equal(t, result[0].AccountFromID, transaction.AccountFromID)
	assert.Equal(t, result[0].Amount, transaction.Amount)
	assert.Equal(t, result[0].AccountToID, transaction.AccountToID)
}

func Test_FindAllByAccountTo_Fail_On_FindAccount(t *testing.T) {
	mockAccountRepo := new(repoMock.MockAccountRepository)
	id := uuid.NewV4().String()
	page := 1
	limit := 10
	sort := "created_at DESC"

	mockAccountRepo.On("Find", id).Return(&entity.Account{}, errors.New("account not found"))

	accountTransaction := service.NewAccountTransaction(nil, mockAccountRepo)

	result, total, err := accountTransaction.FindAllByAccountTo(id, page, limit, sort)
	mockAccountRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, total, 0)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "account not found")
}

func Test_FindAllByAccountTo_Fail_On_FindTransaction(t *testing.T) {
	mockTransactionRepo := new(repoMock.MockTransactionRepository)
	mockAccountRepo := new(repoMock.MockAccountRepository)

	account, _ := entity.NewAccount(3000.93)
	page := 1
	limit := 10
	sort := "created_at DESC"

	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	mockAccountRepo.On("Find", account.ID).Return(account, nil)
	mockTransactionRepo.On("FindByAccountToId", account.ID, pagination).Return([]*entity.Transaction{}, 0, errors.New("transaction empty"))

	accountTransaction := service.NewAccountTransaction(mockTransactionRepo, mockAccountRepo)

	result, total, err := accountTransaction.FindAllByAccountTo(account.ID, page, limit, sort)
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, total, 0)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "transaction empty")
}

func Test_FindOneByAccount_Success(t *testing.T) {
	mockTransactionRepo := new(repoMock.MockTransactionRepository)
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountFrom, _ := entity.NewAccount(3000.93)
	accountTo, _ := entity.NewAccount(200)
	amount := 200.00
	currency := "AKZ"
	transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, amount, currency)

	mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)
	mockTransactionRepo.On("FindOneByAccount", transaction.ID, accountTo.ID).Return(transaction, nil)

	accountTransaction := service.NewAccountTransaction(mockTransactionRepo, mockAccountRepo)

	result, err := accountTransaction.FindOneByAccount(accountTo.ID, transaction.ID)

	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, result.ID, transaction.ID)
	assert.Equal(t, result.AccountFromID, transaction.AccountFromID)
	assert.Equal(t, result.Amount, transaction.Amount)
	assert.Equal(t, result.AccountToID, transaction.AccountToID)
}

func Test_FindOneByAccount_Fail_On_FindAccount(t *testing.T) {
	mockAccountRepo := new(repoMock.MockAccountRepository)
	id := uuid.NewV4().String()
	transactionId := uuid.NewV4().String()

	mockAccountRepo.On("Find", id).Return(&entity.Account{}, errors.New("account not found"))

	accountTransaction := service.NewAccountTransaction(nil, mockAccountRepo)

	result, err := accountTransaction.FindOneByAccount(id, transactionId)
	mockAccountRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "account not found")

}

func Test_FindOneByAccount_Fail_On_FindTransaction(t *testing.T) {
	mockTransactionRepo := new(repoMock.MockTransactionRepository)
	mockAccountRepo := new(repoMock.MockAccountRepository)

	accountTo, _ := entity.NewAccount(3000.93)
	transactionId := uuid.NewV4().String()

	mockAccountRepo.On("Find", accountTo.ID).Return(accountTo, nil)
	mockTransactionRepo.On("FindOneByAccount", transactionId, accountTo.ID).Return(&entity.Transaction{}, errors.New("transaction not found"))

	accountTransaction := service.NewAccountTransaction(mockTransactionRepo, mockAccountRepo)

	result, err := accountTransaction.FindOneByAccount(accountTo.ID, transactionId)

	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "transaction not found")
}
