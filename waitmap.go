// Package waitmap implements a simple thread-safe map.
package waitmap

import (
	"sync"
)

// An entry in a WaitMap. Note that mutx and cond may legally be nil - this
// avoids heavy allocations at the cost of some code complexity.
type entry struct {
	mutx *sync.Mutex
	cond *sync.Cond
	data interface{}
	ok   bool
}

type WaitMap struct {
	lock *sync.Mutex
	ents map[interface{}]*entry
}

func New() *WaitMap {
	return &WaitMap{
		lock: new(sync.Mutex),
		ents: make(map[interface{}]*entry),
	}
}

// Retrieves the value mapped to by k. If no such value is yet available, waits
// until one is, and then returns that. Otherwise, it returns the value
// available at the time Get is called.
func (m *WaitMap) Get(k interface{}) interface{} {
	m.lock.Lock()
	e, ok := m.ents[k]
	if !ok {
		mutx := new(sync.Mutex)
		e = &entry{
			mutx: mutx,
			cond: sync.NewCond(mutx),
			data: nil,
			ok:   false,
		}
		m.ents[k] = e
	}
	m.lock.Unlock()

	// If e.ok is true, e.data exists and can never cease to exist. We need
	// this check to avoid using a nil e.mutx. We could also actually check
	// e.mutx to see if it's nil, but this accomplishes the same thing and
	// will be slightly faster on average (since we will often avoid
	// unnecessarily messing with the mutex).
	if e.ok {
		return e.data
	}
	e.mutx.Lock()
	defer e.mutx.Unlock()
	e.cond.Wait()
	return e.data
}

// Maps the given key and value, waking any waiting calls to Get.
func (m *WaitMap) Set(k interface{}, v interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	e, ok := m.ents[k]
	if !ok {
		e := &entry{nil, nil, v, true}
		m.ents[k] = e
		return
	}
	e.data = v
	if !e.ok {
		e.ok = true
		e.cond.Broadcast()
	}
}

// Returns true if k is a key in the map.
func (m *WaitMap) Check(k interface{}) bool {
	m.lock.Lock()
	e, ok := m.ents[k]
	m.lock.Unlock()
	if !ok {
		return false
	}
	return e.ok
}
