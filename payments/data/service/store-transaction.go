package service

import (
	"errors"

	"github.com/EdlanioJ/kbu/payments/data/repository"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
)

type StoreTransaction struct {
	AccountRepository     repository.AccountRepository
	StoreRepository       repository.StoreRepository
	TransactionRepository repository.TransactionRepository
}

func NewStoreTransaction(
	AccountRepository repository.AccountRepository,
	StoreRepository repository.StoreRepository,
	TransactionRepository repository.TransactionRepository,
) *StoreTransaction {
	return &StoreTransaction{
		AccountRepository:     AccountRepository,
		StoreRepository:       StoreRepository,
		TransactionRepository: TransactionRepository,
	}
}

func (t *StoreTransaction) RegisterStoreTransaction(fromAccountId string, storeId string, amount float64, currency string) (*entity.Transaction, error) {
	accountFrom, err := t.AccountRepository.Find(fromAccountId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	store, err := t.StoreRepository.FindStoreByIdAndStatus(storeId, entity.StoreActive)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = accountFrom.Withdow(amount)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	transaction, err := entity.NewTransaction(accountFrom, nil, nil, store, amount, currency)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = t.TransactionRepository.Register(transaction)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = t.AccountRepository.Save(accountFrom)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return transaction, nil
}

func (t *StoreTransaction) FindAllByStoreId(storeId string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	_, err := t.StoreRepository.Find(storeId)

	if err != nil {
		return nil, 0, errors.New(err.Error())
	}

	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}
	transactions, total, err := t.TransactionRepository.FindByStoreId(storeId, pagination)

	if err != nil {
		return nil, 0, errors.New(err.Error())
	}

	return transactions, total, nil
}
func (t *StoreTransaction) FindOneByStore(storeId string, transactionId string) (*entity.Transaction, error) {
	_, err := t.StoreRepository.Find(storeId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	transaction, err := t.TransactionRepository.FindOneByStore(transactionId, storeId)

	if err != nil {
		return nil, errors.New(err.Error())
	}
	return transaction, nil
}
