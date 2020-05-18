package main

import (
  "fmt"
  "github.com/go-redis/redis"
  lru "github.com/hashicorp/golang-lru"
  "github.com/ricmalta/urlshortner/internal/config"
	"github.com/ricmalta/urlshortner/internal/service"
  "github.com/ricmalta/urlshortner/internal/store"
  "net/http"
  "time"
)

func main() {
	cfg, err := config.LoadConfig("./internal/config/config.yaml")
	if err != nil {
		panic(err.Error())
	}

	// LRU cache instance
  cache, err := lru.New(cfg.Cache.Size)
  if err != nil {
    panic(err)
  }

  // Redis client
  redisClient := redis.NewClient(&redis.Options{
    Addr:     cfg.Redis.Host,
    Password: cfg.Redis.Password,
    DB:       cfg.Redis.Database,
  })
  // check redis connection
  if _, err := redisClient.Ping().Result(); err != nil {
    panic(err)
  }

  // URL store Instance
  urlStore, err := store.NewStore(cache, redisClient)
  if err != nil {
    panic(err)
  }

    // HTTP server instance
  httpServer := &http.Server{
    Addr:              fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
    Handler:           service.NewServiceHandler(urlStore, cfg.Service.BaseURL),
    ReadTimeout:       1 * time.Second,
    WriteTimeout:      1 * time.Second,
    IdleTimeout:       30 * time.Second,
    ReadHeaderTimeout: 2 * time.Second,
  }

  // Short URL service instance
	service, err := service.NewService(httpServer, urlStore)
	if err != nil {
		panic(err.Error())
	}

	// Start service
	service.Start()
}
