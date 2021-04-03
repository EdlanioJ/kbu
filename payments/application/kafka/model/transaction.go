package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Transaction struct {
	ID          string  `json:"id"`
	AccountFrom string  `json:"account_from"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	AccountTo   string  `json:"account_to,omitempty"`
	Store       string  `json:"store,omitempty"`
	Service     string  `json:"service,omitempty"`
}

func (t *Transaction) isValid() error {
	err := validation.ValidateStruct(t,
		validation.Field(t.ID, validation.Required, is.UUIDv4),
		validation.Field(t.AccountFrom, validation.Required, is.UUIDv4),
		validation.Field(t.AccountTo, is.UUIDv4),
		validation.Field(t.Service, is.UUIDv4),
		validation.Field(t.Store, is.UUIDv4),
		validation.Field(t.Status, validation.Required),
		validation.Field(t.Amount, validation.Required, validation.Min(float64(0))),
	)

	return err
}

func (t *Transaction) ParseJson(data []byte) error {
	err := json.Unmarshal(data, t)

	if err != nil {
		return err
	}

	err = t.isValid()

	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) ToJson() ([]byte, error) {
	err := t.isValid()

	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(t)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewTransaction() *Transaction {
	return &Transaction{}
}
