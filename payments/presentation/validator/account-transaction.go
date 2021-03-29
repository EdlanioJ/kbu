package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateRegisterAccountTransactionParams(
	fromId string,
	toId string,
	amount float64,
	currency string,
) error {
	err := validation.Errors{
		"from id":        validation.Validate(fromId, validation.Required, is.UUIDv4),
		"destination id": validation.Validate(toId, validation.Required, is.UUIDv4),
		"amount":         validation.Validate(amount, validation.Required, validation.Min(float64(0))),
		"currency":       validation.Validate(currency, is.CurrencyCode),
	}.Filter()

	return err
}

func ValidateFindAllByAccountToParams(accountId string, page int, limit int, sort string) error {
	err := validation.Errors{
		"account destination": validation.Validate(accountId, validation.Required, is.UUIDv4),
		"page":                validation.Validate(page, validation.Required, validation.Min(int(0))),
		"limit":               validation.Validate(limit, validation.Required, validation.Min(int(-1))),
		"sort":                validation.Validate(sort, validation.Required),
	}.Filter()

	return err
}

func ValidateFindOneByAccountParams(fromId string, transactionId string) error {
	err := validation.Errors{
		"accountFrom": validation.Validate(fromId, validation.Required, is.UUIDv4),
		"transaction": validation.Validate(transactionId, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}
