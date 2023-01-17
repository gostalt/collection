package collection

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/gostalt/collection/join"
)

type collection[T comparable] struct {
	contents []T
}

// Make returns a new empty collection of type T.
func Make[T comparable]() collection[T] {
	return collection[T]{
		contents: []T{},
	}
}

// From returns a new collection from the provided slice.
func From[T comparable](slice []T) collection[T] {
	return collection[T]{
		contents: slice,
	}
}

// All returns the underlying data for the collection.
func (c collection[T]) All() []T {
	return c.contents
}

// Slice returns the underlying data for the collection.
//
// An alias of `All`.
func (c collection[T]) Slice() []T {
	return c.All()
}

// Filter uses the provided predicate to filter the collection, keeping only the
// items for which the predicate returns true.
func (c collection[T]) Filter(predicate func(i int, v T) bool) collection[T] {
	new := Make[T]()

	for i, v := range c.All() {
		if predicate(i, v) {
			new.contents = append(new.All(), v)
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

// SafeFirst works in the same way as First, but returns a collection.ErrNoItem
// if no item was found in the collection (i.e., the collection was empty).
func (c collection[T]) SafeFirst() (T, error) {
	return c.SafeAt(0)
}

// Last returns the last item in the collection. If the collection is empty, a zero
// value of the underlying collection type is returned.
func (c collection[T]) Last() T {
	v, _ := c.SafeLast()
	return v
}

// SafeLast works in the same was as Last, but returns a collection.ErrNoItem if
// no item was found in the collection (i.e., the collection was empty).
func (c collection[T]) SafeLast() (T, error) {
	if c.Empty() {
		return *new(T), ErrNoItem
	}

	return c.All()[len(c.All())-1], nil
}

// FirstWhere returns the first item from the collection that matches the provided
// predicate.
func (c collection[T]) FirstWhere(predicate func(i int, value T) bool) T {
	for i, v := range c.All() {
		if predicate(i, v) {
			return v
		}
	}

	return *new(T)
}

// SafeFirstWhere?

// LastWhere?

// SafeLastWhere?

// Has returns true if the collection contains any item that matches the provided
// predicate. If no nothing matches, or collection is empty, false is returned.
func (c collection[T]) Has(predicate func(i int, value T) bool) bool {
	for i, v := range c.All() {
		if predicate(i, v) {
			return true
		}
	}

	return false
}

// Has no returns true if the collection does not contain an item that matches the
// provided predicate. If nothing matches, or the collection is empty, true is
// returned.
func (c collection[T]) HasNo(predicate func(i int, value T) bool) bool {
	for i, v := range c.All() {
		if predicate(i, v) {
			return false
		}
	}

	return true
}

// Count returns the total length of the collection.
func (c collection[T]) Count() int {
	return len(c.All())
}

// CountWhere returns the number of items in the collection that match the given
// predicate.
func (c collection[T]) CountWhere(predicate func(i int, value T) bool) int {
	count := 0

	for i, v := range c.All() {
		if predicate(i, v) {
			count++
		}
	}

	return count
}

// Append adds the given values to the end of the collection.
func (c collection[T]) Append(value ...T) collection[T] {
	c.contents = append(c.All(), value...)

	return c
}

// Prepend adds the given values to the start of the collection.
func (c collection[T]) Prepend(value ...T) collection[T] {
	return From(value).Append(c.All()...)
}

// At returns the item at the given index. If the index does not exist in the
// collection, a zero value is returned.
func (c collection[T]) At(i int) T {
	v, _ := c.SafeAt(i)
	return v
}

// SafeAt returns the item at the given index. If the index does not exist in the
// collection, a zero value is returned along with collection.ErrNoItem.
func (c collection[T]) SafeAt(i int) (T, error) {
	if c.Empty() || c.Count() < i {
		return *new(T), ErrNoItem
	}

	return c.All()[i], nil
}

// Chan returns a readonly channel for consuming values from the collection.
func (c collection[T]) Chan() <-chan T {
	ch := make(chan T)

	go func(ch chan<- T, c collection[T]) {
		for i := 0; i < c.Count(); i++ {
			ch <- c.At(i)
		}
	}(ch, c)

	return ch
}

// Concat appends the given collection's values to the end of the existing
// collection.
func (c collection[T]) Concat(val collection[T]) collection[T] {
	return c.Append(val.All()...)
}

// Chunk breaks the collection into smaller slices of a given size.
//
// Limitations with generics means it is not possible to return a collection of
// collection (i.e., collection[collection[T]]). If you wish to continue working
// with collections with the returned chunks, you'll need to use From to turn
// them back into collections.
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

// Unique returns all the unique items from the collection.
func (c collection[T]) Unique() collection[T] {
	new := Make[T]()

	for _, v := range c.All() {
		if new.HasNo(func(i int, value T) bool {
			return value == v
		}) {
			new.contents = append(new.contents, v)
		}
	}

	return new
}

// Map iterates through each item of the collection and uses the given function
// to transform the item.
func (c collection[T]) Map(fn func(i int, value T) T) collection[T] {
	new := Make[T]()

	for i, v := range c.contents {
		new = new.Append(fn(i, v))
	}

	return new
}

// Pop removes and returns items from the end of the collection.
func (c *collection[T]) Pop(count int) collection[T] {
	split := c.Split(c.Count() - count)
	c.contents = c.All()[:c.Count()-count]

	return split[1]
}

// Before returns the items before the provided index.
func (c collection[T]) Before(i int) collection[T] {
	return From(c.All()[:i])
}

// After returns the items after the provided index.
func (c collection[T]) After(i int) collection[T] {
	return From(c.All()[i:])
}

// Split returns two collections, split on the given index.
func (c collection[T]) Split(i int) []collection[T] {
	return []collection[T]{
		c.Before(i),
		c.After(i),
	}
}

// Diff returns the values from the original collection that are not found in the
// given collection.
func (c collection[T]) Diff(comp collection[T]) collection[T] {
	return c.Filter(func(i int, v T) bool {
		return comp.HasNo(func(i int, value T) bool {
			return value == v
		})
	})
}

// Join joins the collection's items using the provided join.Method. If a Final
// value is provided to the join.Method, it is used to join the final two elements.
//
// Two standard join methods are provided by the join package:
//   - join.CommaSeparatedJoin, which would result in: "1, 2, 3"
//   - join.ListJoin, which would result in "1, 2 and 3"
func (c collection[T]) Join(format join.Method) string {
	resp := ""

	for i, v := range c.All() {
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

	return From(c.All()[:count])
}

// Empty returns true if the collection contains no items.
func (c collection[T]) Empty() bool {
	if c.Count() == 0 {
		return true
	}

	return false
}

// NotEmpty returns true if the collection contains items.
func (c collection[T]) NotEmpty() bool {
	return !c.Empty()
}

// Random uses the provided *rand.Rand to pick the given number of items from the
// collection. Elements can be picked more than once. Because random elements
// are picked, the count parameter can be larger than the total size of
// the collection.
func (c collection[T]) Random(r *rand.Rand, count int) collection[T] {
	new := From(make([]T, count))
	for i := range new.All() {
		new.Set(i, c.random(r))
	}

	return new
}

// random returns a single item from the underlying contents of the collection.
func (c collection[T]) random(r *rand.Rand) T {
	return c.At(r.Intn(c.Count()))
}

// Set updates the value at the given index to value. If the given index is out of
// range for the collection's underlying slice, the slice is expanded to allow
// the value to be set. Use `SafeSet` to prevent this behaviour and return
// an error if out of bounds.
func (c *collection[T]) Set(index int, value T) {
	if c.Count() >= index {
		c.contents[index] = value
		return
	}

	new := make([]T, index+1)
	for i, v := range c.All() {
		new[i] = v
	}

	new[index] = value
	c.contents = new
}

// SafeSet updates the value at the given index to value. If the given index is
// out of range for the collection's underlying slice, an error is returned and
// the collection is not modified.
func (c *collection[T]) SafeSet(index int, value T) error {
	if c.Count() < index {
		return ErrIndexOutOfRange
	}

	c.contents[index] = value

	return nil
}

// Each iterates over each item inside the collection and passes the index and value
// to the provided func.
func (c collection[T]) Each(fn func(i int, value T)) {
	for i, v := range c.All() {
		fn(i, v)
	}
}

// EachCtx iterates over each item inside the collection and passes the index and
// value to the provided func. If the given context is Done, the iteration stops.
func (c collection[T]) EachCtx(ctx context.Context, fn func(i int, value T)) {
	for i, v := range c.All() {
		select {
		case <-ctx.Done():
			return
		default:
			fn(i, v)
		}
	}
}

// Every returns true if all items inside the collection satisfy the given predicate.
func (c collection[T]) Every(predicate func(i int, value T) bool) bool {
	for i, v := range c.contents {
		if !predicate(i, v) {
			return false
		}
	}

	return true
}

// Reverse returns a new collection with the values in reverse order.
func (c collection[T]) Reverse() collection[T] {
	new := From(make([]T, c.Count()))

	c.Each(func(i int, value T) {
		new.Set(c.Count()-1-i, value)
	})

	return new
}

// Search returns the index of the first item that matches the given predicate.
// If no item is found, -1 is returned.
func (c collection[T]) Search(fn func(i int, value T) bool) int {
	v, _ := c.SafeSearch(fn)
	return v
}

// SafeSearch returns the index of the first item that matches the given predicate.
// If no item is found, -1 and collection.ErrNoItem is returned.
func (c collection[T]) SafeSearch(fn func(i int, value T) bool) (int, error) {
	for i, v := range c.All() {
		if fn(i, v) {
			return i, nil
		}
	}

	return -1, ErrNoItem
}
