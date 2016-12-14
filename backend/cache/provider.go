package cache

type CacheProvider interface {
	Set(key string, value interface{}) error
	SetWithFields(fields Fields, value interface{}) error

	Get(key string, value interface{}) (bool, error)
	GetWithFields(fields Fields, value interface{}) (bool, error)

	Delete(key string) error
	DeleteWithFields(fields Fields) error

	Flush() error
}
