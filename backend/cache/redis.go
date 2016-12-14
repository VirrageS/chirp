package cache

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"

	"github.com/VirrageS/chirp/backend/config"
)

// Default port on which Redis is listening
const DefaultRedisPort = "6379"

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
func (cache *RedisCache) Set(key string, value interface{}) {
	err := cache.client.Set(key, value, cache.config.GetCacheExpirationTime()).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"value": value,
		}).WithError(err).Error("set: failed to set key and value in cache")
	}
}

// Set `value` for specified key where key is created by hashing `fields`
func (cache *RedisCache) SetWithFields(fields Fields, value interface{}) {
	cache.Set(convertFieldsToKey(fields), value)
}

// Get value for specified `key`
func (cache *RedisCache) Get(key string) (interface{}, bool) {
	result, err := cache.client.Get(key).Result()
	if err != nil {
		log.WithField("key", key).WithError(err).Error("Get: failed to get key from cache")
		return nil, false
	}

	return result, true
}

// Get value for specified key where key is created by hashing `fields`
func (cache *RedisCache) GetWithFields(fields Fields) (interface{}, bool) {
	return cache.Get(convertFieldsToKey(fields))
}

// Delete value for specific `key`
func (cache *RedisCache) Delete(key string) {
	err := cache.client.Del(key).Err()
	if err != nil {
		log.WithField("key", key).WithError(err).Error("delete: failed to remove key from cache")
	}
}

// Delete value for specific key where key is created by hashing `fields`
func (cache *RedisCache) DeleteWithFields(fields Fields) {
	cache.Delete(convertFieldsToKey(fields))
}
