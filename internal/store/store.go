package store

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/ricmalta/urlshortner/internal/config"

	"github.com/go-redis/redis"
	lru "github.com/hashicorp/golang-lru"
)

const (
	counterKeyName string = "counter"
)

type Store struct {
	cache       *lru.Cache
	redisClient *redis.Client
	wg          sync.WaitGroup
}

func NewStore(cfg config.Config) (*Store, error) {
	cache, err := lru.New(cfg.Cache.Size)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	// check redis connection
	if _, err := redisClient.Ping().Result(); err != nil {
		return nil, err
	}

	storeInstance := &Store{
		cache:       cache,
		redisClient: redisClient,
	}

	return storeInstance, nil
}

func (store *Store) Add(URL string) (shortKey string, err error) {
	store.wg.Add(1)

	var value string
	err = store.redisClient.Watch(func(tx *redis.Tx) error {
		counterIncr := tx.Incr(counterKeyName)
		if counterIncr.Err() != nil {
			return counterIncr.Err()
		}

		key := strconv.FormatInt(counterIncr.Val(), 36)
		_, err = tx.TxPipelined(func(pipe redis.Pipeliner) error {
			statusCmd := pipe.Set(key, URL, 0)
			if statusCmd.Err() != nil {
				return statusCmd.Err()
			}
			return nil
		})
		// set the return value
		value = key
		store.wg.Done()
		return nil
	}, counterKeyName)

	store.wg.Wait()
	return value, err
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
