package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	ServiceActive  string = "active"
	ServiceDesable string = "disable"
	ServicePending string = "pending"
)

type Service struct {
	Base         `valid:"required"`
	Name         string          `json:"name" gorm:"type:varchar(255)" valid:"notnull"`
	Description  string          `json:"description" gorm:"type:varchar(255)" valid:"-"`
	FromID       string          `json:"from" gorm:"type:uuid" valid:"notnull,uuidv4"`
	TypeID       string          `json:"type" gorm:"type:uuid" valid:"notnull,uuidv4"`
	Status       string          `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	ServicePrice *[]ServicePrice `gorm:"foreignkey:ServiceID" valid:"-"`
	Transaction  *[]Transaction  `gorm:"foreignkey:ServiceID" valid:"-"`
}

func (s *Service) isValid() error {
	_, err := govalidator.ValidateStruct(s)

	if s.Status != ServiceActive && s.Status != ServicePending && s.Status != ServiceDesable {
		return errors.New("invalid status")
	}

	if err != nil {
		return err
	}
	return nil
}

func NewService(name string, description string, fromId string, typeId string) (*Service, error) {
	service := Service{
		Name:        name,
		Description: description,
		FromID:      fromId,
		TypeID:      typeId,
		Status:      ServicePending,
	}
	service.ID = uuid.NewV4().String()
	service.CreatedAt = time.Now()

	err := service.isValid()

	if err != nil {
		return nil, err
	}

	return &service, nil

}
