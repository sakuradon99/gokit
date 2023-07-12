package opt

import (
	"testing"
)

func Test_Optional(t *testing.T) {
	a := OfNullable[int](nil)
	if a.Exists() {
		t.Error("a.Exists() should be false")
	}
	if a.Get() != 0 {
		t.Error("a.Value() should be 0")
	}

	b := OfNullable[int](new(int))
	if !b.Exists() {
		t.Error("b.Exists() should be true")
	}
}
