package pg

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parsePgError(err error) error {

	var pgerr *pq.Error
	if !errors.As(err, &pgerr) {
		return status.Errorf(codes.Internal, "DB error: %v", err)
	}

	switch pgerr.Code {
	case pq.ErrorCode("23505"):
		return errors.Wrap(ErrUnique, pgerr.Constraint)

	case pq.ErrorCode("23503"):
		return errors.Wrap(InvalidForeignKeyError, pgerr.Constraint)
	}

	return err
}
