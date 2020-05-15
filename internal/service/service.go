package service

import (
	"github.com/ricmalta/urlshortner/internal/config"
	"github.com/ricmalta/urlshortner/internal/store"
)

type Service struct {
	store *store.Store
}

func NewService(cfg config.Config) (*Service, error) {
	storeInstance, err := store.NewStore(cfg)
	if err != nil {
		return nil, err
	}

	return &Service{
		store: storeInstance,
	}, nil
}

func (service *Service) Start() {

}
