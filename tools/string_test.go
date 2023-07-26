package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StringToInt(t *testing.T) {
	assert.Equal(t, 0, StringToInt("a"))
	assert.Equal(t, 1, StringToInt("1"))
}

func Test_StringIn(t *testing.T) {
	assert.Equal(t, true, StringIn("a", "a", "b"))
	assert.Equal(t, false, StringIn("c", "a", "b"))
}
