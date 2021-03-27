package repository

import "github.com/EdlanioJ/kbu/payments/domain/entity"

type AccountRepository interface {
	Find(id string) (*entity.Account, error)
	Save(account *entity.Account) error
}
