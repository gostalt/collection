package collection

import (
	"math"
)

type collection[T comparable] struct {
	contents []T
}

func From[T comparable](slice []T) collection[T] {
	return collection[T]{
		contents: slice,
	}
}

func (c collection[T]) All() []T {
	return c.contents
}

// Slice returns the underlying data for the collection.
//
// An alias of `All`.
func (c collection[T]) Slice() []T {
	return c.All()
}

func (c collection[T]) Filter(predicate func(i int, v T) bool) collection[T] {
	new := collection[T]{}

	for i, v := range c.contents {
		if predicate(i, v) {
			new.contents = append(new.contents, v)
		}
	}

	return new
}

// First returns the first item in the collection. If the collection is empty, a
// zero value of the underlying collection type is returned.
func (c collection[T]) First() T {
	v, _ := c.SafeFirst()
	return v
}

// SafeFirst works in the same way as `First`, but returns a `collection.ErrNoItem`
// if no item was found in the collection (i.e., the collection was empty).
func (c collection[T]) SafeFirst() (T, error) {
	if len(c.contents) == 0 {
		return *new(T), ErrNoItem
	}

	return c.contents[0], nil
}

func (c collection[T]) Last() T {
	v, _ := c.SafeLast()
	return v
}

func (c collection[T]) SafeLast() (T, error) {
	if len(c.contents) == 0 {
		return *new(T), ErrNoItem
	}

	return c.contents[len(c.contents)-1], nil
}

// FirstWhere
func (c collection[T]) FirstWhere(predicate func(i int, value T) bool) T {
	for i, v := range c.contents {
		if predicate(i, v) {
			return v
		}
	}

	return *new(T)
}

// SafeFirstWhere?

// LastWhere?

// SafeLastWhere?

func (c collection[T]) Has(predicate func(i int, value T) bool) bool {
	for i, v := range c.contents {
		if predicate(i, v) {
			return true
		}
	}

	return false
}

func (c collection[T]) HasNo(predicate func(i int, value T) bool) bool {
	for i, v := range c.contents {
		if predicate(i, v) {
			return false
		}
	}

	return true
}

func (c collection[T]) Count() int {
	return len(c.contents)
}

func (c collection[T]) CountWhere(predicate func(i int, value T) bool) int {
	count := 0

	for i, v := range c.contents {
		if predicate(i, v) {
			count++
		}
	}

	return count
}

func (c collection[T]) Append(value ...T) collection[T] {
	c.contents = append(c.contents, value...)

	return c
}

func (c collection[T]) At(i int) T {
	v, _ := c.SafeAt(i)
	return v
}

func (c collection[T]) SafeAt(i int) (T, error) {
	if len(c.contents) < i {
		return *new(T), ErrNoItem
	}

	return c.contents[i], nil
}

func (c collection[T]) Chan() <-chan T {
	ch := make(chan T)

	go func(ch chan<- T, c collection[T]) {
		for i := 0; i < c.Count(); i++ {
			ch <- c.At(i)
		}
	}(ch, c)

	return ch
}

func (c collection[T]) Concat(val collection[T]) collection[T] {
	new := c.Append(val.All()...)

	return new
}

func (c collection[T]) Chunk(per int) [][]T {
	count := int(math.Ceil(float64(len(c.contents)) / float64(per)))
	chunks := make([][]T, count)

	for i := range chunks {
		chunks[i] = make([]T, 0, per)
		for j := 0; j < per; j++ {
			offset := i*per + j
			if offset >= c.Count() {
				break
			}

			chunks[i] = append(chunks[i], c.At(offset))
		}
	}
	return chunks
}

func (c collection[T]) Unique() collection[T] {
	new := collection[T]{}

	for _, v := range c.contents {
		if new.HasNo(func(i int, value T) bool {
			return value == v
		}) {
			new.contents = append(new.contents, v)
		}
	}

	return new
}

func (c collection[T]) Map(fn func(i int, value T) T) collection[T] {
	new := collection[T]{}

	for i, v := range c.contents {
		new = new.Append(fn(i, v))
	}

	return new
}

func (c *collection[T]) Pop(count int) collection[T] {
	split := c.Split(c.Count() - count)
	c.contents = c.contents[:len(c.contents)-count]

	return split[1]
}

func (c collection[T]) Split(i int) []collection[T] {
	one := From(c.contents[:i])
	two := From(c.contents[i:])

	return []collection[T]{one, two}
}

func (c collection[T]) Diff(comp collection[T]) collection[T] {
	return c.Filter(func(i int, v T) bool {
		return comp.HasNo(func(i int, value T) bool {
			return value == v
		})
	})
}
