// Package waitmap implements a simple thread-safe map.
package waitmap

import (
	"sync"
)

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
// until one is, and then returns that.
func (m *WaitMap) Get(k interface{}) interface{} {
	m.lock.Lock()
	e, ok := m.ents[k]
	if !ok {
		mutx := new(sync.Mutex)
		e = &entry{
			mutx: mutx,
			cond: sync.NewCond(mutx),
			data: nil,
			ok: false,
		}
		m.ents[k] = e
	}
	m.lock.Unlock()
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
	e.ok = true
	e.cond.Broadcast()
}

// Returns true if k is a key in the map.
func (m *WaitMap) Check(k interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	e, ok := m.ents[k]
	if !ok { return false }
	return e.ok
}

