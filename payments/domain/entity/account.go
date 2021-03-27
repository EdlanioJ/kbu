package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Base    `valid:"required"`
	Balance float64 `json:"balance" gorm:"type:float" valid:"-"`
}

func (a *Account) isValid() error {
	_, err := govalidator.ValidateStruct(a)

	if err != nil {
		return err
	}
	return nil
}

func NewAccount(balance float64) (*Account, error) {
	account := Account{
		Balance: balance,
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()
	err := account.isValid()

	if err != nil {
		return nil, err
	}

	return &account, nil
}
func (a *Account) Deposit(amount float64) error {
	a.Balance += amount

	err := a.isValid()

	if err != nil {
		return err
	}

	return nil
}

func (a *Account) Withdow(amount float64) error {
	if a.Balance < amount {
		return errors.New("account does not have balance")
	}

	a.Balance -= amount

	err := a.isValid()

	if err != nil {
		return err
	}
	return nil
}
