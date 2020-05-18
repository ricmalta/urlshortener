package store

import (
	"fmt"
  "regexp"
  "strconv"

	"github.com/go-redis/redis"
	lru "github.com/hashicorp/golang-lru"
)

const (
	counterKeyName string = "counter"
)

type Store struct {
	cache       *lru.Cache
	redisClient *redis.Client
}

func NewStore(cache *lru.Cache, redisClient *redis.Client) (*Store, error) {
	storeInstance := &Store{
		cache:       cache,
		redisClient: redisClient,
	}

	return storeInstance, nil
}

func (store *Store) Add(URL string) (shortKey string, err error) {
  if match, _ := regexp.MatchString(`^http(s)?://[a-z0-9A-Z\.]+/?`, URL); !match {
    return "", ErrorInvalidInputURL{}
  }
	incrCmd := store.redisClient.Incr(counterKeyName)
  if incrCmd.Err() != nil {
    return "", incrCmd.Err()
  }
  key := strconv.FormatInt(incrCmd.Val(), 36)
  statusCmd := store.redisClient.Set(key, URL, 0)
  if statusCmd.Err() != nil {
    return "", statusCmd.Err()
  }
	return key, err
}

func (store *Store) Get(shortKey string) (url string, err error) {
  value, ok := store.cache.Get(shortKey)
  if ok {
    return fmt.Sprintf("%v", value), nil
  }
  stringCmd := store.redisClient.Get(shortKey)
  if stringCmd.Val() != "" {
    store.cache.Add(shortKey, stringCmd.Val())
    return store.Get(shortKey)
  }
  return "", ErrorNotStoredShortURL{}
}
