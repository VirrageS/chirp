package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type object struct {
	Str string
	Num int
}

var dummyCache CacheProvider = NewDummyCache()

func TestDummyCacheSet(t *testing.T) {
	err := dummyCache.Set("key", &object{
		Str: "wtf",
		Num: 12,
	})

	assert.Nil(t, err)
}

func TestDummyCacheSetWithFields(t *testing.T) {
	err := dummyCache.SetWithFields(Fields{"key", 12, "hello", -1}, &object{
		Str: "wtf",
		Num: 12,
	})

	assert.Nil(t, err)
}

func TestDummyCacheGetNoExists(t *testing.T) {
	var obj object
	exists, err := dummyCache.Get("key", &obj)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestDummyCacheGetExists(t *testing.T) {
	dummyCache.Set("key", &object{
		Str: "wtf",
		Num: 12,
	})

	var obj object
	exists, err := dummyCache.Get("key", &obj)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestDummyCacheGetWithFieldsNoExists(t *testing.T) {
	var obj object
	exists, err := dummyCache.GetWithFields(Fields{"key", "super", 1}, &obj)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestDummyCacheGetWithFieldsExists(t *testing.T) {
	err := dummyCache.SetWithFields(Fields{"key", "super", 1}, &object{
		Str: "wtf",
		Num: 12,
	})
	assert.Nil(t, err)

	var obj object
	exists, err := dummyCache.GetWithFields(Fields{"key", "super", 1}, &obj)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestDummyCacheDeleteNoExists(t *testing.T) {
	err := dummyCache.Delete("key")
	assert.Nil(t, err)
}

func TestDummyCacheDeleteExists(t *testing.T) {
	err := dummyCache.Set("key", &object{
		Str: "wtf",
		Num: 12,
	})
	assert.Nil(t, err)

	err = dummyCache.Delete("key")
	assert.Nil(t, err)
}

func TestDummyCacheDeleteWithFieldsNoExists(t *testing.T) {
	err := dummyCache.DeleteWithFields(Fields{"key", "super", 1})
	assert.Nil(t, err)
}

func TestDummyCacheDeleteWithFieldsExists(t *testing.T) {
	err := dummyCache.SetWithFields(Fields{"key", "super", 1}, &object{
		Str: "wtf",
		Num: 12,
	})
	assert.Nil(t, err)

	err = dummyCache.DeleteWithFields(Fields{"key", "super", 1})
	assert.Nil(t, err)
}
