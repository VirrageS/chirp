package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/VirrageS/chirp/backend/model"
)

type mockCacheConfigProvider struct{}

func (m *mockCacheConfigProvider) GetCacheExpirationTime() time.Duration {
	return time.Second
}

var objectTests = []struct {
	in  interface{}
	out interface{}
}{
	{
		&model.User{
			ID:       1,
			Username: "username",
		},
		&model.User{},
	},
	{
		&[]*model.User{
			{ID: 1, Username: "username1"},
			{ID: 2, Username: "username2"},
		},
		&[]*model.User{},
	},
	{
		&model.PublicUser{
			ID:        2,
			Username:  "username",
			Name:      "name",
			AvatarUrl: "avatars@web.com",
			Following: true,
		},
		&model.PublicUser{},
	},
}

var fieldsTests = []Fields{
	{"key", 12, "hello", -1},
	{"key", 12},
	{-1, "key", 12},
	{0},
	{"key"},
}

var redisCache CacheProvider = NewRedisCache("6379", &mockCacheConfigProvider{})

func TestRedisCacheSet(t *testing.T) {
	redisCache.Flush()

	for _, test := range objectTests {
		err := redisCache.Set("key", test.in)
		assert.Nil(t, err)
	}
}

func TestRedisCacheSetWithFields(t *testing.T) {
	redisCache.Flush()

	for _, test := range objectTests {
		for _, fields := range fieldsTests {
			err := redisCache.SetWithFields(fields, test.in)
			assert.Nil(t, err)
		}
	}
}

func TestRedisCacheGetNotExists(t *testing.T) {
	redisCache.Flush()

	for _, test := range objectTests {
		exists, err := redisCache.Get("key", &test.in)
		assert.False(t, exists)
		assert.Nil(t, err)
	}
}

func TestRedisCacheGetExists(t *testing.T) {
	redisCache.Flush()

	for _, test := range objectTests {
		err := redisCache.Set("key", test.in)
		assert.Nil(t, err)

		exists, err := redisCache.Get("key", test.out)
		assert.Nil(t, err)
		assert.True(t, exists)

		assert.Equal(t, test.in, test.out)
	}
}

func TestRedisCacheGetWithFieldsNotExists(t *testing.T) {
	redisCache.Flush()

	for _, test := range objectTests {
		for _, fields := range fieldsTests {
			exists, err := redisCache.GetWithFields(fields, test.out)
			assert.Nil(t, err)
			assert.False(t, exists)
		}
	}
}

func TestRedisCacheGetWithFieldsExists(t *testing.T) {
	redisCache.Flush()

	for _, test := range objectTests {
		for _, fields := range fieldsTests {
			exists, err := redisCache.GetWithFields(fields, test.out)
			assert.Nil(t, err)
			assert.False(t, exists)
		}
	}
}

func TestRedisCacheDelete(t *testing.T) {
	redisCache.Flush()

	err := redisCache.Delete("key")
	assert.Nil(t, err)
}

func TestDeleteExists(t *testing.T) {
	redisCache.Flush()

	for _, test := range objectTests {
		err := redisCache.Set("key", test.in)
		assert.Nil(t, err)

		err = redisCache.Delete("key")
		assert.Nil(t, err)

		exists, err := redisCache.Get("key", test.out)
		assert.False(t, exists)
		assert.Nil(t, err)
	}
}

func TestRedisCacheDeleteWithFields(t *testing.T) {
	redisCache.Flush()

	for _, fields := range fieldsTests {
		err := redisCache.DeleteWithFields(fields)
		assert.Nil(t, err)
	}
}

func TestRedisCacheDeleteWithFieldsExists(t *testing.T) {
	redisCache.Flush()

	for _, test := range objectTests {
		for _, fields := range fieldsTests {
			err := redisCache.SetWithFields(fields, test.in)
			assert.Nil(t, err)

			err = redisCache.DeleteWithFields(fields)
			assert.Nil(t, err)

			exists, err := redisCache.GetWithFields(fields, test.out)
			assert.False(t, exists)
			assert.Nil(t, err)
		}
	}
}
