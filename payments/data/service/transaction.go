package service

import (
	"errors"

	"github.com/EdlanioJ/kbu/payments/data/repository"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
)

type Transaction struct {
	TransactionRepository repository.TransactionRepository
}

func NewTransaction(TransactionRepository repository.TransactionRepository) *Transaction {

	return &Transaction{
		TransactionRepository: TransactionRepository,
	}

}

func (t *Transaction) Find(id string) (*entity.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(id)

	if err != nil {
		return nil, errors.New(err.Error())
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
		return nil, 0, errors.New(err.Error())
	}

	return transaction, total, nil
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

	transaction.Status = entity.TransactionError

	err = t.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return transaction, nil
}
