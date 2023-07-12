package db

import (
	"errors"
	"github.com/sakuradon99/gokit/internal/db"
	"gorm.io/gorm"
)

type Manager = db.Manager

func RecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
