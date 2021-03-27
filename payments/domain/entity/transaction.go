package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
)

type Transaction struct {
	Base          `valid:"required"`
	Amount        float64  `json:"amount" gorm:"type:float" valid:"notnull"`
	Status        string   `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Currency      string   `json:"currency" gorm:"type:varchar(5)" valid:"notnull"`
	AccountFrom   *Account `valid:"-"`
	AccountFromID string   `json:"account_from,omitempty" gorm:"column:account_from_id;type:uuid;not null" valid:"notnull,uuidv4"`
	Service       *Service `valid:"-"`
	ServiceID     string   `json:"service,omitempty" gorm:"column:service_id;type:uuid;default:null" valid:"optional,uuidv4"`
	Store         *Store   `valid:"-"`
	StoreID       string   `json:"store,omitempty" gorm:"column:store_id;type:uuid;default:null" valid:"optional,uuidv4"`
	AccountTo     *Account `valid:"-"`
	AccountToID   string   `json:"account_to,omitempty" gorm:"column:account_to_id;type:uuid;default:null" valid:"optional,uuidv4"`
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if err != nil {
		return err
	}

	if t.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("invalid status")
	}
	return nil
}

func NewTransaction(accountFrom *Account, accountTo *Account, service *Service, store *Store, amount float64, currency string) (*Transaction, error) {
	if currency == "" {
		currency = "AKZ"
	}

	transaction := Transaction{
		Amount:        amount,
		Currency:      currency,
		AccountFrom:   accountFrom,
		AccountFromID: accountFrom.ID,
		AccountTo:     accountTo,
		Service:       service,
		Store:         store,
		Status:        TransactionPending,
	}

	if accountTo != nil {
		transaction.AccountToID = accountTo.ID
	}

	if service != nil {
		transaction.ServiceID = service.ID
	}

	if store != nil {
		transaction.StoreID = store.ID
	}
	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
