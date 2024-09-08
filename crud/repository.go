package crud

import (
	"context"
	"github.com/sakuradon99/gokit/db"
	"gorm.io/gorm"
)

type DBOption func(db *gorm.DB) *gorm.DB

type Repository[T any] interface {
	Insert(ctx context.Context, entity T, options ...DBOption) error
	InsertBatch(ctx context.Context, entities []T, options ...DBOption) error
	Delete(ctx context.Context, condition T, options ...DBOption) error
	DeleteByID(ctx context.Context, id any, options ...DBOption) error
	DeleteByIDs(ctx context.Context, ids []any, options ...DBOption) error
	Update(ctx context.Context, entity T, options ...DBOption) error
	UpdateBatch(ctx context.Context, entities []T, options ...DBOption) error
	SelectOne(ctx context.Context, condition T, options ...DBOption) (T, error)
	SelectList(ctx context.Context, condition T, options ...DBOption) ([]T, error)
	SelectByID(ctx context.Context, id any, options ...DBOption) (T, error)
	SelectByIDs(ctx context.Context, ids []any, options ...DBOption) ([]T, error)
	SelectCount(ctx context.Context, condition T, options ...DBOption) (int64, error)
}

type RepositoryImpl[T any] struct {
	dbm db.Manager `inject:""`
}

func (r *RepositoryImpl[T]) db(ctx context.Context, options ...DBOption) *gorm.DB {
	gdb := r.dbm.DB(ctx)
	for _, option := range options {
		gdb = option(gdb)
	}
	return gdb
}

func (r *RepositoryImpl[T]) Insert(ctx context.Context, entity T, options ...DBOption) error {
	return r.db(ctx, options...).Create(&entity).Error
}

func (r *RepositoryImpl[T]) InsertBatch(ctx context.Context, entities []T, options ...DBOption) error {
	return r.db(ctx, options...).Create(&entities).Error
}

func (r *RepositoryImpl[T]) Delete(ctx context.Context, condition T, options ...DBOption) error {
	var entity T
	return r.db(ctx, options...).Where(&condition).Delete(&entity).Error
}

func (r *RepositoryImpl[T]) DeleteByID(ctx context.Context, id any, options ...DBOption) error {
	var entity T
	return r.db(ctx, options...).Delete(&entity, id).Error
}

func (r *RepositoryImpl[T]) DeleteByIDs(ctx context.Context, ids []any, options ...DBOption) error {
	var entity T
	return r.db(ctx, options...).Delete(&entity, ids).Error
}

func (r *RepositoryImpl[T]) Update(ctx context.Context, entity T, options ...DBOption) error {
	return r.db(ctx, options...).Save(&entity).Error
}

func (r *RepositoryImpl[T]) UpdateBatch(ctx context.Context, entities []T, options ...DBOption) error {
	return r.db(ctx, options...).Save(&entities).Error
}

func (r *RepositoryImpl[T]) SelectOne(ctx context.Context, condition T, options ...DBOption) (T, error) {
	var entity T
	err := r.db(ctx, options...).Where(&condition).First(&entity).Error
	return entity, err
}

func (r *RepositoryImpl[T]) SelectList(ctx context.Context, condition T, options ...DBOption) ([]T, error) {
	var entities []T
	err := r.db(ctx, options...).Where(&condition).Find(&entities).Error
	return entities, err
}

func (r *RepositoryImpl[T]) SelectByID(ctx context.Context, id any, options ...DBOption) (T, error) {
	var entity T
	err := r.db(ctx, options...).First(&entity, id).Error
	return entity, err
}

func (r *RepositoryImpl[T]) SelectByIDs(ctx context.Context, ids []any, options ...DBOption) ([]T, error) {
	var entities []T
	err := r.db(ctx, options...).Find(&entities, ids).Error
	return entities, err
}

func (r *RepositoryImpl[T]) SelectCount(ctx context.Context, condition T, options ...DBOption) (int64, error) {
	var count int64
	var model T
	err := r.db(ctx, options...).Model(&model).Where(&condition).Count(&count).Error
	return count, err
}
