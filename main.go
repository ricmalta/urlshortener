package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ricmalta/urlshortner/internal/config"
	"github.com/ricmalta/urlshortner/internal/logger"
	"github.com/ricmalta/urlshortner/internal/service"
	"github.com/ricmalta/urlshortner/internal/store"
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
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}

	// Short URL service instance
	service, err := service.NewService(httpServer, urlStore, logger)
	if err != nil {
		panic(err.Error())
	}

	// Start service
	go func() {
		if err := service.Start(); err != nil {
			logger.Fatal(err)
		}
	}()

	// create quit channel and wait until receive the process interrupt
	quitC := make(chan os.Signal, 1)
	signal.Notify(quitC, os.Interrupt, syscall.SIGTERM)
	<-quitC

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := service.Shutdown(ctx); err != nil {
		logger.Fatalf("HTTP server shutdown failed '%v'", err)
	}
}
