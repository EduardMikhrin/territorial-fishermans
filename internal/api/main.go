package api

import (
	"context"
	"github.com/EduardMikhrin/territorial-fishermans/internal/api/server"
	"github.com/EduardMikhrin/territorial-fishermans/internal/api/types"
	"github.com/EduardMikhrin/territorial-fishermans/internal/config"
	"github.com/EduardMikhrin/territorial-fishermans/internal/data/pg"
	"github.com/pkg/errors"

	"sync"
)

const (
	serviceComponent = "component"
	componentServer  = "server"
)

func RunServer(ctx context.Context, cfg config.Config) error {
	logger := cfg.Log()
	wg := new(sync.WaitGroup)

	srv := server.NewService(
		cfg.GRPCListener(),
		cfg.HTTPListener(),
		logger.WithField(serviceComponent, componentServer),
		pg.NewMasterQ(cfg.DB()),
	)

	wg.Add(2)
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		srv.RunHTTP(ctx, errChan)
	}()

	go func() {
		defer wg.Done()
		srv.RunGRPC(ctx, errChan)
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return errors.Wrap(err, "error running server")
		}
	}

	return grpc.ErrStopped
}
