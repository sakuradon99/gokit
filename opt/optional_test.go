package opt

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testJson struct {
	O Optional[int] `json:"o"`
}

func Test_Optional(t *testing.T) {
	o := Empty[int]()
	assert.Equal(t, false, o.Exists())
	assert.Equal(t, 0, o.Get())
	assert.Equal(t, true, o.GetPointer() == nil)

	o = Of[int](1)
	assert.Equal(t, true, o.Exists())
	assert.Equal(t, 1, o.Get())
	v, ok := o.GetAndExists()
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, *o.GetPointer())
	assert.Equal(t, 1, o.GetOrElse(2))

	o = OfNullable[int](nil)
	assert.Equal(t, false, o.Exists())
	assert.Equal(t, 0, o.Get())
	assert.Equal(t, 1, o.GetOrElse(1))
	assert.Panics(t, func() {
		o.MustGet()
	})

	o = OfNullable[int](&[]int{1}[0])
	assert.Equal(t, true, o.Exists())
	assert.Equal(t, 1, o.Get())
	assert.Equal(t, 1, o.MustGet())

	j := testJson{
		O: Of[int](1),
	}
	b, err := json.Marshal(&j)
	assert.Nil(t, err)
	assert.Equal(t, `{"o":1}`, string(b))

	j = testJson{
		O: Empty[int](),
	}
	b, err = json.Marshal(&j)
	assert.Nil(t, err)
	assert.Equal(t, `{"o":null}`, string(b))

	j = testJson{}
	err = json.Unmarshal([]byte(`{"o":1}`), &j)
	assert.Nil(t, err)
	assert.Equal(t, true, j.O.Exists())

	j = testJson{}
	err = json.Unmarshal([]byte(`{"o":null}`), &j)
	assert.Nil(t, err)
	assert.Equal(t, false, j.O.Exists())

	j = testJson{}
	err = json.Unmarshal([]byte(`{"o":true}`), &j)
	assert.NotNil(t, err)
	assert.Equal(t, false, j.O.Exists())

	d := Empty[int]()
	err = d.Scan(1)
	assert.Nil(t, err)
	assert.Equal(t, true, d.Exists())
	assert.Equal(t, 1, d.Get())

	d = Empty[int]()
	err = d.Scan(nil)
	assert.Nil(t, err)
	assert.Equal(t, false, d.Exists())

	d = Empty[int]()
	err = d.Scan("1")
	assert.NotNil(t, err)
	assert.Equal(t, false, d.Exists())

	d = Empty[int]()
	dbVal, err := d.Value()
	assert.Nil(t, err)
	assert.Equal(t, nil, dbVal)

	d = Of[int](1)
	dbVal, err = d.Value()
	assert.Nil(t, err)
	assert.Equal(t, 1, dbVal)
}
