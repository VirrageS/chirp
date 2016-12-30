package cache

type CacheProvider interface {
	Set(key string, value interface{}) error
	SetWithFields(fields Fields, value interface{}) error

	SetInt(key string, value int64) error
	SetIntWithFields(fields Fields, value int64) error

	Increment(key string) error
	IncrementWithFields(fields Fields) error

	Decrement(key string) error
	DecrementWithFields(fields Fields) error

	Get(key string, value interface{}) (bool, error)
	GetWithFields(fields Fields, value interface{}) (bool, error)

	GetInt(key string, value *int64) (bool, error)
	GetIntWithFields(fields Fields, value *int64) (bool, error)

	Delete(key string) error
	DeleteWithFields(fields Fields) error

	Flush() error
}
