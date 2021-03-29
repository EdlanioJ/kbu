package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateRegisterServiceTransactionParams(fromId string, serviceId string, servicePriceId string, amount float64, currency string) error {
	err := validation.Errors{
		"accountFrom":  validation.Validate(fromId, validation.Required, is.UUIDv4),
		"service":      validation.Validate(serviceId, validation.Required, is.UUIDv4),
		"servicePrice": validation.Validate(servicePriceId, is.UUIDv4),
		"amount":       validation.Validate(amount, validation.Min(float64(0))),
		"currency":     validation.Validate(currency, validation.Required, is.CurrencyCode),
	}.Filter()

	return err
}

func ValidateFindAllByServiceIdParams(serviceId string, page int, limit int, sort string) error {
	err := validation.Errors{
		"service": validation.Validate(serviceId, validation.Required, is.UUIDv4),
		"page":    validation.Validate(page, validation.Required, validation.Min(int(0))),
		"limit":   validation.Validate(limit, validation.Required, validation.Min(int(-1))),
		"sort":    validation.Validate(sort, validation.Required),
	}.Filter()

	return err
}

func ValidateFindOneByServiceParams(serviceId string, transactionId string) error {
	err := validation.Errors{
		"service":     validation.Validate(serviceId, validation.Required, is.UUIDv4),
		"transaction": validation.Validate(transactionId, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}
