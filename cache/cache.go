package cache

import (
	"errors"
	"sync"
)

var cache map[string]map[string]any = make(map[string]map[string]any)

var (
	ErrAlreadyCached = errors.New("Already Cached.")
)

var writeMutex sync.Mutex
var readMutex sync.Mutex

// Inserts a spell into the cache.
// Needs a module, a name and the spell to insert.
//
// Returns an error when it's already cached. You can ignore by checking
// ErrAlreadyCached.
func Insert(module string, name string, spell any) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()
	if Retrieve(module, name) != nil {
		return ErrAlreadyCached
	}
	if cache[module] == nil {
		cache[module] = make(map[string]any)
	}
	cache[module][name] = spell
	return nil
}

// Retrieves a spell from the cache, if it is present. Returns nil when the
// module is not in the cache.
func Retrieve(module string, name string) any {
	readMutex.Lock()
	defer readMutex.Unlock()
	if cache[module] == nil {
		return nil
	}
	if spell, present := cache[module][name]; present {
		return spell
	}
	return nil
}
