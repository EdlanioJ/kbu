package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateFindParams(id string) error {

	err := validation.Errors{
		"id": validation.Validate(id, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}

func ValidateFindAllParams(page int, limit int, sort string) error {
	err := validation.Errors{
		"page":  validation.Validate(page, validation.Required, validation.Min(int(0))),
		"limit": validation.Validate(limit, validation.Required),
		"sort":  validation.Validate(sort, validation.Required),
	}.Filter()

	return err
}

func ValidateCompleteParams(id string) error {
	err := validation.Errors{
		"transaction": validation.Validate(id, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}

func ValidateErrorParams(id string) error {
	err := validation.Errors{
		"transaction": validation.Validate(id, validation.Required, is.UUIDv4),
	}.Filter()

	return err
}
