package cache

// Key represents abstract identificator for cache at which values/sets can be stored.
type Key []interface{}

// Value represents single value which can be added or retrieved from cache.
type Value interface{}

// Values represents array of values which can be added or retrieved from set in cache.
type Values interface{}

// Entry represents of single <key, value> pair which can be added or retrieved from cache.
type Entry struct {
	Key   Key
	Value Value
}

// Accessor is interface which defines all cache functions used in system.
// All implementations of different types of cache should implement these methods.
type Accessor interface {
	// Set adds or updates values for each key in `entries`.
	Set(entries ...Entry) error

	// GetSingle retrieves value for key in `entry`.
	// If key exists `true` is returned, `false` otherwise.
	GetSingle(key Key, value Value) (bool, error)

	// Get retrieves values for each key in `entries`.
	// Returns array of booleans which tells if key existed or not.
	//
	// Returned array is always of the same length as `entries`.
	// Returned error is first error which occured during executing function.
	Get(entries ...Entry) ([]bool, error)

	// Delete removes `key` and associated value from cache.
	Delete(keys ...Key) error

	// Incr increments value for each key in `keys`.
	Incr(keys ...Key) error

	// Decr decrements value for each key in `keys`.
	Decr(keys ...Key) error

	// SAdd adds array of new `values` for set stored at `key`.
	SAdd(key Key, values Values) error

	// SMembers gets all values from set stored at `key`.
	// `values` should be an array of values to which members will be unmarshalled.
	SMembers(key Key, values Values) (bool, error)

	// SRemove removes array of `values` from set stored at `key`.
	SRemove(key Key, values Values) error

	// Flush performs full clean on cache.
	// In the result all keys and values should be removed from cache.
	Flush() error
}
