package collection

import (
	"fmt"
	"math"

	"github.com/gostalt/collection/join"
)

type collection[T comparable] struct {
	contents []T
}

func Make[T comparable]() collection[T] {
	return collection[T]{
		contents: []T{},
	}
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
	new := Make[T]()

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
	if c.Empty() {
		return *new(T), ErrNoItem
	}

	return c.contents[0], nil
}

func (c collection[T]) Last() T {
	v, _ := c.SafeLast()
	return v
}

func (c collection[T]) SafeLast() (T, error) {
	if c.Empty() {
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

func (c collection[T]) Prepend(value ...T) collection[T] {
	return From(value).Append(c.All()...)
}

func (c collection[T]) At(i int) T {
	v, _ := c.SafeAt(i)
	return v
}

func (c collection[T]) SafeAt(i int) (T, error) {
	if c.Count() < i {
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
	return c.Append(val.All()...)
}

func (c collection[T]) Chunk(per int) [][]T {
	count := int(math.Ceil(float64(c.Count()) / float64(per)))
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
	new := Make[T]()

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
	new := Make[T]()

	for i, v := range c.contents {
		new = new.Append(fn(i, v))
	}

	return new
}

func (c *collection[T]) Pop(count int) collection[T] {
	split := c.Split(c.Count() - count)
	c.contents = c.contents[:c.Count()-count]

	return split[1]
}

func (c collection[T]) Before(i int) collection[T] {
	return From(c.contents[:i])
}

func (c collection[T]) After(i int) collection[T] {
	return From(c.contents[i:])
}

func (c collection[T]) Split(i int) []collection[T] {
	return []collection[T]{
		c.Before(i),
		c.After(i),
	}
}

func (c collection[T]) Diff(comp collection[T]) collection[T] {
	return c.Filter(func(i int, v T) bool {
		return comp.HasNo(func(i int, value T) bool {
			return value == v
		})
	})
}

func (c collection[T]) Join(format join.Method) string {
	resp := ""

	for i, v := range c.contents {
		resp = resp + fmt.Sprintf("%v", v)
		if i == c.Count()-1 {
			continue
		}
		if i == c.Count()-2 {
			if format.Final != "" {
				resp = resp + format.Final
			} else {
				resp = resp + format.Between
			}
			continue
		}

		resp = resp + format.Between
	}

	return resp
}

// FirstX returns the first X items from the collection as a new collection. If
// the collection has fewer than the requested number of items, the original
// collection is returned.
func (c collection[T]) FirstX(count int) collection[T] {
	if c.Count() <= count {
		return c
	}

	return From(c.contents[:count])
}

func (c collection[T]) Empty() bool {
	if c.Count() == 0 {
		return true
	}

	return false
}
