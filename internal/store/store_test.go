package store_test

import (
  "github.com/alicebob/miniredis"
  "github.com/go-redis/redis"
  lru "github.com/hashicorp/golang-lru"
  "github.com/ricmalta/urlshortner/internal/store"
  "github.com/stretchr/testify/assert"
  "testing"
)

const addValidURL = "https://www.example.com/test"
const addInvalidURL = "--https://www.example.com/test"

func TestCreateStoreInstance(t *testing.T) {
  cache, err := lru.New(10)
  if err != nil {
    assert.Error(t, err)
  }

  mr, err := miniredis.Run()
  if err != nil {
    panic(err)
  }

  redisClient := redis.NewClient(&redis.Options{
    Addr: mr.Addr(),
  })

  storeInstance, err := store.NewStore(cache, redisClient)
  assert.Nil(t, err, "return no error")
  assert.NotNil(t, storeInstance, "returns a valid store instance")
}

func TestAddInvalidURL(t *testing.T) {
  cache, err := lru.New(1)
  if err != nil {
    t.Error(err)
  }

  mr, err := miniredis.Run()
  if err != nil {
    t.Error(err)
  }

  redisClient := redis.NewClient(&redis.Options{
    Addr: mr.Addr(),
  })

  storeInstance, err := store.NewStore(cache, redisClient)
  if err != nil {
    t.Error(err)
  }

  shortKey, err := storeInstance.Add(addInvalidURL)
  assert.NotNil(t, err, "return error")
  assert.Equal(t, store.ErrorInvalidInputURL{} , err, "error type of store.ErrorInvalidInputURL")
  assert.Equal(t, "", shortKey, "short key should be empty")

}

func TestAddValidURL(t *testing.T) {
  cache, err := lru.New(1)
  if err != nil {
    t.Error(err)
  }

  mr, err := miniredis.Run()
  if err != nil {
    t.Error(err)
  }

  redisClient := redis.NewClient(&redis.Options{
    Addr: mr.Addr(),
  })

  storeInstance, err := store.NewStore(cache, redisClient)
  if err != nil {
    t.Error(err)
  }

  shortKey, err := storeInstance.Add(addValidURL)
  assert.Nil(t, err, "return no error")
  assert.NotEqual(t, "aa", shortKey, "short key should not be empty")
}

func TestGetExistingURL(t *testing.T) {}

func TestGetUnexistingURL(t *testing.T) {}

