package utils

type AdvancedArray[T any] struct {
	index int
	arr   []T
}

func NewAdvancedArray[T any](arr []T) *AdvancedArray[T] {
	return &AdvancedArray[T]{
		arr: arr,
	}
}

func (a *AdvancedArray[T]) Next() (T, bool) {
	if a.index >= len(a.arr) {
		var zero T
		return zero, false
	}

	val := a.arr[a.index]
	a.index++

	return val, true
}

func (a *AdvancedArray[T]) Back() {
	if a.index > 0 {
		a.index--
	}
}
