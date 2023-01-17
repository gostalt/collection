package collection

type i interface {
	int | int8 | int16 | int32 | int64
}

type f interface {
	float32 | float64
}

type numeric interface {
	i | f
}

type numericCollection[T numeric] struct {
	collection[T]
}

// FromNumeric creates a new numericCollection from the provided slice.
func FromNumeric[T numeric](slice []T) numericCollection[T] {
	c := From(slice)

	return numericCollection[T]{
		c,
	}
}

// FromRange creates a new numericCollection of ints between the first and last
// parameters, inclusive, with a step of 1. For example, FromRange(1, 3) would
// return a new numericCollection with an underlying slice of []int{1, 2, 3}.
//
// Specifying a range where the first parameter is larger than the last parameter
// results in a decrementing range.
func FromRange(first int, last int) numericCollection[int] {
	desc := false
	if first > last {
		first, last = last, first
		desc = true
	}

	len := last - first + 1
	c := FromNumeric(make([]int, len))
	for i := 0; i < len; i++ {
		c.Set(i, first+i)
	}

	if desc {
		return FromNumeric(c.Reverse().All())
	}

	return c
}

// Average returns a mean average of the collection.
func (c numericCollection[T]) Average() float64 {
	return c.Average64()
}

// Average32 returns a mean average of the collection, as a float32 type.
func (c numericCollection[T]) Average32() float32 {
	var sum T
	count := len(c.contents)

	for _, v := range c.contents {
		sum = sum + v
	}

	return float32(sum) / float32(count)
}

// Average64 returns a mean average of the collection, as a float64 type.
func (c numericCollection[T]) Average64() float64 {
	var sum T
	count := len(c.contents)

	for _, v := range c.contents {
		sum = sum + v
	}

	return float64(sum) / float64(count)
}

// Min returns the smallest number in the collection. If the collection is empty,
// a zero value is returned.
func (c numericCollection[T]) Min() T {
	if c.Empty() {
		return 0
	}

	min := c.At(0)

	for _, v := range c.contents {
		if v < min {
			min = v
		}
	}

	return min
}

// Max returns the largest number in the collection. If the collection is empty,
// a zero value is returned.
func (c numericCollection[T]) Max() T {
	if c.Empty() {
		return 0
	}

	max := c.At(0)

	for _, v := range c.contents {
		if v > max {
			max = v
		}
	}

	return max
}

// Sum returns the total value of all of the values inside the collection.
func (c numericCollection[T]) Sum() T {
	var total T = 0

	for _, v := range c.contents {
		total = total + v
	}

	return total
}
