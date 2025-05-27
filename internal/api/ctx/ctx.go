package ctx

import (
	"context"
	"github.com/EduardMikhrin/territorial-fishermans/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	loggerKey ctxKey = iota
	dbKey
)

func LoggerProvider(logger *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, loggerKey, logger)
	}
}

func Logger(ctx context.Context) *logan.Entry {
	return ctx.Value(loggerKey).(*logan.Entry)
}

func DBProvider(q data.MasterQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {

		return context.WithValue(ctx, dbKey, q)
	}
}

func DB(ctx context.Context) data.MasterQ {
	return ctx.Value(dbKey).(data.MasterQ).New()
}
