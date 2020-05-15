package store

import (
	"github.com/ricmalta/urlshortner/internal/config"

	"github.com/go-redis/redis"
	lru "github.com/hashicorp/golang-lru"
)

type Store struct {
	URLs        map[string]string
	cache       *lru.Cache
	redisClient *redis.Client
}

func NewStore(cfg config.Config) (*Store, error) {
	cache, err := lru.New(cfg.Cache.Size)
	if err != nil {
		return nil, err
	}

	return &Store{
		URLs:  make(map[string]string),
		cache: cache,
		redisClient: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Host,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Database,
		}),
	}, nil
}

func (store *Store) Get() {}

func (store *Store) Add() {}

func (store *Store) generateKey() {}
