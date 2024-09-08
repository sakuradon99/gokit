package crud

import (
	"context"
)

type Service[T any] interface {
	Save(ctx context.Context, entity T) error
	SaveBatch(ctx context.Context, entities []T) error
	Remove(ctx context.Context, condition T) error
	RemoveByID(ctx context.Context, id any) error
	RemoveByIDs(ctx context.Context, ids []any) error
	Update(ctx context.Context, entity T) error
	UpdateBatch(ctx context.Context, entities []T) error
	GetOne(ctx context.Context, condition T) (T, error)
	List(ctx context.Context, condition T) ([]T, error)
	GetByID(ctx context.Context, id any) (T, error)
	ListByIDs(ctx context.Context, ids []any) ([]T, error)
	Count(ctx context.Context, condition T) (int64, error)
}

type ServiceImpl[T any] struct {
	repository RepositoryImpl[T]
}

func (s *ServiceImpl[T]) Save(ctx context.Context, entity T) error {
	return s.repository.Insert(ctx, entity)
}

func (s *ServiceImpl[T]) SaveBatch(ctx context.Context, entities []T) error {
	return s.repository.InsertBatch(ctx, entities)
}

func (s *ServiceImpl[T]) Remove(ctx context.Context, condition T) error {
	return s.repository.Delete(ctx, condition)
}

func (s *ServiceImpl[T]) RemoveByID(ctx context.Context, id any) error {
	return s.repository.DeleteByID(ctx, id)
}

func (s *ServiceImpl[T]) RemoveByIDs(ctx context.Context, ids []any) error {
	return s.repository.DeleteByIDs(ctx, ids)
}

func (s *ServiceImpl[T]) Update(ctx context.Context, entity T) error {
	return s.repository.Update(ctx, entity)
}

func (s *ServiceImpl[T]) UpdateBatch(ctx context.Context, entities []T) error {
	return s.repository.UpdateBatch(ctx, entities)
}

func (s *ServiceImpl[T]) GetOne(ctx context.Context, condition T) (T, error) {
	return s.repository.SelectOne(ctx, condition)
}

func (s *ServiceImpl[T]) List(ctx context.Context, condition T) ([]T, error) {
	return s.repository.SelectList(ctx, condition)
}

func (s *ServiceImpl[T]) GetByID(ctx context.Context, id any) (T, error) {
	return s.repository.SelectByID(ctx, id)
}

func (s *ServiceImpl[T]) ListByIDs(ctx context.Context, ids []any) ([]T, error) {
	return s.repository.SelectByIDs(ctx, ids)
}

func (s *ServiceImpl[T]) Count(ctx context.Context, condition T) (int64, error) {
	return s.repository.SelectCount(ctx, condition)
}
