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
	TransactionCanceled  string = "canceled"

	TransactionToUser    string = "to_user"
	TransactionToService string = "to_service"
	TransactionToStore   string = "to_store"
)

type Transaction struct {
	Base          `valid:"required"`
	Amount        float64  `json:"amount" gorm:"type:float" valid:"notnull"`
	Status        string   `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Currency      string   `json:"currency" gorm:"type:varchar(5)" valid:"notnull"`
	AccountFrom   *Account `valid:"-"`
	AccountFromID string   `json:"account_from" gorm:"column:account_from_id;type:uuid;not null" valid:"notnull,uuidv4"`
	AccountTo     *Account `valid:"-"`
	AccountToID   string   `json:"account_to" gorm:"column:account_to_id;type:uuid;default:null" valid:"notnull,uuidv4"`
	Type          string   `json:"type" gorm:"type:varchar(30)" valid:"notnull"`
	ExternalID    string   `json:"external_id" gorm:"column:external_id;type:uuid" valid:"notnull,uuidv4"`
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if err != nil {
		return err
	}

	if t.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if t.Type != TransactionToUser && t.Type != TransactionToService && t.Type != TransactionToStore {
		return errors.New("invalid type transaction")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionCanceled {
		return errors.New("invalid status")
	}
	return nil
}

func NewTransaction(accountFrom *Account, accountTo *Account, externalID, transactionType, currency string, amount float64) (*Transaction, error) {
	if currency == "" {
		currency = "AOA"
	}

	transaction := Transaction{
		Amount:        amount,
		Currency:      currency,
		AccountFrom:   accountFrom,
		AccountFromID: accountFrom.ID,
		AccountTo:     accountTo,
		AccountToID:   accountTo.ID,
		ExternalID:    externalID,
		Type:          transactionType,
		Status:        TransactionPending,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
