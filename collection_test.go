package collection_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/gostalt/collection"
	"github.com/gostalt/collection/join"
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

	empty := collection.From([]string{}).Has(func(i int, value string) bool {
		return value == "anything"
	})

	assert.Equal(t, false, empty)
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

	empty := collection.From([]string{}).HasNo(func(i int, value string) bool {
		return value == "anything"
	})

	assert.Equal(t, true, empty)
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

func TestPop(t *testing.T) {
	orig := collection.From([]int{1, 2, 3, 4, 5})
	single := orig.Pop(1)

	assert.Equal(t, []int{1, 2, 3, 4}, orig.All())
	assert.Equal(t, []int{5}, single.All())

	multi := orig.Pop(2)

	assert.Equal(t, []int{1, 2}, orig.All())
	assert.Equal(t, []int{3, 4}, multi.All())
}

func TestSplit(t *testing.T) {
	values := collection.From([]int{1, 2, 3, 4, 5, 6}).Split(3)

	assert.Equal(t, []int{1, 2, 3}, values[0].All())
	assert.Equal(t, []int{4, 5, 6}, values[1].All())
}

func TestDiff(t *testing.T) {
	first := collection.From([]int{1, 2, 3, 4, 5})
	diff := first.Diff(collection.From([]int{2, 5}))

	assert.Equal(t, []int{1, 3, 4}, diff.All())
}

func TestJoin(t *testing.T) {
	cs := collection.From([]string{"first", "second", "third"}).Join(join.CommaSeparatedJoin)
	assert.Equal(t, "first, second, third", cs)

	list := collection.From([]string{"first", "second", "third"}).Join(join.ListJoin)
	assert.Equal(t, "first, second and third", list)

	custom := collection.From([]string{"first", "second", "third"}).Join(join.Method{Between: "… ", Final: " & "})
	assert.Equal(t, "first… second & third", custom)
}

func TestFirstX(t *testing.T) {
	two := collection.From([]int{1, 2, 3, 4, 5}).FirstX(2)
	assert.Equal(t, collection.From([]int{1, 2}).All(), two.All())

	one := collection.From([]int{1}).FirstX(2)
	assert.Equal(t, 1, one.Count())
	assert.Equal(t, collection.From([]int{1}).All(), one.All())
}

func TestEmpty(t *testing.T) {
	truthy := collection.Make[string]().Empty()
	assert.Equal(t, true, truthy)

	falsy := collection.From([]int{1}).Empty()
	assert.Equal(t, false, falsy)
}

func TestNotEmpty(t *testing.T) {
	falsy := collection.Make[string]().NotEmpty()
	assert.Equal(t, false, falsy)

	truthy := collection.From([]int{1}).NotEmpty()
	assert.Equal(t, true, truthy)
}

func TestPrepend(t *testing.T) {
	col := collection.From([]int{2, 3, 4, 5}).Prepend(1)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, col.All())
}

func TestSet(t *testing.T) {
	col := collection.From([]int{1, 2, 5})
	col.Set(2, 3)
	assert.Equal(t, []int{1, 2, 3}, col.All())

	col.Set(4, 5)
	assert.Equal(t, []int{1, 2, 3, 0, 5}, col.All())
}

func TestSafeSet(t *testing.T) {
	col := collection.From([]int{1, 2, 5})
	err := col.SafeSet(2, 3)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, col.All())

	err = col.SafeSet(4, 5)
	assert.ErrorIs(t, err, collection.ErrIndexOutOfRange)
	assert.Equal(t, []int{1, 2, 3}, col.All())
}

func TestEach(t *testing.T) {
	col := collection.From([]int{1, 2, 3, 4, 5})
	incr := 0

	col.Each(func(i int, value int) {
		incr = incr + value
	})

	assert.Equal(t, 15, incr)
}

func TestEachCtx(t *testing.T) {
	col := collection.From([]int{1, 2, 3, 4, 5})
	incr := 0
	ctx, cancel := context.WithCancel(context.Background())

	col.EachCtx(ctx, func(i int, value int) {
		incr = incr + value
		if value == 3 {
			cancel()
		}
	})

	assert.Equal(t, 6, incr)
}

func TestEvery(t *testing.T) {
	truthy := collection.From([]int{1, 3, 5, 7, 9}).Every(func(i int, value int) bool {
		return value%2 == 1
	})
	assert.Equal(t, true, truthy)

	falsy := collection.From([]string{"dog", "cat", "lion"}).Every(func(i int, value string) bool {
		return len(value) == 3
	})
	assert.Equal(t, false, falsy)

	empty := collection.Make[string]().Every(func(i int, value string) bool {
		return false
	})
	assert.Equal(t, true, empty)
}

func TestRandom(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)

	col := collection.From([]int{1, 2, 3, 4, 5})

	assert.Equal(t, []int{2, 3}, col.Random(r, 2).All())
	assert.Equal(t, []int{3, 5}, col.Random(r, 2).All())
	assert.Equal(t, []int{2, 4, 1, 1, 2, 1, 5, 2, 3, 5}, col.Random(r, 10).All())
}

func TestReverse(t *testing.T) {
	col := collection.From([]int{1, 2, 3, 4, 5})
	assert.Equal(t, []int{5, 4, 3, 2, 1}, col.Reverse().All())
}

func TestSearch(t *testing.T) {
	res := collection.FromRange(1, 5).Search(func(i int, value int) bool {
		return value == 3
	})

	assert.Equal(t, 2, res)

	notFound := collection.FromRange(1, 5).Search(func(i int, value int) bool {
		return value == 12
	})

	assert.Equal(t, -1, notFound)
}

func TestSafeSearch(t *testing.T) {
	res, err := collection.FromRange(1, 5).SafeSearch(func(i int, value int) bool {
		return value == 3
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, res)

	notFound, err := collection.FromRange(1, 5).SafeSearch(func(i int, value int) bool {
		return value == 12
	})

	assert.ErrorIs(t, err, collection.ErrNoItem)
	assert.Equal(t, -1, notFound)
}
