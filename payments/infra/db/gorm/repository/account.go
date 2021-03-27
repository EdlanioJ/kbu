package repository

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/jinzhu/gorm"
)

type AccountRepositoryGORM struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepositoryGORM {
	return &AccountRepositoryGORM{
		DB: db,
	}
}

func (a *AccountRepositoryGORM) Find(id string) (*entity.Account, error) {
	account := &entity.Account{}

	err := a.DB.First(account, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return account, nil
}
func (a *AccountRepositoryGORM) Save(account *entity.Account) error {
	err := a.DB.Save(account).Error

	if err != nil {
		return err
	}

	return nil
}
