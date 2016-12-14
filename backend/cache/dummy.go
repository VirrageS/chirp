package cache

type DummyCache struct{}

func NewDummyCache() CacheProvider {
	return &DummyCache{}
}

func (cache *DummyCache) Set(key string, value interface{}) {
}

func (cache *DummyCache) SetWithFields(fields Fields, value interface{}) {
}

func (cache *DummyCache) Get(key string) (interface{}, bool) {
	return nil, false
}

func (cache *DummyCache) GetWithFields(fields Fields) (interface{}, bool) {
	return nil, false
}

func (cache *DummyCache) Delete(key string) {
}

func (cache *DummyCache) DeleteWithFields(fields Fields) {
}
