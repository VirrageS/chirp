package cache

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/VirrageS/chirp/backend/config"
)

type RedisCache struct {
	client *redis.Client
	config config.CacheConfigProvider
}

// Creates new CacheProvider from Redis client
func NewRedisCache(config config.RedisConfigProvider) CacheProvider {
	address := fmt.Sprintf("%s:%s", config.GetHost(), config.GetPort())
	password := config.GetPassword()
	db := config.GetDB()

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.WithError(err).Fatal("Error connecting to cache instance.")
	}

	return &RedisCache{
		client: client,
		config: config,
	}
}

// Set `value` for specified `key`
func (cache *RedisCache) Set(key string, value interface{}) error {
	return cache.set(key, value, cache.config.GetCacheExpirationTime())
}

// Set `value` for specified key where key is created by hashing `fields`
func (cache *RedisCache) SetWithFields(fields Fields, value interface{}) error {
	return cache.Set(convertFieldsToKey(fields), value)
}

// Set `value` for specified `key` without expiration time
func (cache *RedisCache) SetWithoutExpiration(key string, value interface{}) error {
	return cache.set(key, value, 0)
}

// Set `value` for specified key where key is created by hashing `fields`, without expiration time
func (cache *RedisCache) SetWithFieldsWithoutExpiration(fields Fields, value interface{}) error {
	return cache.SetWithoutExpiration(convertFieldsToKey(fields), value)
}

// Atomic increment of a value for specified `key`
func (cache *RedisCache) Increment(key string) error {
	err := cache.client.Incr(key).Err()
	if err != nil {
		log.WithField("key", key).WithError(err).Error("increment: failed to increment key in cache")
		return err
	}

	return nil
}

// Increment value for specified key where key is created by hashing `fields`
func (cache *RedisCache) IncrementWithFields(fields Fields) error {
	return cache.Increment(convertFieldsToKey(fields))
}

// Atomic decrement of a value for specified `key`
func (cache *RedisCache) Decrement(key string) error {
	err := cache.client.Decr(key).Err()
	if err != nil {
		log.WithField("key", key).WithError(err).Error("decrement: failed to decrement key in cache")
		return err
	}

	return nil
}

// Decrement value for specified key where key is created by hashing `fields`
func (cache *RedisCache) DecrementWithFields(fields Fields) error {
	return cache.Decrement(convertFieldsToKey(fields))
}

// Get value for specified `key`
func (cache *RedisCache) Get(key string, value interface{}) (bool, error) {
	var err error
	var val int64

	result, err := cache.client.Get(key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.WithField("key", key).WithError(err).Error("get: failed to get key from cache")
		return false, err
	}

	switch value := value.(type) {
	case *int64:
		*value, err = strconv.ParseInt(result, 10, 64)
	case *int32:
		val, err = strconv.ParseInt(result, 10, 32)
		*value = int32(val) // We can just cast it here, because on error val will be = 0
	case *int16:
		val, err = strconv.ParseInt(result, 10, 16)
		*value = int16(val)
	case *int8:
		val, err = strconv.ParseInt(result, 10, 8)
		*value = int8(val)
	case *int:
		val, err = strconv.ParseInt(result, 10, 0)
		*value = int(val)
	default:
		err = msgpack.Unmarshal([]byte(result), value)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"value": result,
		}).WithError(err).Error("get: failed to unmarshal value")
		return false, err
	}

	return true, nil
}

// Get value for specified key where key is created by hashing `fields`
func (cache *RedisCache) GetWithFields(fields Fields, value interface{}) (bool, error) {
	return cache.Get(convertFieldsToKey(fields), value)
}

// Delete value for specific `key`
func (cache *RedisCache) Delete(key string) error {
	err := cache.client.Del(key).Err()
	if err != nil {
		log.WithField("key", key).WithError(err).Error("delete: failed to remove key from cache")
		return err
	}

	return nil
}

// Delete value for specific key where key is created by hashing `fields`
func (cache *RedisCache) DeleteWithFields(fields Fields) error {
	return cache.Delete(convertFieldsToKey(fields))
}

// Flush all cache
func (cache *RedisCache) Flush() error {
	return cache.client.FlushAll().Err()
}

func (cache *RedisCache) set(key string, value interface{}, expirationTime time.Duration) error {
	var data interface{}

	switch value := value.(type) {
	case int64, int32, int16, int8, int:
		data = value
	default:
		var err error
		data, err = msgpack.Marshal(value)
		if err != nil {
			log.WithFields(log.Fields{
				"key":   key,
				"value": value,
			}).WithError(err).Error("set: failed to marshal value")
			return err
		}
	}

	err := cache.client.Set(key, data, expirationTime).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"value": value,
		}).WithError(err).Error("set: failed to set key and value in cache")
		return err
	}

	return nil
}
