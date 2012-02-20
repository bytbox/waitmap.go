// Package waitmap implements a simple thread-safe map.
package waitmap

import (
	"sync"
)

// TODO use sync/atomic wherever possible

// An entry in a WaitMap.
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

func NewCap(c int) *WaitMap {
	return &WaitMap{
		lock: new(sync.Mutex),
		ents: make(map[interface{}]*entry, c),
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
	e.cond.Wait()
	e.mutx.Unlock()
	return e.data
}

// Maps the given key and value, waking any waiting calls to Get. Returns false
// (and changes nothing) if the key is already in the map.
func (m *WaitMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	e, ok := m.ents[k]
	if !ok {
		mutx := new(sync.Mutex)
		e := &entry{
			mutx: mutx,
			cond: sync.NewCond(mutx),
			data: v,
			ok:   true,
		}
		m.ents[k] = e
		m.lock.Unlock()
		return true
	}
	if e.ok {
		m.lock.Unlock()
		return false
	}
	e.mutx.Lock()
	m.lock.Unlock()
	e.data = v
	e.ok = true
	e.cond.Broadcast()
	e.mutx.Unlock()
	return true
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
