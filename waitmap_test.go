package waitmap

import (
	"testing"
)

func TestFlatGet(t *testing.T) {
	m := New()
	m.Set(1, 2)
	m.Set(2, 3)
	if m.Get(1) != 2 {
		t.Errorf("Wrong result from Get(1)")
	}
	if m.Get(2) != 3 {
		t.Errorf("Wrong result from Get(2)")
	}
}

func TestFlatCheck(t *testing.T) {
	m := New()
	if m.Check(1) {
		t.Errorf("Check(1) returned true; false expected")
	}

	m.Set(1, 2)
	if !m.Check(1) {
		t.Errorf("Check(1) returned false; true expected")
	}
	if m.Check(2) {
		t.Errorf("Check(2) returned true; false expected")
	}
}

func TestSimple(t *testing.T) {
	m := New()
	var v interface{}
	c := make(chan interface{})
	go func() {
		v = m.Get(1)
		c <- nil
	}()
	go func() {
		m.Set(1, "hi")
	}()
	<-c
	if v != "hi" {
		t.Errorf("m.Get(1) returned incorrect valeu")
	}
}

func BenchmarkFlatFailedCheck(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Check(0)
	}
}

func BenchmarkFlatWaitedGet(b *testing.B) {
	m := New()
	go func() { m.Get(0) }()
	m.Set(0, "ho")
	for i := 0; i < b.N; i++ {
		m.Get(0)
	}
}

func BenchmarkFlatSimpleGet(b *testing.B) {
	m := New()
	m.Set(0, "ho")
	for i := 0; i < b.N; i++ {
		m.Get(0)
	}
}

func BenchmarkFlatSimpleSet(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Set(0, 0)
	}
}

func BenchmarkFlatIncrementalSet(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Set(i, 0)
	}
}
