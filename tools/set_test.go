package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Set(t *testing.T) {
	set := NewUnstableSet[int]()
	set.Put(1)
	set.Put(2)

	assert.Equal(t, true, set.Contains(1))
	assert.Equal(t, false, set.Contains(3))
	assert.Equal(t, 2, len(set.ToSlice()))
	assert.Equal(t, 2, set.Len())
}
