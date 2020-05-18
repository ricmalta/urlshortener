package service

import (
  "fmt"
  "github.com/ricmalta/urlshortner/internal/store"
  "net/http"
  "os"
)

type Service struct {
	urlStore   *store.Store
	httpServer *http.Server
	quitC      chan os.Signal
}

func NewService(httpServer *http.Server, urlStore *store.Store) (*Service, error) {
	return &Service{
		urlStore:   urlStore,
		httpServer: httpServer,
	}, nil
}

func (service *Service) Start() error {
	fmt.Printf("HTTP server started at %s", service.httpServer.Addr)
	if err := service.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

