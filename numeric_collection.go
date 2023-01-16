package collection

type numeric interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type numericCollection[T numeric] struct {
	collection[T]
}

func FromNumeric[T numeric](slice []T) numericCollection[T] {
	c := From(slice)

	return numericCollection[T]{
		c,
	}
}

func (c numericCollection[T]) Average() float64 {
	return c.Average64()
}

func (c numericCollection[T]) Average32() float32 {
	var sum T
	count := len(c.contents)

	for _, v := range c.contents {
		sum = sum + v
	}

	return float32(sum) / float32(count)
}

func (c numericCollection[T]) Average64() float64 {
	var sum T
	count := len(c.contents)

	for _, v := range c.contents {
		sum = sum + v
	}

	return float64(sum) / float64(count)
}

func (c numericCollection[T]) Min() T {
	var min T
	if c.Count() == 0 {
		return 0
	}

	min = c.At(0)

	for _, v := range c.contents {
		if v < min {
			min = v
		}
	}

	return min
}

func (c numericCollection[T]) Max() T {
	var max T
	if c.Count() == 0 {
		return 0
	}

	max = c.At(0)

	for _, v := range c.contents {
		if v > max {
			max = v
		}
	}

	return max
}

func (c numericCollection[T]) Sum() T {
	var total T = 0

	for _, v := range c.contents {
		total = total + v
	}

	return total
}
