package cache

type DummyCache struct{}

func NewDummyCache() CacheProvider {
	return &DummyCache{}
}

func (cache *DummyCache) Set(key string, value interface{}) error {
	return nil
}

func (cache *DummyCache) SetWithFields(fields Fields, value interface{}) error {
	return nil
}

func (cache *DummyCache) SetWithoutExpiration(key string, value interface{}) error {
	return nil
}

func (cache *DummyCache) SetWithFieldsWithoutExpiration(fields Fields, value interface{}) error {
	return nil
}

func (cache *DummyCache) Increment(key string) error {
	return nil
}

func (cache *DummyCache) IncrementWithFields(fields Fields) error {
	return nil
}

func (cache *DummyCache) Decrement(key string) error {
	return nil
}

func (cache *DummyCache) DecrementWithFields(fields Fields) error {
	return nil
}

func (cache *DummyCache) Get(key string, value interface{}) (bool, error) {
	return false, nil
}

func (cache *DummyCache) GetWithFields(fields Fields, value interface{}) (bool, error) {
	return false, nil
}

func (cache *DummyCache) Delete(key string) error {
	return nil
}

func (cache *DummyCache) DeleteWithFields(fields Fields) error {
	return nil
}

func (cache *DummyCache) Flush() error {
	return nil
}
