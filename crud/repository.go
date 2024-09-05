package crud

import (
	"context"
	"github.com/sakuradon99/gokit/db"
	"gorm.io/gorm/clause"
)

type Repository[T any] interface {
	Insert(ctx context.Context, entity T) error
	DeleteByID(ctx context.Context, id any, options ...RemoveOption) error
	DeleteByIDs(ctx context.Context, ids []any, options ...RemoveOption) error
	UpdateByID(ctx context.Context, entity T) error
	SelectByID(ctx context.Context, id any, options ...QueryOption) (T, error)
	SelectByIDs(ctx context.Context, ids []any, options ...QueryOption) ([]T, error)
}

type RepositoryImpl[T any] struct {
	dbm db.Manager `inject:""`
}

func (r *RepositoryImpl[T]) Insert(ctx context.Context, entity T) error {
	return r.dbm.DB(ctx).Create(&entity).Error
}

func (r *RepositoryImpl[T]) DeleteByID(ctx context.Context, id any, options ...RemoveOption) error {
	opt := buildRemoveOptions(options...)
	d := r.dbm.DB(ctx)
	if opt.RemoveAssociations {
		d = d.Select(clause.Associations)
	}

	var entity T
	return r.dbm.DB(ctx).Delete(&entity, id).Error
}

func (r *RepositoryImpl[T]) DeleteByIDs(ctx context.Context, ids []any, options ...RemoveOption) error {
	opt := buildRemoveOptions(options...)
	d := r.dbm.DB(ctx)
	if opt.RemoveAssociations {
		d = d.Select(clause.Associations)
	}

	var entity T
	return r.dbm.DB(ctx).Delete(&entity, ids).Error
}

func (r *RepositoryImpl[T]) UpdateByID(ctx context.Context, entity T) error {
	return r.dbm.DB(ctx).Save(&entity).Error
}

func (r *RepositoryImpl[T]) SelectByID(ctx context.Context, id any, options ...QueryOption) (T, error) {
	opt := buildQueryOptions(options...)
	d := r.dbm.DB(ctx)
	if opt.PreloadAssociations {
		d = d.Preload(clause.Associations)
	}

	var entity T
	err := d.First(&entity, id).Error
	return entity, err
}
func (r *RepositoryImpl[T]) SelectByIDs(ctx context.Context, ids []any, options ...QueryOption) ([]T, error) {
	opt := buildQueryOptions(options...)
	d := r.dbm.DB(ctx)
	if opt.PreloadAssociations {
		d = d.Preload(clause.Associations)
	}

	var entities []T
	err := r.dbm.DB(ctx).Find(&entities, ids).Error
	return entities, err
}
