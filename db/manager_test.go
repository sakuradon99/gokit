package db

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func Test_RecordNotFound(t *testing.T) {
	err := gorm.ErrRecordNotFound
	assert.Equal(t, true, RecordNotFound(err))
	err = assert.AnError
	assert.Equal(t, false, RecordNotFound(err))
}
