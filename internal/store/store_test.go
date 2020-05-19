package store_test

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ricmalta/urlshortner/internal/store"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

const (
	addValidURL   = "https://www.example.com/test"
	addInvalidURL = "--https://www.example.com/test"
)

var (
	cache         *lru.Cache
	mr            *miniredis.Miniredis
	redisClient   *redis.Client
	logger        *logrus.Logger
	storeInstance *store.Store
)

func init() {
	cache, err := lru.New(10)
	if err != nil {
		panic(err)
	}
	mr, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	logger, _ = test.NewNullLogger()
	storeInstance, err = store.NewStore(cache, redisClient, logger)
}

func TestCreateStoreInstance(t *testing.T) {
	storeInstance, err := store.NewStore(cache, redisClient, logger)
	assert.Nil(t, err, "return no error")
	assert.NotNil(t, storeInstance, "returns a valid store instance")
}

func TestAddInvalidURL(t *testing.T) {
	shortKey, err := storeInstance.Add(addInvalidURL)
	assert.NotNil(t, err, "return error")
	assert.Equal(t, store.ErrorInvalidInputURL{}, err, "error type of store.ErrorInvalidInputURL")
	assert.Equal(t, "", shortKey, "short key should be empty")

}

func TestAddValidURL(t *testing.T) {
	shortKey, err := storeInstance.Add(addValidURL)
	assert.Nil(t, err, "return no error")
	assert.NotEqual(t, "aa", shortKey, "short key should not be empty")
}

func TestGetExistingURL(t *testing.T) {}

func TestGetUnexistingURL(t *testing.T) {}
