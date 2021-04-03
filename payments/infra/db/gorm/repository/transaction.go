package repository

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/jinzhu/gorm"
)

type TransactionRepositoryGORM struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepositoryGORM {
	return &TransactionRepositoryGORM{
		DB: db,
	}
}

func (t *TransactionRepositoryGORM) Register(transaction *entity.Transaction) error {
	err := t.DB.Omit("AccountFrom", "Service", "Store", "AccountTo").Create(transaction).Error

	if err != nil {
		return err
	}

	return nil
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

func (t *TransactionRepositoryGORM) FindByType(transactionID, transactionType string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}

	err := t.DB.First(transaction, "id = ? AND type = ?", transactionID, transactionType).Error

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepositoryGORM) FindAllByType(transactionType string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int
	err := t.DB.
		Where("type = ?", transactionType).
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

func (t *TransactionRepositoryGORM) FindByExternalID(transactionID, externalID string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}

	err := t.DB.First(transaction, "id = ? AND external_id = ?", transactionID, externalID).Error

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepositoryGORM) FindAllByExternalID(externalID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int

	err := t.DB.
		Where("external_id = ?", externalID).
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

func (t *TransactionRepositoryGORM) FindByFromAccountID(transactionID, accountID string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}

	err := t.DB.First(transaction, "id = ? AND account_from_id = ?", transactionID, accountID).Error

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepositoryGORM) FindAllByFromAccountID(accountID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int

	err := t.DB.
		Where("account_from_id = ?", accountID).
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

func (t *TransactionRepositoryGORM) FindByToAccountID(transactionID, accountID string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}
	err := t.DB.First(transaction, "id = ? AND account_to_id = ?", transactionID, accountID).Error

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepositoryGORM) FindAllByToAccountID(accountID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction

	limit := pagination.Limit
	sort := pagination.Sort
	page := pagination.Page

	var totalTransaction int

	err := t.DB.
		Where("account_to_id = ?", accountID).
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
