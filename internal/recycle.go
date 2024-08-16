package internal

type Reusable interface {
	Recycle()
}

func RecycleVal(v any) {
	if r, ok := v.(Reusable); ok {
		r.Recycle()
	}
}
