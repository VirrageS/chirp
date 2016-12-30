package cache

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
	"gopkg.in/vmihailenco/msgpack.v2"

	"fmt"
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
	bytes, err := msgpack.Marshal(value)
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"value": value,
		}).WithError(err).Error("set: failed to marshal value")
		return err
	}

	err = cache.client.Set(key, bytes, cache.config.GetCacheExpirationTime()).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"value": value,
		}).WithError(err).Error("set: failed to set key and value in cache")
		return err
	}

	return nil
}

// Set `value` for specified key where key is created by hashing `fields`
func (cache *RedisCache) SetWithFields(fields Fields, value interface{}) error {
	return cache.Set(convertFieldsToKey(fields), value)
}

// Set integer `value` for a specified `key`
func (cache *RedisCache) SetInt(key string, value int64) error {
	err := cache.client.Set(key, value, cache.config.GetCacheExpirationTime()).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"value": value,
		}).WithError(err).Error("setInt: failed to set key and value in cache")
		return err
	}

	return nil
}

// Set integer `value` for specified key where key is created by hashing `fields`
func (cache *RedisCache) SetIntWithFields(fields Fields, value int64) error {
	return cache.SetInt(convertFieldsToKey(fields), value)
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
	bytes, err := cache.client.Get(key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.WithField("key", key).WithError(err).Error("Get: failed to get key from cache")
		return false, err
	}

	err = msgpack.Unmarshal([]byte(bytes), value)
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"value": value,
		}).WithError(err).Error("set: failed to marshal value")
		return false, err
	}

	return true, nil
}

// Get value for specified key where key is created by hashing `fields`
func (cache *RedisCache) GetWithFields(fields Fields, value interface{}) (bool, error) {
	return cache.Get(convertFieldsToKey(fields), value)
}

// Get integer value for specified `key`
func (cache *RedisCache) GetInt(key string, value *int64) (bool, error) {
	val, err := cache.client.Get(key).Int64()
	if err == redis.Nil {
		*value = int64(0)
		return false, nil
	} else if err != nil {
		log.WithField("key", key).WithError(err).Error("GetInt: failed to get key from cache")
		*value = int64(0)
		return false, err
	}

	*value = val
	return true, nil
}

// Get integer value for specified where key is created by hashing `fields`
func (cache *RedisCache) GetIntWithFields(fields Fields, value *int64) (bool, error) {
	return cache.GetInt(convertFieldsToKey(fields), value)
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
