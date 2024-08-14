package async

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Future(t *testing.T) {
	f1 := GoFuture(func() (int, error) {
		time.Sleep(1 * time.Second)
		return 1, nil
	})
	f2 := GoFuture(func() (int, error) {
		time.Sleep(2 * time.Second)
		return 0, errors.New("error")
	})
	f3 := GoFuture(func() (int, error) {
		panic("panic")
	})

	ctx := context.Background()
	err := GetFirstFutureError(ctx, f1, f2, f3)
	assert.NotNil(t, err)
	assert.Equal(t, 1, f1.Value(ctx))
	assert.True(t, f1.IsDone())
	assert.Equal(t, "error", f2.Err(ctx).Error())
	assert.Equal(t, "panic in future", f3.Err(ctx).Error())
	var e FuturePanicError
	ok := errors.As(f3.Err(ctx), &e)
	assert.True(t, ok)
	assert.Equal(t, "panic", e.Recover())
}
