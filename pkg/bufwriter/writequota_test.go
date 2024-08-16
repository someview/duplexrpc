package bufwriter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteQuota_GetWithCtx(t *testing.T) {
	t.Run("buffer is enough", func(t *testing.T) {
		closeCh := make(chan struct{}, 1)
		bufSize := 1024
		qw := newWriteQuota(int32(bufSize), closeCh)
		err := qw.GetWithCtx(context.Background(), 1024)
		assert.Nil(t, err)
	})
	t.Run("buffer is not enough", func(t *testing.T) {
		closeCh := make(chan struct{}, 1)
		bufSize := 1024
		qw := newWriteQuota(int32(bufSize), closeCh)
		assert.True(t, qw.Available() > 0)
		err := qw.GetWithCtx(context.TODO(), bufSize*2)
		assert.Nil(t, err, "配额大于0,放行当前额度")
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		assert.True(t, qw.Available() < 0)
		err = qw.GetWithCtx(ctx, bufSize*2)
		assert.NotNil(t, err, "配额小于0,当前请求被拒绝")
		ok := qw.TryGet(bufSize * 2)
		assert.False(t, ok)
	})
}

func TestWriteQuota_Close(t *testing.T) {
	closeCh := make(chan struct{}, 1)
	bufSize := -1
	qw := newWriteQuota(int32(bufSize), closeCh)
	close(closeCh)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := qw.GetWithCtx(context.TODO(), bufSize*2)
		assert.NotNil(t, err)
	}()
	go func() {
		defer wg.Done()
		err := qw.GetWithCtx(context.TODO(), bufSize*2)
		assert.NotNil(t, err)
	}()
	wg.Wait()
	t.Log("close成功通知所有因为没有获取到配额的协程")
}

func TestWriteQuota_Release(t *testing.T) {
	closeCh := make(chan struct{}, 1)
	bufSize := -1
	qw := newWriteQuota(int32(bufSize), closeCh)
	close(closeCh)
	quota := 1000 - bufSize
	qw.Release(quota)
	assert.Equal(t, qw.Available(), 1000)
	var wg sync.WaitGroup
	routineNum := 10
	quotaPerRoutine := quota / routineNum
	wg.Add(routineNum)
	for i := 0; i < routineNum; i++ {
		go func() {
			defer wg.Done()
			err := qw.GetWithCtx(context.TODO(), quotaPerRoutine)
			assert.Nil(t, err)
		}()
	}
	wg.Wait()
	assert.Equal(t, qw.Available(), 0)
}
