package cache

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/VirrageS/chirp/backend/config"
)

type redisCache struct {
	client *redis.Client
	config config.CacheConfigProvider
}

// NewRedisCache constructs Redis cache which implements all Accessor functions.
func NewRedisCache(config config.RedisConfigProvider) Accessor {
	address := fmt.Sprintf("%s:%s", config.GetHost(), config.GetPort())
	password := config.GetPassword()
	db := config.GetDB()

	client := redis.NewClient(&redis.Options{
		Addr:       address,
		Password:   password,
		DB:         db,
		MaxRetries: 3,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.WithError(err).Error("Error connecting to cache instance.")
		return nil
	}

	return &redisCache{
		client: client,
		config: config,
	}
}

// Set adds all entries to cache. Value for key is either updated or created.
// Operations are pipelined so there is only one message sent to Redis instance.
func (cache *redisCache) Set(entries ...Entry) error {
	pipe := cache.client.Pipeline()
	for _, entry := range entries {
		data, err := cache.marshalValue(entry.Value)
		if err != nil {
			log.WithField("value", entry.Value).WithError(err).Error("Set: failed to marshal value")
			return err
		}

		pipe.Set(cache.convertKeyToHash(entry.Key), data, 0)
	}

	if _, err := pipe.Exec(); err != nil {
		log.WithField("entries", entries).WithError(err).Error("Set: failed to set values in cache")
		return err
	}

	return nil
}

// GetSingle retrieves single `value` stored at `key`.
// If key does not exist and no error occured returned bool is set to false.
func (cache *redisCache) GetSingle(key Key, value Value) (bool, error) {
	result, err := cache.client.Get(cache.convertKeyToHash(key)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.WithField("key", key).WithError(err).Error("GetSingle: failed to get key from cache")
		return false, err
	}

	err = cache.unmarshalValue(result, value)
	if err != nil {
		log.WithField("value", value).WithError(err).Error("GetSingle: failed to unmarshal value")
		return false, err
	}

	return true, nil
}

// Get retrieves multiple values.
// Operations are pipelined so there is only one message sent to Redis instance.
// Returned array of bools is set appropiretly depending on wheter key existed or not.
func (cache *redisCache) Get(entries ...Entry) ([]bool, error) {
	results := make([]*redis.StringCmd, len(entries))
	existing := make([]bool, len(entries))
	for i := range existing {
		existing[i] = true
	}

	pipe := cache.client.Pipeline()
	for i, entry := range entries {
		results[i] = pipe.Get(cache.convertKeyToHash(entry.Key))
	}

	if _, err := pipe.Exec(); err != nil && err != redis.Nil {
		log.WithField("entries", entries).WithError(err).Error("Get: failed to set values in cache")
		return nil, err
	}

	for i, entry := range entries {
		result, err := results[i].Result()
		if err == redis.Nil {
			existing[i] = false
			continue
		}

		err = cache.unmarshalValue(result, entry.Value)
		if err != nil {
			log.WithField("value", entry.Value).WithError(err).Error("Get: failed to unmarshal value")
			return nil, err
		}
	}

	return existing, nil
}

// Delete deletes multiple keys from cache.
// Operations are pipelined so there is only one message sent to Redis instance.
func (cache *redisCache) Delete(keys ...Key) error {
	pipe := cache.client.Pipeline()
	for _, key := range keys {
		pipe.Del(cache.convertKeyToHash(key))
	}

	if _, err := pipe.Exec(); err != nil {
		log.WithField("keys", keys).WithError(err).Error("Delete: failed to remove keys from cache")
		return err
	}

	return nil
}

// Incr increases values by 1 for each key. If key does not exist value is set to 1.
// Operations are pipelined so there is only one message sent to Redis instance.
func (cache *redisCache) Incr(keys ...Key) error {
	pipe := cache.client.Pipeline()
	for _, key := range keys {
		pipe.Incr(cache.convertKeyToHash(key))
	}

	if _, err := pipe.Exec(); err != nil {
		log.WithField("keys", keys).WithError(err).Error("Incr: failed to increment keys in cache")
		return err
	}

	return nil
}

// Decr decreases values by 1 for each key in `keys`. If key does not exist value is set to -1.
// Operations are pipelined so there is only one message sent to Redis instance.
func (cache *redisCache) Decr(keys ...Key) error {
	pipe := cache.client.Pipeline()
	for _, key := range keys {
		pipe.Decr(cache.convertKeyToHash(key))
	}

	if _, err := pipe.Exec(); err != nil {
		log.WithField("keys", keys).WithError(err).Error("Decr: failed to decrement keys in cache")
		return err
	}

	return nil
}

// SAdd adds array of `values` to set stored at `key`.
// Since it is set multiple values will squashed into single.
func (cache *redisCache) SAdd(key Key, values Values) error {
	members, err := cache.marshalValues(values)
	if err != nil {
		log.WithField("values", values).WithError(err).Error("SAdd: failed to marshal values")
		return err
	}

	return cache.client.SAdd(cache.convertKeyToHash(key), members...).Err()
}

// SMembers returns array of `values` from set stored at `key`.
// Returned bool is set to false when no error occured and `key` does not exist.
func (cache *redisCache) SMembers(key Key, values Values) (bool, error) {
	hash := cache.convertKeyToHash(key)
	exists, err := cache.client.Exists(hash).Result()
	if err != nil {
		log.WithField("key", key).WithError(err).Error("SMembers: failed to check if key exists")
		return false, err
	}

	if !exists {
		return false, nil
	}

	results, err := cache.client.SMembers(hash).Result()
	if err != nil {
		return false, err
	}

	err = cache.unmarshalValues(results, values)
	if err != nil {
		log.WithField("values", values).WithError(err).Error("SMembers: failed to unmarshal value")
		return false, err
	}

	return true, nil
}

// SRemove removes all `values` from set stored at `key`.
func (cache *redisCache) SRemove(key Key, values Values) error {
	members, err := cache.marshalValues(values)
	if err != nil {
		log.WithField("values", values).WithError(err).Error("SRemove: failed to marshal values")
		return err
	}

	return cache.client.SRem(cache.convertKeyToHash(key), members...).Err()
}

// Flush performs full clean on cache.
func (cache *redisCache) Flush() error {
	return cache.client.FlushAll().Err()
}

// unmarshalValue unmarshals single value from string.
func (cache *redisCache) unmarshalValue(result string, value Value) error {
	var (
		val int64
		err error
	)

	switch value := value.(type) {
	case *int64:
		*value, err = strconv.ParseInt(result, 10, 64)
	case *int32:
		val, err = strconv.ParseInt(result, 10, 32)
		*value = int32(val)
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

	return err
}

// marshalValue marshals single value.
func (cache *redisCache) marshalValue(value Value) (interface{}, error) {
	var (
		data interface{}
		err  error
	)

	switch value := value.(type) {
	case int64, int32, int16, int8, int:
		data = value
	default:
		data, err = msgpack.Marshal(value)
	}

	return data, err
}

// marshalValues marshals array of values
func (cache *redisCache) marshalValues(values Values) ([]interface{}, error) {
	var members []interface{}

	switch v := values.(type) {
	case []int64:
		members = make([]interface{}, 0, len(v))
		for _, value := range v {
			data, err := cache.marshalValue(value)
			if err != nil {
				return nil, err
			}

			members = append(members, data)
		}
	default:
		panic("marshalValues: invalid `values` type")
	}

	return members, nil
}

// unmarshalValues unmarshals array of values from array of strings.
func (cache *redisCache) unmarshalValues(results []string, values Values) error {
	for _, result := range results {
		var err error

		switch v := values.(type) {
		case *[]int64:
			var item int64
			err = cache.unmarshalValue(result, &item)
			*v = append(*v, item)
		default:
			panic("unmarshalValues: unsupported values")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// convertKeyToHash joins multiple fields in `key` to single hash
func (cache *redisCache) convertKeyToHash(key Key) string {
	stringFields := make([]string, 0, len(key))

	for _, field := range key {
		switch field.(type) {
		case string:
			stringFields = append(stringFields, field.(string))
		case int:
			stringFields = append(stringFields, strconv.FormatInt(int64(field.(int)), 10))
		case int32:
			stringFields = append(stringFields, strconv.FormatInt(int64(field.(int32)), 10))
		case int64:
			stringFields = append(stringFields, strconv.FormatInt(int64(field.(int64)), 10))
		default:
			panic("convertKeyToHash: unsupported field in key")
		}
	}

	return strings.Join(stringFields, ":")
}
