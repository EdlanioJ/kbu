package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type ServicePrice struct {
	Base        `valid:"required"`
	Description string   `json:"description" gorm:"type:varchar(255)" valid:"-"`
	Amount      float64  `json:"amount" gorm:"type:float" valid:"notnull"`
	Currency    string   `json:"currency" gorm:"type:varchar(5)" valid:"notnull"`
	Service     *Service `valid:"-"`
	ServiceID   string   `gorm:"column:service_id;type:uuid;not null" valid:"notnull,uuidv4"`
}

func (s *ServicePrice) isValid() error {
	_, err := govalidator.ValidateStruct(s)

	if s.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}
	if err != nil {
		return err
	}

	return nil
}

func NewServicePrice(service *Service, description string, amount float64, currency string) (*ServicePrice, error) {
	if currency == "" {
		currency = "AOA"
	}
	servicePrice := ServicePrice{
		Service:     service,
		ServiceID:   service.ID,
		Description: description,
		Amount:      amount,
		Currency:    currency,
	}

	servicePrice.ID = uuid.NewV4().String()
	servicePrice.CreatedAt = time.Now()

	err := servicePrice.isValid()

	if err != nil {
		return nil, err
	}

	return &servicePrice, nil
}
