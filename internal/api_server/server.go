package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	api "github.com/dcm-project/dcm-placement-api/api/v1alpha1"
	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/config"
	handlers "github.com/dcm-project/dcm-placement-api/internal/handlers/v1alpha1"
	"github.com/dcm-project/dcm-placement-api/internal/opa"
	"github.com/dcm-project/dcm-placement-api/internal/service"
	"github.com/dcm-project/dcm-placement-api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
	"go.uber.org/zap"
)

const (
	gracefulShutdownTimeout = 5 * time.Second
)

type Server struct {
	cfg          *config.Config
	store        store.Store
	listener     net.Listener
	opaValidator *opa.Validator
}

// New returns a new instance of a migration-planner server.
func New(
	cfg *config.Config,
	store store.Store,
	listener net.Listener,
	opaValidator *opa.Validator,
) *Server {
	return &Server{
		cfg:          cfg,
		store:        store,
		listener:     listener,
		opaValidator: opaValidator,
	}
}

func oapiErrorHandler(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, fmt.Sprintf("API Error: %s", message), statusCode)
}

func (s *Server) Run(ctx context.Context) error {
	zap.S().Named("api_server").Info("Initializing API server")
	swagger, err := api.GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to load swagger spec: %w", err)
	}
	// Skip server name validation
	swagger.Servers = nil

	oapiOpts := oapimiddleware.Options{
		ErrorHandler: oapiErrorHandler,
	}
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.Recoverer,
		oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapiOpts),
	)

	h := handlers.NewServiceHandler(
		service.NewHealthService(),
		service.NewPlacementService(),
	)
	server.HandlerFromMux(server.NewStrictHandler(h, nil), router)
	srv := http.Server{Addr: s.cfg.Service.Address, Handler: router}

	go func() {
		<-ctx.Done()
		zap.S().Named("api_server").Infof("Shutdown signal received: %s", ctx.Err())
		ctxTimeout, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		_ = srv.Shutdown(ctxTimeout)
	}()

	zap.S().Named("api_server").Infof("Listening on %s...", s.listener.Addr().String())
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}

	return nil
}
