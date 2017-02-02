package cache

type fakeCache struct{}

// NewFakeCache creates new instance of fake cache which imitates caching values.
// Implements all Accessor methods but each of them is doing nothing.
func NewFakeCache() Accessor {
	return &fakeCache{}
}

func (cache *fakeCache) Set(entries ...Entry) error {
	return nil
}

func (cache *fakeCache) GetSingle(key Key, value Value) (bool, error) {
	return false, nil
}

func (cache *fakeCache) Get(entries ...Entry) ([]bool, error) {
	results := make([]bool, len(entries))
	for i := range results {
		results[i] = false
	}

	return results, nil
}

func (cache *fakeCache) Delete(keys ...Key) error {
	return nil
}

func (cache *fakeCache) Incr(keys ...Key) error {
	return nil
}

func (cache *fakeCache) Decr(keys ...Key) error {
	return nil
}

func (cache *fakeCache) SAdd(key Key, values Values) error {
	return nil
}

func (cache *fakeCache) SMembers(key Key, values Values) (bool, error) {
	return false, nil
}

func (cache *fakeCache) SRemove(key Key, values Values) error {
	return nil
}

func (cache *fakeCache) Flush() error {
	return nil
}
