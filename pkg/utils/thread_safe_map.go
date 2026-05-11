package utils

import "sync"

type ThreadSafeMap[K comparable, V any] struct {
	mutex sync.RWMutex

	innerMap map[K]V
}

func NewThreadSafeMap[K comparable, V any]() *ThreadSafeMap[K, V] {
	return &ThreadSafeMap[K, V]{
		innerMap: make(map[K]V),
	}
}

func (m *ThreadSafeMap[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, ok := m.innerMap[key]
	return value, ok
}

func (m *ThreadSafeMap[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.innerMap[key] = value
}

func (m *ThreadSafeMap[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.innerMap, key)
}

func (m *ThreadSafeMap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	keys := make([]K, 0, len(m.innerMap))
	for key := range m.innerMap {
		keys = append(keys, key)
	}

	return keys
}

func (m *ThreadSafeMap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values := make([]V, 0, len(m.innerMap))
	for _, value := range m.innerMap {
		values = append(values, value)
	}

	return values
}

func (m *ThreadSafeMap[K, V]) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.innerMap)
}

func (m *ThreadSafeMap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.innerMap = make(map[K]V)
}

func (m *ThreadSafeMap[K, V]) IsEmpty() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.innerMap) == 0
}

func (m *ThreadSafeMap[K, V]) Has(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.innerMap[key]
	return ok
}
