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
