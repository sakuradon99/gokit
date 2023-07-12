package db

import (
	"context"
	"gorm.io/gorm"
)

const (
	KeyDB = "db_tx"
)

type Manager interface {
	DB(ctx context.Context) *gorm.DB
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}

type ManagerImpl struct {
	db *gorm.DB `inject:""`
}

func (m *ManagerImpl) DB(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(KeyDB).(*gorm.DB)
	if !ok {
		db = m.db.WithContext(ctx)
	}
	return db
}

func (m *ManagerImpl) Transaction(ctx context.Context, f func(ctx context.Context) error) error {
	return m.DB(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, KeyDB, tx)
		return f(txCtx)
	})
}
