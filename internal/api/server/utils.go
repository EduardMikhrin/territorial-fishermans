package server

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

func validatePositiveID() validation.Rule {
	return validation.By(func(value interface{}) error {
		v, ok := value.(uint64)
		if !ok {
			return errors.New("value is not int64")
		}

		if v == 0 {
			return errors.New("value is must be more than 0")
		}

		return nil
	})
}
