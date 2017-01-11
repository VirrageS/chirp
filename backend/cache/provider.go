package cache

type CacheProvider interface {
	Set(key string, value interface{}) error
	SetWithFields(fields Fields, value interface{}) error

	SetWithoutExpiration(key string, value interface{}) error
	SetWithFieldsWithoutExpiration(fields Fields, value interface{}) error

	Increment(key string) error
	IncrementWithFields(fields Fields) error

	Decrement(key string) error
	DecrementWithFields(fields Fields) error

	Get(key string, value interface{}) (bool, error)
	GetWithFields(fields Fields, value interface{}) (bool, error)

	Delete(key string) error
	DeleteWithFields(fields Fields) error

	Flush() error
}
