package collection

type collection[T any] struct {
	contents []T
}

func From[T any](slice []T) collection[T] {
	return collection[T]{
		contents: slice,
	}
}

func (c collection[T]) All() []T {
	return c.contents
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
