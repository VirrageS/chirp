package cache

type CacheProvider interface {
	Set(key string, value interface{})
	SetWithFields(fields Fields, value interface{})

	Get(key string) (interface{}, bool)
	GetWithFields(fields Fields) (interface{}, bool)

	Delete(key string)
	DeleteWithFields(fields Fields)
}
