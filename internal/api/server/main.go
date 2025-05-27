package server

import (
	"context"
	docs "github.com/EduardMikhrin/territorial-fishermans/api"
	"github.com/EduardMikhrin/territorial-fishermans/internal/api/ctx"
	types "github.com/EduardMikhrin/territorial-fishermans/internal/api/types"
	"github.com/EduardMikhrin/territorial-fishermans/internal/data"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"

	"github.com/ignite/cli/ignite/pkg/openapiconsole"
	"gitlab.com/distributed_lab/logan/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"sync"
	"time"
)

var _ types.FishermanServiceServer = grpcImpl{}

type grpcImpl struct{}

type service struct {
	grpc net.Listener
	http net.Listener

	logger       *logan.Entry
	ctxExtenders []func(context.Context) context.Context
}

func NewService(listnerGRPC net.Listener, listenerHTTP net.Listener, log *logan.Entry, db data.MasterQ) types.Server {
	return service{
		grpc:   listnerGRPC,
		http:   listenerHTTP,
		logger: log,
		ctxExtenders: []func(context.Context) context.Context{
			ctx.LoggerProvider(log),
			ctx.DBProvider(db),
		},
	}
}

func (s service) RunGRPC(ctx context.Context, errChan chan<- error) {
	srv := s.grpcServer()
	once := sync.Once{}
	go func() {
		select {
		case <-ctx.Done():
			srv.GracefulStop()
			s.logger.Info("grpc serving stopped: context canceled")
			once.Do(func() {
				errChan <- nil
			})
		}
	}()

	s.logger.Info("grpc serving started")
	err := srv.Serve(s.grpc)

	once.Do(func() {
		errChan <- errors.Wrap(err, "grpc serving error occurred")
	})
}

func (s service) RunHTTP(ctx context.Context, errChan chan<- error) {
	handler, err := s.httpServer(ctx)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to start http server")
		return
	}
	srv := &http.Server{Handler: handler}
	once := sync.Once{}

	go func() {
		select {
		case <-ctx.Done():
			shutdownDeadline, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			if err = srv.Shutdown(shutdownDeadline); err != nil {
				once.Do(func() {
					errChan <- errors.Wrap(err, "failed to shutdown http server")
				})
				return
			}
			s.logger.Info("http serving stopped: context canceled")
			once.Do(func() {
				errChan <- nil
			})
		}
	}()

	s.logger.Info("http serving started")
	if err = srv.Serve(s.http); !errors.Is(err, http.ErrServerClosed) {
		once.Do(func() {
			errChan <- errors.Wrap(err, "failed to start http server")
		})
		return
	}

	once.Do(func() {
		errChan <- nil
	})

}

func (s service) httpServer(ctx context.Context) (http.Handler, error) {
	router := chi.NewRouter()
	router.Use(ape.LoganMiddleware(s.logger), ape.RecoverMiddleware(s.logger))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	grpcRouter := runtime.NewServeMux()
	if err := types.RegisterFishermanServiceHandlerServer(ctx, grpcRouter, grpcImpl{}); err != nil {
		return nil, errors.Wrap(err, "failed to register user operation service handler")
	}

	router.With(ape.CtxMiddleware(s.ctxExtenders...)).Mount("/", grpcRouter)
	router.Mount("/static/api.swagger.json", http.FileServer(http.FS(docs.Docs)))
	router.HandleFunc("/api", openapiconsole.Handler("Signer service", "/static/api.swagger.json"))

	return router, nil

}

func (s service) grpcServer() *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			ContextExtenderInterceptor(s.ctxExtenders...),
			LoggerInterceptor(s.logger),
			// RecoveryInterceptor should be the last one
			RecoveryInterceptor(s.logger),
		),
	)

	types.RegisterFishermanServiceServer(srv, grpcImpl{})
	reflection.Register(srv)

	return srv
}
