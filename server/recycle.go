package server

type recyclable interface {
	Recycle()
}

func recycleValue(v any) {
	if r, ok := v.(recyclable); ok {
		r.Recycle()
	}
}
