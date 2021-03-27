package service

import (
	"errors"

	"github.com/EdlanioJ/kbu/payments/data/repository"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
)

type ServiceTransaction struct {
	AccountRepository      repository.AccountRepository
	TransactionRepository  repository.TransactionRepository
	ServiceRepository      repository.ServiceRepository
	ServicePriceRepository repository.ServicePriceRepository
}

func NewServiceTransaction(
	AccountRepository repository.AccountRepository,
	ServicePriceRepository repository.ServicePriceRepository,
	ServiceRepository repository.ServiceRepository,
	TransactionRepository repository.TransactionRepository,
) *ServiceTransaction {
	return &ServiceTransaction{
		AccountRepository:      AccountRepository,
		TransactionRepository:  TransactionRepository,
		ServiceRepository:      ServiceRepository,
		ServicePriceRepository: ServicePriceRepository,
	}
}

func (t *ServiceTransaction) RegisterServiceTransaction(fromId string, serviceId string, servicePriceId string, amount float64, currency string) (*entity.Transaction, error) {
	var finalAmount float64

	if servicePriceId == "" && amount == 0 {
		return nil, errors.New("missing amount or servicePriceId")
	}

	accountFrom, err := t.AccountRepository.Find(fromId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	service, err := t.ServiceRepository.FindServiceByIdAndStatus(serviceId, entity.ServiceActive)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	if servicePriceId != "" {
		servicePrice, err := t.ServicePriceRepository.Find(servicePriceId)

		if err != nil {
			return nil, errors.New(err.Error())
		}

		finalAmount = servicePrice.Amount
	} else {
		finalAmount = amount
	}

	err = accountFrom.Withdow(finalAmount)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	transaction, err := entity.NewTransaction(accountFrom, nil, service, nil, finalAmount, currency)

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
func (t *ServiceTransaction) FindAllByServiceId(serviceId string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	_, err := t.ServiceRepository.Find(serviceId)

	if err != nil {
		return nil, 0, errors.New(err.Error())
	}

	pagination := &entity.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	transaction, total, err := t.TransactionRepository.FindByServiceId(serviceId, pagination)

	if err != nil {
		return nil, 0, errors.New(err.Error())
	}
	return transaction, total, nil
}

func (t *ServiceTransaction) FindOneByService(serviceId string, transactionId string) (*entity.Transaction, error) {
	_, err := t.ServiceRepository.Find(serviceId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	transaction, err := t.TransactionRepository.FindOneByService(transactionId, serviceId)

	if err != nil {
		return nil, errors.New(err.Error())
	}
	return transaction, nil
}
