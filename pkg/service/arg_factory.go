package service

type ArgFactory interface {
	New() any
	Recycle(any)
}

type TArgsFactory[T any] struct{}

func (f TArgsFactory[T]) New() any {
	return new(T)
}

func (f TArgsFactory[T]) Recycle(v any) {}
