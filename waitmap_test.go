package waitmap

import (
	"testing"
)

func TestFlatGet(t *testing.T) {
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

}
