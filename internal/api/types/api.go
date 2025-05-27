package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal = status.Error(codes.Internal, "internal server error")
	ErrStopped  = errors.New("service stopped")
)

type Server interface {
	RunGRPC(ctx context.Context, errCh chan<- error)
	RunHTTP(ctx context.Context, errCh chan<- error)
}
