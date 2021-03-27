package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	StoreActive  string = "active"
	StoreDesable string = "disable"
	StorePending string = "pending"
)

type Store struct {
	Base        `valid:"required"`
	Name        string         `json:"name" gorm:"type:varchar(255)" valid:"notnull"`
	Description string         `json:"description" gorm:"type:varchar(255)" valid:"-"`
	FromID      string         `valid:"notnull,uuidv4"`
	TypeID      string         `json:"type" gorm:"type:uuid" valid:"notnull,uuidv4"`
	Status      string         `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Transaction *[]Transaction `gorm:"foreignkey:StoreID" valid:"-"`
}

func (s *Store) isValid() error {
	_, err := govalidator.ValidateStruct(s)

	if s.Status != StoreActive && s.Status != StorePending && s.Status != ServiceDesable {
		return errors.New("invalid status")
	}

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func NewStore(name string, description string, fromId string, typeId string) (*Store, error) {
	store := Store{
		Name:        name,
		Description: description,
		FromID:      fromId,
		TypeID:      typeId,
		Status:      StorePending,
	}

	store.ID = uuid.NewV4().String()
	store.CreatedAt = time.Now()

	err := store.isValid()
	if err != nil {
		return nil, err
	}

	return &store, nil
}
