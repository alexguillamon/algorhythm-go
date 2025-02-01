package orderedset

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type Set[T comparable] struct {
	set   mapset.Set[T] // For fast membership tests.
	order []T           // To keep track of insertion order.
}

func NewSet[T comparable](items ...T) Set[T] {
	os := Set[T]{
		set:   mapset.NewSet[T](),
		order: make([]T, 0, len(items)),
	}
	for _, item := range items {
		os.Add(item)
	}
	return os
}

func (os *Set[T]) Add(items ...T) {
	for _, item := range items {
		if !os.set.Contains(item) {
			os.set.Add(item)
			os.order = append(os.order, item)
		}
	}
}

func (os *Set[T]) Contains(item T) bool {
	return os.set.Contains(item)
}

func (os *Set[T]) Remove(item T) {
	if os.set.Contains(item) {
		os.set.Remove(item)
		// Remove item from the order slice (linear time).
		for i, v := range os.order {
			if v == item {
				os.order = append(os.order[:i], os.order[i+1:]...)
				break
			}
		}
	}
}

func (os *Set[T]) ToSlice() []T {
	return os.order
}

func (os Set[T]) Union(other Set[T]) Set[T] {
	result := NewSet[T]()
	// Add all items from the left set.
	for _, item := range os.order {
		result.Add(item)
	}
	// Add items from the right set only if they are not already included.
	for _, item := range other.order {
		if !result.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func (os Set[T]) Difference(other Set[T]) Set[T] {
	result := NewSet[T]()
	for _, item := range os.order {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}
