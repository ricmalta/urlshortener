package store

import (
	"fmt"
  "github.com/sirupsen/logrus"
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
  logger *logrus.Logger
}

func NewStore(cache *lru.Cache, redisClient *redis.Client, logger *logrus.Logger) (*Store, error) {
	storeInstance := &Store{
		cache:       cache,
		redisClient: redisClient,
    logger: logger,
	}
	return storeInstance, nil
}

func (store *Store) Add(URL string) (shortKey string, err error) {
  if match, _ := regexp.MatchString(`^http(s)?://[a-z0-9A-Z\.]+/?`, URL); !match {
    return "", ErrorInvalidInputURL{}
  }
	incrCmd := store.redisClient.Incr(counterKeyName)
  if incrCmd.Err() != nil {
    store.logger.Error(err)
    return "", incrCmd.Err()
  }
  key := strconv.FormatInt(incrCmd.Val(), 36)
  statusCmd := store.redisClient.Set(key, URL, 0)
  if statusCmd.Err() != nil {
    store.logger.Error(err)
    return "", statusCmd.Err()
  }
  store.logger.Infof("new URL '%s' with short key '%s'", URL, key)
	return key, err
}

func (store *Store) Get(shortKey string) (url string, err error) {
  value, ok := store.cache.Get(shortKey)
  if ok {
    store.logger.Infof("return '%s' from memory cache", shortKey)
    return fmt.Sprintf("%v", value), nil
  }
  stringCmd := store.redisClient.Get(shortKey)
  if stringCmd.Val() != "" {
    store.logger.Infof("fetch '%s' from data store", shortKey)
    store.cache.Add(shortKey, stringCmd.Val())
    return store.Get(shortKey)
  }
  return "", ErrorNotStoredShortURL{}
}

func (store *Store) Shutdown() {

}
