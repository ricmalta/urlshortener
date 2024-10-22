package service

import (
	"context"
	"net/http"

	"github.com/ricmalta/urlshortner/internal/store"
	"github.com/sirupsen/logrus"
)

type Service struct {
	urlStore   *store.Store
	httpServer *http.Server
	logger     *logrus.Logger
}

func NewService(httpServer *http.Server, urlStore *store.Store, logger *logrus.Logger) (*Service, error) {
	return &Service{
		urlStore:   urlStore,
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

func (service *Service) Start() error {
	service.logger.Infof("HTTP server started at %s", service.httpServer.Addr)
	if err := service.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (service *Service) Shutdown(ctx context.Context) error {
	service.logger.Info("HTTP server shutting down gracefully")
	return service.httpServer.Shutdown(ctx)
}
