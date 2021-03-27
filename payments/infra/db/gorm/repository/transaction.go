package repository

import (
	"errors"

	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/jinzhu/gorm"
)

type TransactionRepositoryGORM struct {
	DB *gorm.DB
}

func (t *TransactionRepositoryGORM) Register(transaction *entity.Transaction) error {
	err := t.DB.Omit("AccountFrom", "Service", "Store", "AccountTo").Create(transaction).Error

	if err != nil {
		return err
	}

	return nil
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepositoryGORM {
	return &TransactionRepositoryGORM{
		DB: db,
	}
}

func (t *TransactionRepositoryGORM) Save(transaction *entity.Transaction) error {
	err := t.DB.Omit("AccountFrom", "Service", "Store", "AccountTo").Save(transaction).Error

	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryGORM) Find(id string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}
	err := t.DB.First(transaction, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepositoryGORM) FindOneByAccount(transactionId string, accountId string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}

	err := t.DB.First(transaction, "id = ? AND account_to_id = ?", transactionId, accountId).Error

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepositoryGORM) FindOneByService(transactionId string, serviceId string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}

	err := t.DB.First(transaction, "id = ? AND service_id = ?", transactionId, serviceId).Error

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepositoryGORM) FindOneByStore(transactionId string, storeId string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}

	err := t.DB.First(transaction, "id = ? AND store_id = ?", transactionId, storeId).Error

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepositoryGORM) FindAll(pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int
	err := t.DB.
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&transactions).
		Count(&totalTransaction).
		Error

	if err != nil {
		return nil, 0, err
	}
	return transactions, totalTransaction, nil
}

func (t *TransactionRepositoryGORM) FindByAccountFromId(accountId string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int

	err := t.DB.
		Where("account_from_id = ?", accountId).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&transactions).
		Count(&totalTransaction).
		Error

	if err != nil {
		return nil, 0, err
	}
	return transactions, totalTransaction, nil
}

func (t *TransactionRepositoryGORM) FindByAccountToId(accountId string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int

	err := t.DB.
		Where("account_to_id = ?", accountId).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&transactions).
		Count(&totalTransaction).
		Error

	if err != nil {
		return nil, 0, err
	}
	return transactions, totalTransaction, nil
}

func (t *TransactionRepositoryGORM) FindByServiceId(serviceId string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int

	err := t.DB.
		Where("service_id = ?", serviceId).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&transactions).
		Count(&totalTransaction).
		Error

	if err != nil {
		return nil, 0, err
	}
	return transactions, totalTransaction, nil
}

func (t *TransactionRepositoryGORM) FindByStoreId(storeId string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int

	err := t.DB.
		Where("store_id = ?", storeId).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&transactions).
		Count(&totalTransaction).
		Error

	if err != nil {
		return nil, 0, errors.New(err.Error())
	}
	return transactions, totalTransaction, nil
}
