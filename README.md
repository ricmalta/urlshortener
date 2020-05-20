# Simple URL Shortner

This simple URL shortner service is build arround to demostrate a scalable approach arround Redis. Its build with Golang and using the as minimum dependencies as time had permited.

## TL;DR

To try the service locally, ensure you haver Docker and Docker compose installed and run:

```shell
$ docker-compose up
```

The service is running on TCP port 3000 by default.

## HTTP service

It's HTTP JSON based API with two endpoints:

1. POST / to create a new short URL.
2. GET /{short_key} to retreive the redirect from the short URL to the original URL.

## Implemenation details

The implementation is made arround two main design decisions:

1. To create a new short URL, it's used a sequencial int64 key, based on a Redis counter, which represents the real limit of the possible samples, and then converting the sequencial base 10 result in a base 36 number to use as the original URL key in a key-value format. The base 36 representation is more of a way to ensure the shortness of the URL, because it's not possible to go further than the int64 sample limit.
2. To redirect the short URL to the orignal URL the following, the first step is to try the LRU cache, if the key is not in memory, then it´s fetched from Redis. There always the change of the key to not exist in neither in the cache nor the Redis key value store.

The LRU cache is configurable in size, in order to decide the size of each running node in a production environment.

The Redis counter is a conscious limitation. This specific solution basecode is hard limited to the size of the this storage type limitation, which is a int64.

The project was build using Go 1.14 and go mod vendoring

## Repo relevant files and directories

```shell
/docker -> configuration to be used by docker compose
  /internal -> all the project private packages
  /config -> configuration package and default configuartion file
  /logger -> basic logging package to use cross package
  /service -> HTTP service interface and handlers
  /store -> URL store and biz logic implementation
/Dockerfile -> Docker configuration with the optimized image for runtime
/main.go -> main package

```

## Run and environment details

The main package accepts a -config flag to specify the configuration file, by default is the one in the internal config package. Being Redis a dependicy, to run locally, you need a Redis server running.

```shell
$ go run main.go
```

```shell
$ go run main.go -config path_to_config_file.yaml
```

## Unit tests

There is unit tests for the config, service and store packages.

```shell
$ go test ./...
```
