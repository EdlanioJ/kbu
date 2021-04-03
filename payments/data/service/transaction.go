package service

import (
	"errors"

	"github.com/EdlanioJ/kbu/payments/data/repository"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
)

type Transaction struct {
	TransactionRepository repository.TransactionRepository
	AccountRepository     repository.AccountRepository
}

func NewTransaction(
	TransactionRepository repository.TransactionRepository,
	AccountRepository repository.AccountRepository,
) *Transaction {

	return &Transaction{
		TransactionRepository: TransactionRepository,
		AccountRepository:     AccountRepository,
	}
}

func (t *Transaction) Register(fromID, toID, externalID, transactionType, currency string, amount float64) (*entity.Transaction, error) {
	accountFrom, err := t.AccountRepository.Find(fromID)

	if err != nil {
		return nil, err
	}

	if accountFrom == nil {
		return nil, errors.New("no account from was found")
	}

	accountTo, err := t.AccountRepository.Find(toID)

	if err != nil {
		return nil, err
	}

	if accountTo == nil {
		return nil, errors.New("no account destination was found")
	}

	err = accountFrom.Withdow(amount)

	if err != nil {
		return nil, err
	}

	transaction, err := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

	if err != nil {
		return nil, err
	}

	err = t.TransactionRepository.Register(transaction)

	if err != nil {
		return nil, err
	}

	err = t.AccountRepository.Save(accountFrom)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) Find(id string) (*entity.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(id)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) FindAll(page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}
	transaction, total, err := t.TransactionRepository.FindAll(pagination)

	if err != nil {
		return nil, 0, err
	}

	return transaction, total, nil
}

func (t *Transaction) FindByType(transactionType, transactionID string) (*entity.Transaction, error) {
	transaction, err := t.TransactionRepository.FindByType(transactionID, transactionType)

	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t *Transaction) FindAllByType(transactionType string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	transactions, total, err := t.TransactionRepository.FindAllByType(transactionType, pagination)

	if err != nil {
		return nil, 0, err
	}
	return transactions, total, nil
}

func (t *Transaction) FindByExternalID(externalID, transactionID string) (*entity.Transaction, error) {
	transaction, err := t.TransactionRepository.FindByExternalID(transactionID, externalID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) FindAllByExternalID(externalID string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	transactions, total, err := t.TransactionRepository.FindAllByExternalID(externalID, pagination)

	if err != nil {
		return nil, 0, err
	}
	return transactions, total, nil
}

func (t *Transaction) FindAllByFromAccountID(accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	transactions, total, err := t.TransactionRepository.FindAllByFromAccountID(accountID, pagination)

	if err != nil {
		return nil, 0, err
	}
	return transactions, total, nil
}

func (t *Transaction) FindByFromAccountID(accountID, transactionID string) (*entity.Transaction, error) {
	transaction, err := t.TransactionRepository.FindByFromAccountID(transactionID, accountID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) FindAllByToAccountID(accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	transactions, total, err := t.TransactionRepository.FindAllByToAccountID(accountID, pagination)

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (t *Transaction) FindByToAccountID(accountID, transactionID string) (*entity.Transaction, error) {
	transaction, err := t.TransactionRepository.FindByToAccountID(transactionID, accountID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) Complete(transactionId string) (*entity.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionId)

	if err != nil {
		return nil, err
	}

	transaction.Status = entity.TransactionCompleted

	err = t.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return transaction, nil
}

func (t *Transaction) Error(transactionId string) (*entity.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	transaction.Status = entity.TransactionCanceled

	err = t.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return transaction, nil
}
