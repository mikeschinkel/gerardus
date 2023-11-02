package parser

import (
	"sync"
)

// See: https://hjr265.me/blog/synchronization-constructs-in-go-standard-library/

// SafeMap is a generic map that uses sync.RWMutex for concurrency-safe access
type SafeMap[K comparable, V any] struct {
	sync.RWMutex
	data map[K]V
}

// NewSafeMap instantiates a new instance of for SafeMap{KeyType,ValueType], e.g.
//
//	myMap := NewSafeMap[int,string]()
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *SafeMap[K, V]) Has(k K) (ok bool) {
	m.RWMutex.RLock()
	_, ok = m.data[k]
	m.RWMutex.RUnlock()
	return ok
}

func (m *SafeMap[K, V]) Load(k K) (v V, ok bool) {
	m.RWMutex.RLock()
	v, ok = m.data[k]
	m.RWMutex.RUnlock()
	return v, ok
}

func (m *SafeMap[K, V]) Save(k K, v V) {
	m.RWMutex.Lock()
	m.data[k] = v
	m.RWMutex.Unlock()
}
