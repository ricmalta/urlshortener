package service

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ricmalta/urlshortner/internal/config"
	"github.com/ricmalta/urlshortner/internal/store"
)

type Service struct {
	urlStore   *store.Store
	httpServer *http.Server
	config     config.Config
	quitC      chan os.Signal
}

func NewService(cfg config.Config) (*Service, error) {
	urlStore, err := store.NewStore(cfg)
	if err != nil {
		return nil, err
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", cfg.HTTP.Port),
		Handler:           NewServiceHandler(urlStore),
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return &Service{
		urlStore:   urlStore,
		httpServer: httpServer,
		config:     cfg,
	}, nil
}

func (service *Service) Start() error {
	fmt.Printf("HTTP server started at port %d", service.config.HTTP.Port)
	if err := service.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

/* TODO
   Add URL
   Test add URL
   Get URL
   Test Get URL
   Graceful shutdown

*/
