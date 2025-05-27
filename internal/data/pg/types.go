package pg

import "github.com/pkg/errors"

var (
	InvalidForeignKeyError = errors.New("invalid foreign key")
	ErrUnique              = errors.New("unique constraint violation")
)
