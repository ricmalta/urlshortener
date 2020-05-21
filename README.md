# Simple URL Shortner

This simple URL shortener service is built around to demonstrate a scalable approach around Redis. Its build with Golang and using the as minimum dependencies as time had permitted.

## TL;DR

To try the service locally, ensure you have Docker and Docker compose installed and run:

```shell
$ docker-compose up
```

The service is running on TCP port 3000 by default.

## HTTP service

It's HTTP JSON based API with two endpoints:

1. POST / to create a new short URL.
2. GET /{short_key} to retrieve the redirect from the short URL to the original URL.

## Implemenation details

The implementation is made around two main design decisions:

1. To create a new short URL, it's used a sequential int64 key, based on a Redis counter, which represents the real limit of the possible samples, and then converting the sequential base 10 results in a base 36 number to use as the original URL key in a key-value format. The base 36 representation is more of a way to ensure the shortness of the URL because it's not possible to go further than the int64 sample limit.
2. To redirect the short URL to the original URL the following, the first step is to try the LRU cache, if the key is not in memory, then it´s fetched from Redis. There always the change of the key to not exist in neither in the cache nor the Redis key-value store.

The LRU cache is configurable in size, to decide the size of each running node in a production environment.

The Redis counter is a conscious limitation. This specific solution base code is hard limited to the size of this storage type limitation, which is an int64.

The project was build using Go 1.14 and go mod vendoring.

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

The main package accepts a -config flag to specify the configuration file, by default is the one in the internal config package. Being Redis a dependency, to run locally, you need a Redis server running.

```shell
$ go run main.go
```

```shell
$ go run main.go -config path_to_config_file.yaml
```

## Unit tests

There are unit tests for the config, service and store packages.

```shell
$ go test ./...
```
