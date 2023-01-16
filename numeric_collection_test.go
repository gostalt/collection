package collection_test

import (
	"testing"

	"github.com/gostalt/collection"
	"github.com/stretchr/testify/assert"
)

func TestAverage(t *testing.T) {
	avg := collection.FromNumeric([]int{1, 2, 3, 4}).Average()
	assert.Equal(t, 2.5, avg)

	f32 := collection.FromNumeric([]int{1, 2, 3, 4, 5}).Average32()
	assert.Equal(t, float32(3), f32)

	f64 := collection.FromNumeric([]int{1, 2, 3, 4, 5, 6}).Average64()
	assert.Equal(t, 3.5, f64)
}

func TestMin(t *testing.T) {
	min := collection.FromNumeric([]int{3, 2, 8, 1, 2}).Min()

	assert.Equal(t, 1, min)
}

func TestMax(t *testing.T) {
	max := collection.FromNumeric([]int{3, 2, 8, 1, 2}).Max()

	assert.Equal(t, 8, max)
}

func TestSum(t *testing.T) {
	sum := collection.FromNumeric([]int{1, 2, 3, 4}).Sum()

	assert.Equal(t, 10, sum)
}
