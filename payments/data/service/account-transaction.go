package service

import (
	"errors"

	"github.com/EdlanioJ/kbu/payments/data/repository"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
)

type AccountTransaction struct {
	TransactionRepository repository.TransactionRepository
	AccountRepository     repository.AccountRepository
}

func NewAccountTransaction(TransactionRepository repository.TransactionRepository, AccountRepository repository.AccountRepository) *AccountTransaction {
	return &AccountTransaction{
		TransactionRepository: TransactionRepository,
		AccountRepository:     AccountRepository,
	}
}

func (t *AccountTransaction) RegisterAccountTransaction(fromId string, toId string, amount float64, currency string) (*entity.Transaction, error) {
	accountFrom, err := t.AccountRepository.Find(fromId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	accountTo, err := t.AccountRepository.Find(toId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = accountFrom.Withdow(amount)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	transaction, err := entity.NewTransaction(accountFrom, accountTo, nil, nil, amount, currency)

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

func (t *AccountTransaction) FindAllByAccountTo(accountId string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	_, err := t.AccountRepository.Find(accountId)

	if err != nil {
		return nil, 0, errors.New(err.Error())
	}

	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	transaction, total, err := t.TransactionRepository.FindByAccountToId(accountId, pagination)

	if err != nil {
		return nil, 0, errors.New(err.Error())
	}
	return transaction, total, nil
}

func (t *AccountTransaction) FindOneByAccount(accountId string, transactionId string) (*entity.Transaction, error) {
	_, err := t.AccountRepository.Find(accountId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	transaction, err := t.TransactionRepository.FindOneByAccount(transactionId, accountId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return transaction, nil
}
