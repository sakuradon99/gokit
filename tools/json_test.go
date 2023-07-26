package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testJson struct {
	A int
	B string
}

func Test_JSON(t *testing.T) {
	j := testJson{
		A: 1,
		B: "2",
	}

	jsonStr, err := MarshalJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, `{"A":1,"B":"2"}`, jsonStr)

	j2, err := UnmarshalJSON[testJson](jsonStr)
	assert.Nil(t, err)
	assert.Equal(t, j, j2)

	jsonStr, err = MarshalJSON(complex(2, 3))
	assert.NotNil(t, err)
	assert.Equal(t, "", jsonStr)
}
