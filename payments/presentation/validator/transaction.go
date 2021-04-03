package validator

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func RegisterParams(accountFrom, accountTo, externalID, transactionType, currency string, amount float64) error {

	err := validation.Errors{
		"account_from": validation.Validate(accountFrom, validation.Required, is.UUIDv4),
		"account_to":   validation.Validate(accountTo, validation.Required, is.UUIDv4),
		"reference_id": validation.Validate(externalID, validation.Required, is.UUIDv4),
		"type": validation.Validate(transactionType, validation.Required, validation.In(
			entity.TransactionToService,
			entity.TransactionToStore,
			entity.TransactionToUser,
		)),
		"currency": validation.Validate(currency, validation.Required, is.CurrencyCode),
		"amount":   validation.Validate(amount, validation.Required, validation.Min(float64(0))),
	}.Filter()

	return err
}

func GetParams(id string) error {

	err := validation.Errors{
		"id": validation.Validate(id, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}

func GetAllParams(page int, limit int, sort string) error {
	err := validation.Errors{
		"page":  validation.Validate(page, validation.Required, validation.Min(int(0))),
		"limit": validation.Validate(limit, validation.Required),
		"sort":  validation.Validate(sort, validation.Required),
	}.Filter()

	return err
}

func GetByTypeParams(transactionID, transactionType string) error {
	err := validation.Errors{
		"transaction_id": validation.Validate(transactionID, validation.Required, is.UUIDv4),
		"type": validation.Validate(transactionType, validation.Required, validation.In(
			entity.TransactionToService,
			entity.TransactionToStore,
			entity.TransactionToUser,
		)),
	}.Filter()

	return err
}

func ListByTypeParams(transactionType string, page, limit int, sort string) error {
	err := validation.Errors{
		"type": validation.Validate(transactionType, validation.Required, validation.In(
			entity.TransactionToService,
			entity.TransactionToStore,
			entity.TransactionToUser,
		)),
		"page":  validation.Validate(page, validation.Required, validation.Min(int(0))),
		"limit": validation.Validate(limit, validation.Required, validation.Min(int(-1))),
		"sort":  validation.Validate(sort, validation.Required),
	}.Filter()

	return err
}

func GetByExternalIDParams(transactionID, externalID string) error {
	err := validation.Errors{
		"transaction_id": validation.Validate(transactionID, validation.Required, is.UUIDv4),
		"reference_id":   validation.Validate(externalID, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}

func ListByExternalIDParams(externalID string, page, limit int, sort string) error {
	err := validation.Errors{
		"reference_id": validation.Validate(externalID, validation.Required, is.UUIDv4),
		"page":         validation.Validate(page, validation.Required, validation.Min(int(0))),
		"limit":        validation.Validate(limit, validation.Required, validation.Min(int(-1))),
		"sort":         validation.Validate(sort, validation.Required),
	}.Filter()

	return err
}

func GetByAccountFromParams(transactionID, accountID string) error {
	err := validation.Errors{
		"transaction_id": validation.Validate(transactionID, validation.Required, is.UUIDv4),
		"account_id":     validation.Validate(accountID, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}

func ListByAccountFromParams(accountID string, page, limit int, sort string) error {
	err := validation.Errors{
		"account_id": validation.Validate(accountID, validation.Required, is.UUIDv4),
		"page":       validation.Validate(page, validation.Required, validation.Min(int(0))),
		"limit":      validation.Validate(limit, validation.Required, validation.Min(int(-1))),
		"sort":       validation.Validate(sort, validation.Required),
	}.Filter()

	return err
}

func GetByAccoutToParams(transactionID, accountID string) error {
	err := validation.Errors{
		"transaction_id": validation.Validate(transactionID, validation.Required, is.UUIDv4),
		"account_id":     validation.Validate(accountID, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}

func ListByAccoutToParams(accountID string, page, limit int, sort string) error {
	err := validation.Errors{
		"account_id": validation.Validate(accountID, validation.Required, is.UUIDv4),
		"page":       validation.Validate(page, validation.Required, validation.Min(int(0))),
		"limit":      validation.Validate(limit, validation.Required, validation.Min(int(-1))),
		"sort":       validation.Validate(sort, validation.Required),
	}.Filter()

	return err
}

func CompleteParams(id string) error {
	err := validation.Errors{
		"transaction": validation.Validate(id, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}

func ErrorParams(id string) error {
	err := validation.Errors{
		"transaction": validation.Validate(id, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}
