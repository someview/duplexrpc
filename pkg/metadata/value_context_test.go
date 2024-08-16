package metadata

import (
	"context"
	"testing"
)

func BenchmarkWithValue(b *testing.B) {
	task := func(ctx context.Context, with func(context.Context, any, any) context.Context) context.Context {
		for i := 0; i < 10; i++ {
			ctx = with(ctx, i, i)
		}
		return ctx
	}

	b.Run("官方WithValue", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctx := task(context.TODO(), context.WithValue)
			RecycleContext(ctx)
		}
	})

	b.Run("官方WithValue-Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			b.ResetTimer()
			for pb.Next() {
				ctx := task(context.TODO(), context.WithValue)
				RecycleContext(ctx)
			}
		})
	})

	b.Run("valueContext", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctx := task(context.TODO(), WithValue)
			RecycleContext(ctx)
		}
	})

	b.Run("valueContext-Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			b.ResetTimer()
			for pb.Next() {
				ctx := task(context.TODO(), WithValue)
				RecycleContext(ctx)
			}
		})
	})

}
