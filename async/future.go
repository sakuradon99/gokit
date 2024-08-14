package async

import (
	"context"
	"sync"
)

type futureErr interface {
	Err(ctx context.Context) error
}

type Future[T any] interface {
	futureErr
	Set(value T, err error)
	Value(ctx context.Context) T
	IsDone() bool
}

type future[T any] struct {
	value  T
	err    error
	done   bool
	doneCh chan struct{}
	mu     sync.Mutex
	once   sync.Once
}

func (f *future[T]) Set(value T, err error) {
	f.once.Do(func() {
		f.mu.Lock()
		defer f.mu.Unlock()

		if f.done {
			return
		}

		f.value = value
		f.err = err
		f.done = true
		close(f.doneCh)
	})
}

func (f *future[T]) Value(ctx context.Context) T {
	select {
	case <-ctx.Done():
		var defaultVal T
		return defaultVal
	case <-f.doneCh:
		return f.value
	}
}

func (f *future[T]) IsDone() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.done
}

func (f *future[T]) Err(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-f.doneCh:
		return f.err
	}
}

func NewFuture[T any]() Future[T] {
	return &future[T]{
		doneCh: make(chan struct{}),
		once:   sync.Once{},
	}
}

type FuturePanicError struct {
	recover any
}

func (f FuturePanicError) Error() string {
	return "panic in future"
}

func (f FuturePanicError) Recover() any {
	return f.recover
}

func GoFuture[T any](fn func() (T, error)) Future[T] {
	f := NewFuture[T]()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				var defaultVal T
				var err error

				if e, ok := r.(error); ok {
					err = e
				} else {
					err = FuturePanicError{recover: r}
				}

				f.Set(defaultVal, err)
			}
		}()
		value, err := fn()
		f.Set(value, err)
	}()
	return f
}

func GetFirstFutureError(ctx context.Context, futures ...futureErr) error {
	for _, f := range futures {
		if err := f.Err(ctx); err != nil {
			return err
		}
	}
	return nil
}
