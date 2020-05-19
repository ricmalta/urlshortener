package main

import (
  "flag"
  "fmt"
  "github.com/go-redis/redis"
  lru "github.com/hashicorp/golang-lru"
  "github.com/ricmalta/urlshortner/internal/logger"
  "github.com/ricmalta/urlshortner/internal/config"
	"github.com/ricmalta/urlshortner/internal/service"
  "github.com/ricmalta/urlshortner/internal/store"
  "net/http"
  "time"
)

func main() {
  configPath := flag.String("config", "./internal/config/config.yaml", "config file path")
  flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err.Error())
	}

  logger, err := logger.New(cfg.LogLevel)
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
  urlStore, err := store.NewStore(cache, redisClient, logger)
  if err != nil {
    panic(err)
  }

  serviceHandler := service.NewHandler(urlStore, logger, cfg.Service.BaseURL)

  // HTTP server instance
  httpServer := &http.Server{
    Addr:              fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
    Handler:           serviceHandler.Router,
    ReadTimeout:       1 * time.Second,
    WriteTimeout:      1 * time.Second,
    IdleTimeout:       30 * time.Second,
    ReadHeaderTimeout: 2 * time.Second,
  }

  // Short URL service instance
	service, err := service.NewService(httpServer, urlStore, logger)
	if err != nil {
		panic(err.Error())
	}

	// Start service
	service.Start()
}
