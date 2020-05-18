package main

import (
	"github.com/ricmalta/urlshortner/internal/config"
	"github.com/ricmalta/urlshortner/internal/service"
)

func main() {
	cfg, err := config.LoadConfig("./internal/config")
	if err != nil {
		panic(err.Error())
	}

	service, err := service.NewService(cfg)
	if err != nil {
		panic(err.Error())
	}

	service.Start()
}
