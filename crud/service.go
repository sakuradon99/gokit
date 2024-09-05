package crud

import "context"

type Service[T any] interface {
	Save(ctx context.Context, entity T) error
	RemoveByID(ctx context.Context, id any, options ...RemoveOption) error
	RemoveByIDs(ctx context.Context, ids []any, options ...RemoveOption) error
	UpdateByID(ctx context.Context, entity T) error
	GetByID(ctx context.Context, id any, options ...QueryOption) (T, error)
	ListByIDs(ctx context.Context, ids []any, options ...QueryOption) ([]T, error)
}

type ServiceImpl[T any] struct {
	repository RepositoryImpl[T]
}

func (s *ServiceImpl[T]) Save(ctx context.Context, entity T) error {
	return s.repository.Insert(ctx, entity)
}

func (s *ServiceImpl[T]) RemoveByID(ctx context.Context, id any, options ...RemoveOption) error {
	return s.repository.DeleteByID(ctx, id, options...)
}

func (s *ServiceImpl[T]) RemoveByIDs(ctx context.Context, ids []any, options ...RemoveOption) error {
	return s.repository.DeleteByIDs(ctx, ids, options...)
}

func (s *ServiceImpl[T]) UpdateByID(ctx context.Context, entity T) error {
	return s.repository.UpdateByID(ctx, entity)
}

func (s *ServiceImpl[T]) GetByID(ctx context.Context, id any, options ...QueryOption) (T, error) {
	return s.repository.SelectByID(ctx, id, options...)
}

func (s *ServiceImpl[T]) ListByIDs(ctx context.Context, ids []any, options ...QueryOption) ([]T, error) {
	return s.repository.SelectByIDs(ctx, ids, options...)
}
