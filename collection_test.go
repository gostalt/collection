package collection_test

import (
	"testing"

	"github.com/gostalt/collection"
	"github.com/stretchr/testify/assert"
)

func TestCanAccessUnderlyingSlice(t *testing.T) {
	v := collection.From([]int{1, 2, 3})

	assert.Equal(t, []int{1, 2, 3}, v.All())
}

func TestFilter(t *testing.T) {
	v := collection.
		From([]int{1, 1, 2, 2}).
		Filter(func(i int, value int) bool {
			return value == 2
		})

	assert.Equal(t, []int{2, 2}, v.All())
}

func TestFirst(t *testing.T) {
	v := collection.From([]int{3, 2, 1}).First()
	assert.Equal(t, 3, v)

	v = collection.From([]int{}).First()
	assert.Equal(t, 0, v)
}

func TestSafeFirst(t *testing.T) {
	_, err := collection.From([]string{}).SafeFirst()
	assert.ErrorIs(t, err, collection.ErrNoItem)

	v, err := collection.From([]string{"hello", "world"}).SafeFirst()
	assert.NoError(t, err)
	assert.Equal(t, "hello", v)
}

func TestLast(t *testing.T) {
	v := collection.From([]int{1, 2, 3}).Last()
	assert.Equal(t, 3, v)

	v = collection.From([]int{}).Last()
	assert.Equal(t, 0, v)
}

func TestSafeLast(t *testing.T) {
	_, err := collection.From([]string{}).SafeLast()
	assert.ErrorIs(t, err, collection.ErrNoItem)

	v, err := collection.From([]string{"hello", "world"}).SafeLast()
	assert.NoError(t, err)
	assert.Equal(t, "world", v)
}

func TestFirstWhereReturnsFirstSuccessfulMatch(t *testing.T) {
	v := collection.From([]int{1, 3, 5, 7, 8}).FirstWhere(func(i int, value int) bool {
		return value%2 == 0
	})

	assert.Equal(t, 8, v)
}

func TestFirstWhereOfEmptyReturnsZeroValue(t *testing.T) {
	v := collection.From([]int{}).FirstWhere(func(i int, value int) bool {
		return value%2 == 0
	})

	assert.Equal(t, 0, v)
}

func TestHas(t *testing.T) {
	success := collection.From([]int{1, 3, 5}).Has(func(i int, value int) bool {
		return value == 3
	})

	assert.Equal(t, true, success)

	failure := collection.From([]string{"hello", "world"}).Has(func(i int, value string) bool {
		return value == "mars"
	})

	assert.Equal(t, false, failure)
}

func TestHasNo(t *testing.T) {
	failure := collection.From([]int{1, 2, 3}).HasNo(func(i int, value int) bool {
		return value == 3
	})

	assert.Equal(t, false, failure)

	success := collection.From([]string{"thomas", "smith"}).HasNo(func(i int, value string) bool {
		return value == "another"
	})

	assert.Equal(t, true, success)
}

func TestCount(t *testing.T) {
	count := collection.From([]int{1, 3, 5}).Count()
	assert.Equal(t, 3, count)

	count = collection.From([]string{"single"}).Count()
	assert.Equal(t, 1, count)
}

func TestCountWhere(t *testing.T) {
	count := collection.From([]int{1, 2, 3}).CountWhere(func(i int, value int) bool {
		return value%2 == 1
	})
	assert.Equal(t, 2, count)

	count = collection.From([]string{"single", "double", "something"}).CountWhere(func(i int, value string) bool {
		return value[0] == 's'
	})
	assert.Equal(t, 2, count)
}

func TestAppend(t *testing.T) {
	v := collection.From([]int{1, 2, 3})
	new := v.Append(4, 5)

	assert.Equal(t, 3, v.Count())
	assert.Equal(t, 5, new.Count())
	assert.Equal(t, []int{1, 2, 3}, v.All())
	assert.Equal(t, []int{1, 2, 3, 4, 5}, new.All())
}

func TestAt(t *testing.T) {
	v := collection.From([]string{"first", "second", "third"})
	assert.Equal(t, "first", v.At(0))
	assert.Equal(t, "second", v.At(1))
	assert.Equal(t, "third", v.At(2))
}

func TestSafeAt(t *testing.T) {
	c := collection.From([]string{"first", "second", "third"})
	v, err := c.SafeAt(0)
	assert.Equal(t, "first", v)
	assert.NoError(t, err)

	v, err = c.SafeAt(4)
	assert.Equal(t, "", v)
	assert.ErrorIs(t, err, collection.ErrNoItem)
}

func TestChan(t *testing.T) {
	var vals []int

	col := collection.From([]int{1, 2, 3, 4})
	ch := col.Chan()

	for i := 0; i < col.Count(); i++ {
		select {
		case v := <-ch:
			vals = append(vals, v)
		}
	}

	assert.Equal(t, col.All(), vals)
}

func TestConcat(t *testing.T) {
	first := collection.From([]int{1, 2, 3})
	second := collection.From([]int{4, 5, 6})

	new := first.Concat(second)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, new.All())
}

func TestChunk(t *testing.T) {
	orig := collection.From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	chunks := orig.Chunk(4)

	assert.Equal(t, orig.All()[0:4], chunks[0])
	assert.Equal(t, orig.All()[4:8], chunks[1])
	assert.Equal(t, orig.All()[8:10], chunks[2])
}

func TestUnique(t *testing.T) {
	orig := collection.From([]int{1, 2, 3, 1, 1, 2, 2, 3, 3}).Unique()

	assert.Equal(t, []int{1, 2, 3}, orig.All())
}

func TestMap(t *testing.T) {
	doubled := collection.From([]int{1, 2, 3, 4, 5}).Map(func(i int, value int) int {
		return value * 2
	})

	assert.Equal(t, []int{2, 4, 6, 8, 10}, doubled.All())

	pluralised := collection.From([]string{"lion", "tiger", "bear"}).Map(func(i int, value string) string {
		return value + "s"
	})

	assert.Equal(t, []string{"lions", "tigers", "bears"}, pluralised.All())
}
