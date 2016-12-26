package cache

import (
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
func NewRedisCache(port string, config config.CacheConfigProvider) CacheProvider {
	// TODO: read user data, host and port from config file
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + port,
		Password: "", // no password set
		DB:       0,  // use default DB
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
