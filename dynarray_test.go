package dynarray

import (
	"fmt"
	"iter"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew25(t *testing.T) {
	dyn := New[int](kSkipFactor * kSkipFactor)
	assert.Equal(t, 3, len(dyn.nodeLists))
	assert.Equal(t, kSkipFactor*kSkipFactor, dyn.Len())
}

func TestNew24(t *testing.T) {
	dyn := New[int](kSkipFactor*kSkipFactor - 1)
	assert.Equal(t, 2, len(dyn.nodeLists))
	assert.Equal(t, kSkipFactor*kSkipFactor-1, dyn.Len())
}

func TestNew1(t *testing.T) {
	dyn := New[int](1)
	assert.Equal(t, 1, len(dyn.nodeLists))
	assert.Equal(t, 1, dyn.Len())
}

func TestNew0(t *testing.T) {
	dyn := New[int](0)
	assert.Equal(t, 0, len(dyn.nodeLists))
	assert.Equal(t, 0, dyn.Len())
}

func TestNewZero(t *testing.T) {
	var dyn DynArray[int]
	assert.Equal(t, 0, len(dyn.nodeLists))
	assert.Equal(t, 0, dyn.Len())
}

func TestNewPanic(t *testing.T) {
	assert.Panics(t, func() { New[int](-1) })
	assert.Panics(t, func() { New[int](math.MaxInt) })
}

func TestLevelCountForNewValueNode(t *testing.T) {
	dyn := New[int](kSkipFactor * kSkipFactor)
	counts := make(map[int]int)
	for range 200 {
		levelCount := dyn.levelCountForNewValueNode()
		counts[levelCount] = counts[levelCount] + 1
	}
	assert.Equal(t, map[int]int{3: 8, 2: 32, 1: 160}, counts)
}

func TestSet(t *testing.T) {
	dyn := New[string](10)
	dyn.Set(6, "hello")

	assert.Panics(t, func() { dyn.Set(0, "goodbye") })
	dyn.Set(12, "goodbye")

	val, ok := dyn.Get(6)
	assert.Equal(t, "hello", val)
	assert.True(t, ok)

	val, ok = dyn.Get(12)
	assert.Equal(t, "goodbye", val)
	assert.True(t, ok)

	assert.Equal(t, 12, dyn.Len())
}

func TestSetExtensive(t *testing.T) {
	dyn := New[int](100)
	dyn.Set(9, 3)
	dyn.Set(36, 6)
	dyn.Set(81, 9)
	dyn.Set(4, 2)
	dyn.Set(25, 5)
	dyn.Set(64, 8)
	dyn.Set(1, 1)
	dyn.Set(16, 4)
	dyn.Set(49, 7)
	dyn.Set(100, 10)
	assert.Equal(t, 1, dyn.GetValue(1))
	assert.Equal(t, 2, dyn.GetValue(4))
	assert.Equal(t, 3, dyn.GetValue(9))
	assert.Equal(t, 4, dyn.GetValue(16))
	assert.Equal(t, 5, dyn.GetValue(25))
	assert.Equal(t, 6, dyn.GetValue(36))
	assert.Equal(t, 7, dyn.GetValue(49))
	assert.Equal(t, 8, dyn.GetValue(64))
	assert.Equal(t, 9, dyn.GetValue(81))
	assert.Equal(t, 0, dyn.GetValue(99))
	assert.Equal(t, 10, dyn.GetValue(100))
}

func TestSetPanics(t *testing.T) {
	dyn := New[int](math.MaxInt - 1)
	assert.Panics(t, func() { dyn.Set(math.MaxInt, 1) })
}

func TestSetByCursor(t *testing.T) {
	dyn := New[int](25)
	c := dyn.NewCursor()
	defer c.Close()
	assert.NoError(t, c.Advance(8))
	c.Set(3)
	assert.NoError(t, c.Advance(7))
	c.Set(4)
	assert.NoError(t, c.Advance(9))
	c.Set(5)
	assert.NoError(t, c.Advance(-24))
	c.Set(1)
	assert.NoError(t, c.Advance(3))
	c.Set(2)
	assert.NoError(t, c.Advance(21))
	assert.Equal(t, 25, c.Index())
	assert.Equal(t, 5, c.GetOrElse(1234))
	assert.NoError(t, c.Advance(-4))
	assert.Equal(t, 21, c.Index())
	assert.Equal(t, 1234, c.GetOrElse(1234))
	assert.NoError(t, c.Advance(-5))
	assert.Equal(t, 4, c.GetValue())
	assert.NoError(t, c.Advance(-8))
	assert.Equal(t, 8, c.Index())
	assert.Equal(t, 0, c.GetValue())
	assert.NoError(t, c.Advance(-4))
	assert.Equal(t, 2, c.GetValue())
	assert.NoError(t, c.Advance(-3))
	assert.Error(t, c.Advance(-1))
	assert.Equal(t, 1, c.GetValue())
	assert.Equal(t, 1, c.Index())
	c.Advance(26)
	assert.Equal(t, 27, c.Index())
	assert.Equal(t, 0, c.GetValue())
	assert.Empty(t, verify(c))
}

func TestGetPanics(t *testing.T) {
	dyn := New[int](100)
	assert.Equal(t, 4936, dyn.GetOrElse(102, 4936))
	assert.Panics(t, func() { dyn.Get(0) })
}

func TestGetSetClearAtEnd(t *testing.T) {
	dyn := New[int](100)
	_, ok := dyn.Get(101)
	assert.False(t, ok)
	_, ok = dyn.Get(104)
	assert.False(t, ok)
	dyn.Clear(101)
	dyn.Clear(104)
	dyn.Set(101, 7)
	assert.Equal(t, 101, dyn.Len())
	dyn.Set(104, 9)
	assert.Equal(t, 7, dyn.GetOrElse(101, 4936))
	assert.Equal(t, 104, dyn.Len())
	assert.Equal(t, 9, dyn.GetValue(104))
}

func TestAddAtEnd(t *testing.T) {
	var dyn DynArray[int]
	c := dyn.NewCursorAt(5)
	defer c.Close()
	_, ok := c.Get()
	assert.False(t, ok)
	value := 10
	for range 6 {
		c.Set(value)
		c.Advance(1)
		value += 10
	}
	assert.Empty(t, verify(c))
	assert.Equal(t, 10, dyn.GetValue(5))
	assert.Equal(t, 20, dyn.GetValue(6))
	assert.Equal(t, 30, dyn.GetValue(7))
	assert.Equal(t, 40, dyn.GetValue(8))
	assert.Equal(t, 50, dyn.GetValue(9))
	assert.Equal(t, 60, dyn.GetValue(10))
}

func TestClear(t *testing.T) {
	dyn := New[int](50)
	for i := range 50 {
		dyn.Set(i+1, 10*(i+1))
	}
	dyn.Clear(37)
	dyn.Clear(52)
	assert.Equal(t, 0, dyn.GetValue(37))
	assert.Equal(t, 360, dyn.GetValue(36))
	assert.Equal(t, 440, dyn.GetValue(44))
}

func TestClearWithCursor(t *testing.T) {
	dyn := New[int](50)
	c := dyn.NewCursor()
	defer c.Close()
	for i := range 50 {
		c.Set(10 * (i + 1))
		c.Advance(1)
	}
	assert.Equal(t, 51, c.Index())
	c.Goto(37)
	c.Clear()
	assert.Equal(t, 0, c.GetValue())
	assert.Equal(t, 37, c.Index())
	c.Goto(52)
	c.Clear()
	assert.Equal(t, 0, c.GetValue())
	assert.Equal(t, 52, c.Index())
	c.Goto(36)
	assert.Equal(t, 360, c.GetValue())
	assert.Equal(t, 36, c.Index())
	c.Goto(44)
	assert.Equal(t, 440, c.GetValue())
	assert.Equal(t, 44, c.Index())
	assert.Panics(t, func() { c.Goto(0) })
	assert.Equal(t, 44, c.Index())
	assert.Empty(t, verify(c))
}

func TestInsert(t *testing.T) {
	dyn := New[int](100)
	for i := range 50 {
		dyn.Set(2*i+2, 10*(2*i+2))
	}
	assert.NoError(t, dyn.Insert(50, 30))
	assert.Equal(t, 130, dyn.Len())
	for i := 1; i < 50; i++ {
		if i%2 != 0 {
			_, ok := dyn.Get(i)
			assert.False(t, ok)
		} else {
			val, ok := dyn.Get(i)
			assert.True(t, ok)
			assert.Equal(t, 10*i, val)
		}
	}
	for i := 50; i < 80; i++ {
		_, ok := dyn.Get(i)
		assert.False(t, ok)
	}
	for i := 80; i <= 130; i++ {
		if i%2 != 0 {
			_, ok := dyn.Get(i)
			assert.False(t, ok)
		} else {
			val, ok := dyn.Get(i)
			assert.True(t, ok)
			assert.Equal(t, 10*(i-30), val)
		}
	}
}

func TestInsert0(t *testing.T) {
	dyn := New[int](100)
	for i := range 50 {
		dyn.Set(2*i+2, 10*(2*i+2))
	}
	assert.NoError(t, dyn.Insert(50, 0))
	assert.Equal(t, 100, dyn.Len())
	for i := 1; i <= 100; i++ {
		if i%2 != 0 {
			_, ok := dyn.Get(i)
			assert.False(t, ok)
		} else {
			val, ok := dyn.Get(i)
			assert.True(t, ok)
			assert.Equal(t, 10*i, val)
		}
	}
}

func TestInsertError(t *testing.T) {
	dyn := New[int](100)
	assert.Error(t, dyn.Insert(0, 3))
	assert.Error(t, dyn.Insert(102, 3))
	assert.Error(t, dyn.Insert(100, -1))
	assert.Error(t, dyn.Insert(100, math.MaxInt-100))
	assert.Equal(t, 100, dyn.Len())
	assert.NoError(t, dyn.Insert(101, 3))
	assert.Equal(t, 103, dyn.Len())
}

func TestInsertWithCursor(t *testing.T) {
	dyn := New[int](100)
	for i := range 50 {
		dyn.Set(2*i+2, 10*(2*i+2))
	}
	c := dyn.NewCursorAt(100)
	defer c.Close()
	assert.NoError(t, c.Insert(2))
	assert.Equal(t, 102, c.Index())
	c.Set(1020)
	c.Advance(-2)
	c.Set(1000)
	assert.Empty(t, verify(c))
	assert.Equal(t, 102, c.Len())
	assert.Equal(t, 102, dyn.Len())
	for i := 1; i <= 102; i++ {
		if i%2 != 0 {
			_, ok := dyn.Get(i)
			assert.False(t, ok)
		} else {
			val, ok := dyn.Get(i)
			assert.True(t, ok)
			assert.Equal(t, 10*i, val)
		}
	}
}

func TestInsertWithCursor0(t *testing.T) {
	dyn := New[int](100)
	for i := range 50 {
		dyn.Set(2*i+2, 10*(2*i+2))
	}
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(90)
	assert.NoError(t, c.Insert(0))
	assert.Equal(t, 90, c.Index())
	assert.Equal(t, 900, c.GetValue())
	assert.Empty(t, verify(c))
	assert.Equal(t, 100, c.Len())
	assert.Equal(t, 100, dyn.Len())
	for i := 1; i <= 100; i++ {
		if i%2 != 0 {
			_, ok := dyn.Get(i)
			assert.False(t, ok)
		} else {
			val, ok := dyn.Get(i)
			assert.True(t, ok)
			assert.Equal(t, 10*i, val)
		}
	}
}

func TestInsertAtEmpty(t *testing.T) {
	dyn := New[int](100)
	for i := range 50 {
		dyn.Set(2*i+2, 10*(2*i+2))
	}
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(97)
	assert.NoError(t, c.Insert(2))
	assert.Equal(t, 99, c.Index())
	assert.Zero(t, c.GetValue())
	c.Advance(1)
	assert.Equal(t, 980, c.GetValue())
	c.Advance(-4)
	assert.Equal(t, 960, c.GetValue())
	assert.Empty(t, verify(c))
	assert.Equal(t, 102, c.Len())
}

func TestFirstValue(t *testing.T) {
	var dyn DynArray[int]
	_, ok := dyn.FirstValue()
	assert.False(t, ok)
}

func TestFirstValue2(t *testing.T) {
	dyn := New[int](100)
	_, ok := dyn.FirstValue()
	assert.False(t, ok)
}

func TestIterateValues(t *testing.T) {
	dyn := New[int](1000)
	for i := range 31 {
		dyn.Set((i+1)*(i+1), i+1)
	}
	c, ok := dyn.FirstValue()
	if ok {
		defer c.Close()
	}
	idx := 0
	for ; ok; ok = c.NextValue() {
		idx++
		assert.Equal(t, idx*idx, c.Index())
		assert.Equal(t, idx, c.GetValue())
	}
	assert.Equal(t, 31, idx)
	assert.Empty(t, verify(c))
	c.Goto(901)
	assert.True(t, c.NextValue())
	assert.Equal(t, 961, c.Index())
	assert.Equal(t, 31, c.GetValue())
	assert.Empty(t, verify(c))
}

func TestSnapshot(t *testing.T) {
	dyn := New[int](1000)
	for i := range 31 {
		dyn.Set((i+1)*(i+1), i+1)
	}
	s := dyn.MakeSnapshot()
	dyn.Insert(1, 100) // s remains unchanged even though we modified dyn.
	assert.Equal(t, 1, s.GetOrElse(1, 0))
	assert.Equal(t, 0, s.GetOrElse(2, 0))
	assert.Equal(t, 0, s.GetOrElse(3, 0))
	assert.Equal(t, 2, s.GetOrElse(4, 0))
	assert.Equal(t, 31, s.GetOrElse(961, 0))
	assert.Equal(t, 0, s.GetOrElse(962, 0))
	assert.Equal(t, 1000, s.Len())
	assert.Panics(t, func() { s.GetOrElse(0, 0) })
	assert.Equal(
		t,
		map[int]int{841: 29, 900: 30, 961: 31},
		iterToMap(s.AllFrom(840)))
	it := s.AllFrom(841)
	assert.Equal(
		t,
		map[int]int{841: 29, 900: 30, 961: 31},
		iterToMap(it))
	assert.Equal(
		t,
		map[int]int{841: 29, 900: 30},
		iterToMapTo(it, 961))
	assert.Equal(
		t,
		map[int]int{961: 31},
		iterToMap(s.AllFrom(961)))
	assert.Empty(
		t,
		iterToMap(s.AllFrom(962)))
	assert.Equal(t, 31, len(iterToMap(s.All())))
	assert.Equal(t, 31, len(iterToMap(s.AllFrom(0))))
}

func TestFromSnapshot(t *testing.T) {
	dyn := New[int](1000)
	for i := range 31 {
		dyn.Set((i+1)*(i+1), i+1)
	}
	dynCopy := FromSnapshot(dyn.MakeSnapshot())
	assert.Equal(t, toMap(dyn), toMap(dynCopy))
	assert.Equal(t, dyn.Len(), dynCopy.Len())
}

func TestZeroSnapshot(t *testing.T) {
	var s Snapshot[string]
	assert.Zero(t, s.Len())
	assert.Equal(t, "novalue", s.GetOrElse(1, "novalue"))
	assert.Empty(t, iterToMap(s.All()))
}

func TestDelete(t *testing.T) {
	dyn := New[int](1000)
	for i := range 31 {
		dyn.Set((i+1)*(i+1), i+1)
	}
	assert.NoError(t, dyn.Delete(2, 500))
	assert.Equal(t, 1, dyn.GetValue(1))
	sqrt := 23
	for i := 2; i <= 500; i++ {
		if i != sqrt*sqrt-500 {
			_, ok := dyn.Get(i)
			assert.False(t, ok)
		} else {
			val, ok := dyn.Get(i)
			assert.Equal(t, sqrt, val)
			assert.True(t, ok)
			sqrt++
		}
	}
	assert.Equal(t, 500, dyn.Len())
}

func TestDelete0(t *testing.T) {
	dyn := New[int](1000)
	for i := range 31 {
		dyn.Set((i+1)*(i+1), i+1)
	}
	assert.NoError(t, dyn.Delete(2, 0))
	sqrt := 1
	for i := 1; i <= 1000; i++ {
		if i != sqrt*sqrt {
			_, ok := dyn.Get(i)
			assert.False(t, ok)
		} else {
			val, ok := dyn.Get(i)
			assert.Equal(t, sqrt, val)
			assert.True(t, ok)
			sqrt++
		}
	}
	assert.Equal(t, 1000, dyn.Len())
}

func TestDeleteWithCursor(t *testing.T) {
	dyn := New[int](1000)
	for i := range 31 {
		dyn.Set((i+1)*(i+1), i+1)
	}
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(2)
	assert.NoError(t, c.Delete(500))
	assert.Equal(t, 2, c.Index())
	assert.Equal(t, 0, c.GetValue())
	sqrt := 23
	for c.Index() <= 500 {
		if c.Index() != sqrt*sqrt-500 {
			_, ok := c.Get()
			assert.False(t, ok)
		} else {
			val, ok := c.Get()
			assert.Equal(t, sqrt, val)
			assert.True(t, ok)
			sqrt++
		}
		c.Advance(1)
	}
	assert.Equal(t, 500, c.Len())
	assert.Empty(t, verify(c))
}

func TestDeleteAtValue(t *testing.T) {
	dyn := New[int](1000)
	for i := range 31 {
		dyn.Set((i+1)*(i+1), i+1)
	}
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(4)
	assert.NoError(t, c.Delete(525))
	assert.Equal(t, 4, c.Index())
	assert.Equal(t, 23, c.GetValue())
	sqrt := 23
	for c.Index() <= 475 {
		if c.Index() != sqrt*sqrt-525 {
			_, ok := c.Get()
			assert.False(t, ok)
		} else {
			val, ok := c.Get()
			assert.Equal(t, sqrt, val)
			assert.True(t, ok)
			sqrt++
		}
		c.Advance(1)
	}
	assert.Equal(t, 475, c.Len())
	assert.Empty(t, verify(c))
}

func TestDeleteWithCursor0(t *testing.T) {
	dyn := New[int](1000)
	for i := range 31 {
		dyn.Set((i+1)*(i+1), i+1)
	}
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(2)
	assert.NoError(t, c.Delete(0))
	assert.Equal(t, 2, c.Index())
	assert.Equal(t, 0, c.GetValue())
	sqrt := 2
	for c.Index() <= 1000 {
		if c.Index() != sqrt*sqrt {
			_, ok := c.Get()
			assert.False(t, ok)
		} else {
			val, ok := c.Get()
			assert.Equal(t, sqrt, val)
			assert.True(t, ok)
			sqrt++
		}
		c.Advance(1)
	}
	assert.Equal(t, 1000, c.Len())
	assert.Empty(t, verify(c))
}

func TestDeleteError(t *testing.T) {
	dyn := New[int](1000)
	assert.Error(t, dyn.Delete(1002, 0))
	assert.Error(t, dyn.Delete(0, 0))
	assert.Error(t, dyn.Delete(500, 502))
	assert.Error(t, dyn.Delete(500, -1))
	assert.Equal(t, 1000, dyn.Len())
	assert.NoError(t, dyn.Delete(1001, 0))
	assert.NoError(t, dyn.Delete(1000, 1))
	assert.Equal(t, 999, dyn.Len())
}

func TestInsertDeleteOnEmptyArray(t *testing.T) {
	var dyn DynArray[int]
	c := dyn.NewCursor()
	defer c.Close()
	assert.NoError(t, c.Insert(10))
	assert.NoError(t, c.Delete(0))
	assert.Equal(t, 11, c.Index())
	assert.Equal(t, 10, c.Len())
	c.Goto(1)
	assert.NoError(t, c.Delete(10))
	assert.Equal(t, 1, c.Index())
	assert.Equal(t, 0, c.Len())
	assert.Empty(t, verify(c))
}

func TestInsertDeleteOnEmptyArrayWithValues(t *testing.T) {
	var dyn DynArray[int]
	c := dyn.NewCursor()
	defer c.Close()
	assert.NoError(t, c.Insert(10))
	assert.NoError(t, c.Advance(-10))
	for range 10 {
		c.Set(1)
		c.Advance(1)
	}
	assert.Equal(t, 11, c.Index())
	assert.Empty(t, verify(c))
	c.Goto(1)
	for range 10 {
		assert.Equal(t, 1, c.GetValue())
		c.Advance(1)
	}
	c.Goto(1)
	assert.NoError(t, c.Delete(10))
	assert.Empty(t, verify(c))
	_, ok := dyn.FirstValue()
	assert.False(t, ok)
	assert.Zero(t, dyn.Len())
}

func TestMultiClear(t *testing.T) {
	dyn := New[int](100)
	for i := 1; i <= 100; i++ {
		dyn.Set(i, 1)
	}
	c := dyn.NewCursor()
	defer c.Close()
	for c.Index() <= c.Len() {
		c.Clear()
		c.Advance(1)
	}
	assert.Empty(t, verify(c))
	_, ok := dyn.FirstValue()
	assert.False(t, ok)
}

func TestMultiCursorsSet(t *testing.T) {
	dyn := New[int](10)
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(8)
	d := dyn.NewCursor()
	defer d.Close()
	d.Goto(8)
	c.Set(1)
	assert.Equal(t, 1, d.GetValue())
	assert.Empty(t, verify(c))
	assert.Empty(t, verify(d))
}

func TestMultiCursorsClear(t *testing.T) {
	dyn := New[int](10)
	dyn.Set(8, 1)
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(8)
	d := dyn.NewCursor()
	defer d.Close()
	d.Goto(8)
	assert.Equal(t, 1, d.GetValue())
	c.Clear()
	assert.Equal(t, 0, d.GetValue())
	assert.Empty(t, verify(c))
	assert.Empty(t, verify(d))
}

func TestMultiCursorInsert(t *testing.T) {
	dyn := New[int](100)
	for i := 1; i <= 100; i++ {
		dyn.Set(i, 10*i)
	}
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(60)
	d := dyn.NewCursor()
	defer d.Close()
	d.Goto(60)
	e := dyn.NewCursor()
	defer e.Close()
	e.Goto(61)
	f := dyn.NewCursor()
	defer f.Close()
	f.Goto(59)
	assert.NoError(t, c.Insert(20))
	assert.Equal(t, 80, c.Index())
	assert.Equal(t, 600, c.GetValue())
	assert.Equal(t, 80, d.Index())
	assert.Equal(t, 600, d.GetValue())
	assert.Equal(t, 81, e.Index())
	assert.Equal(t, 610, e.GetValue())
	assert.Equal(t, 59, f.Index())
	assert.Equal(t, 590, f.GetValue())
	assert.Empty(t, verify(c))
	assert.Empty(t, verify(d))
	assert.Empty(t, verify(e))
	assert.Empty(t, verify(f))
}

func TestMultiCursorDelete(t *testing.T) {
	dyn := New[int](100)
	for i := 1; i <= 100; i++ {
		dyn.Set(i, 10*i)
	}
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(60)
	d := dyn.NewCursor()
	defer d.Close()
	d.Goto(60)
	e := dyn.NewCursor()
	defer e.Close()
	e.Goto(80)
	f := dyn.NewCursor()
	defer f.Close()
	f.Goto(59)
	assert.NoError(t, c.Delete(20))
	assert.Equal(t, 60, c.Index())
	assert.Equal(t, 800, c.GetValue())
	assert.True(t, d.IsClosed())
	assert.Equal(t, 60, e.Index())
	assert.Equal(t, 800, e.GetValue())
	assert.Equal(t, 59, f.Index())
	assert.Equal(t, 590, f.GetValue())
	assert.Empty(t, verify(c))
	assert.Empty(t, verify(e))
	assert.Empty(t, verify(f))
}

func TestCloseCursor(t *testing.T) {
	dyn := New[int](100)
	c := dyn.NewCursor()
	c.Set(3)
	assert.False(t, c.IsClosed())
	c.Close()
	c.Close() // Idempotent
	assert.True(t, c.IsClosed())
	assert.Panics(t, func() { c.Advance(1) })
	assert.Panics(t, func() { c.Clear() })
	assert.Panics(t, func() { c.Goto(1) })
	assert.Panics(t, func() { c.Delete(1) })
	assert.Panics(t, func() { c.Get() })
	assert.Panics(t, func() { c.GetValue() })
	assert.Panics(t, func() { c.Index() })
	assert.Panics(t, func() { c.Insert(1) })
	assert.Panics(t, func() { c.Len() })
	assert.Panics(t, func() { c.NextValue() })
	assert.Panics(t, func() { c.Set(1) })
}

func TestSparseDelete(t *testing.T) {
	dyn := New[int](1000)
	dyn.Set(343, 1)
	dyn.Set(344, 1)
	dyn.Set(686, 1)
	dyn.Set(687, 1)
	dyn.Set(924, 1)
	assert.Equal(
		t, map[int]int{343: 1, 344: 1, 686: 1, 687: 1, 924: 1}, toMap(dyn))
	c := dyn.NewCursor()
	defer c.Close()
	c.Goto(343)
	c.Delete(2)
	assert.Equal(t, map[int]int{684: 1, 685: 1, 922: 1}, toMap(dyn))
	c.Goto(682)
	c.Delete(4)
	assert.Equal(t, map[int]int{918: 1}, toMap(dyn))
	c.Goto(916)
	c.Delete(3)
	assert.Empty(t, toMap(dyn))
	c.Goto(987)
	c.Delete(5)
	assert.Empty(t, toMap(dyn))
	assert.Equal(t, 986, dyn.Len())
}

func TestSparseInsert(t *testing.T) {
	dyn := New[int](1000)
	assert.Empty(t, toMap(dyn))
	c := dyn.NewCursor()
	defer c.Close()
	c.Insert(2)
	c.Advance(-2)
	c.Set(2)
	c.Advance(1)
	c.Set(3)
	assert.Equal(t, map[int]int{1: 2, 2: 3}, toMap(dyn))
	c.Goto(587)
	c.Set(5)
	c.Insert(1)
	assert.Equal(t, map[int]int{1: 2, 2: 3, 588: 5}, toMap(dyn))
	assert.Equal(t, 1003, dyn.Len())
	c.Goto(1003)
	c.Insert(5)
	c.Set(7)
	assert.Equal(t, map[int]int{1: 2, 2: 3, 588: 5, 1008: 7}, toMap(dyn))
	assert.Equal(t, 1008, dyn.Len())
	c.Goto(1)
	c.Clear()
	c.Advance(1)
	c.Clear()
	c.Goto(588)
	c.Clear()
	c.Goto(1008)
	c.Clear()
	assert.Empty(t, toMap(dyn))
	assert.Equal(t, 1008, dyn.Len())
}

func TestZeroArrayClear(t *testing.T) {
	var dyn DynArray[int]
	dyn.Clear(1)
}

func TestZeroArrayDelete(t *testing.T) {
	var dyn DynArray[int]
	assert.NoError(t, dyn.Delete(1, 0))
}

func TestZeroArrayFirstValue(t *testing.T) {
	var dyn DynArray[int]
	_, ok := dyn.FirstValue()
	assert.False(t, ok)
}

func TestZeroArrayGet(t *testing.T) {
	var dyn DynArray[int]
	assert.Zero(t, dyn.GetValue(1))
}

func TestZeroArrayInsert(t *testing.T) {
	var dyn DynArray[int]
	assert.NoError(t, dyn.Insert(1, 3))
	assert.Zero(t, dyn.GetValue(2))
	assert.Equal(t, 3, dyn.Len())
}

func TestZeroArraySet(t *testing.T) {
	var dyn DynArray[int]
	dyn.Set(2, 3)
	assert.Equal(t, 3, dyn.GetValue(2))
}

func TestTruncate(t *testing.T) {
	var dyn DynArray[int]
	for i := 1; i <= 6; i++ {
		dyn.Set(i, 3*i)
	}
	assert.NoError(t, dyn.Truncate(3))
	assert.Equal(t, map[int]int{1: 3, 2: 6, 3: 9}, toMap(&dyn))
	assert.Equal(t, 3, dyn.Len())
	assert.NoError(t, dyn.Truncate(3))
	assert.Equal(t, map[int]int{1: 3, 2: 6, 3: 9}, toMap(&dyn))
	assert.Equal(t, 3, dyn.Len())
	assert.NoError(t, dyn.Truncate(6))
	assert.Equal(t, map[int]int{1: 3, 2: 6, 3: 9}, toMap(&dyn))
	assert.Equal(t, 6, dyn.Len())
}

func TestTruncateError(t *testing.T) {
	var dyn DynArray[int]
	assert.NoError(t, dyn.Truncate(0))
	assert.Error(t, dyn.Truncate(-1))
	assert.Error(t, dyn.Truncate(math.MaxInt))
	assert.NoError(t, dyn.Truncate(math.MaxInt-1))
}

func TestRandomInitRetrieveByCursor(t *testing.T) {
	dyn := New[int](1000)

	// Insert in random places (every 377th place).
	for i := 1; i <= dyn.Len(); i++ {
		dyn.Set((377*i-1)%dyn.Len()+1, i)
	}
	c := dyn.NewCursor()
	defer c.Close()
	for c.Index() <= c.Len() {

		// 313 is the multiplicate inverse of 377 mod 1000
		assert.Equal(t, (c.Index()*313-1)%dyn.Len()+1, c.GetValue())
		c.Advance(1)
	}
}

func TestRandomInitClearByCursor(t *testing.T) {
	dyn := New[int](1000)

	// Insert in random places (every 377th place).
	for i := 1; i <= dyn.Len(); i++ {
		dyn.Set((377*i-1)%dyn.Len()+1, i)
	}

	// Clear everything by cursor
	c, ok := dyn.FirstValue()
	if ok {
		defer c.Close()
	}
	for ; ok; ok = c.NextValue() {
		c.Clear()
	}
	assert.Empty(t, verify(c))

	// Verify no values in DynArray
	_, ok = dyn.FirstValue()
	assert.False(t, ok)
}

func TestRandomInitRetrieveByIndex(t *testing.T) {
	dyn := New[int](1000)

	// Insert in random places (every 377th place).
	for i := 1; i <= dyn.Len(); i++ {
		dyn.Set((377*i-1)%dyn.Len()+1, i)
	}
	for i := 1; i <= dyn.Len(); i++ {

		// 313 is the multiplicate inverse of 377 mod 1000
		assert.Equal(t, (i*313-1)%dyn.Len()+1, dyn.GetValue(i))
	}
}

func TestRandomInitClearByIndex(t *testing.T) {
	dyn := New[int](1000)

	// Insert in random places (every 377th place).
	for i := 1; i <= dyn.Len(); i++ {
		dyn.Set((377*i-1)%dyn.Len()+1, i)
	}
	for i := 1; i <= dyn.Len(); i++ {
		dyn.Clear((549*i-1)%dyn.Len() + 1)
	}

	// Verify no values in DynArray
	_, ok := dyn.FirstValue()
	assert.False(t, ok)
}

func TestRandomInitClearByDelete(t *testing.T) {
	dyn := New[int](1000)

	// Insert in random places (every 377th place).
	for i := 1; i <= dyn.Len(); i++ {
		dyn.Set((377*i-1)%dyn.Len()+1, i)
	}
	dyn.Delete(1, 1000)

	assert.Zero(t, dyn.Len())

	// Verify no values in DynArray
	_, ok := dyn.FirstValue()
	assert.False(t, ok)
}

func TestNearMaxInt(t *testing.T) {
	var dyn DynArray[int]
	dyn.Set(math.MaxInt-3, 3)
	dyn.Set(math.MaxInt-2, 2)
	dyn.Set(math.MaxInt-1, 1)
	assert.Error(t, dyn.Insert(1, 1))
	assert.Equal(
		t,
		map[int]int{math.MaxInt - 3: 3, math.MaxInt - 2: 2, math.MaxInt - 1: 1},
		toMap(&dyn))
	dyn.Clear(math.MaxInt - 3)
	dyn.Clear(math.MaxInt - 2)
	dyn.Clear(math.MaxInt - 1)
	assert.Empty(t, toMap(&dyn))
	assert.NoError(t, dyn.Delete(1, 3))
	assert.Equal(t, math.MaxInt-4, dyn.Len())
}

func TestAllAndAllFrom(t *testing.T) {
	var dyn DynArray[string]
	iter := dyn.All()
	assert.Empty(t, iterToMap(iter))
	dyn.Set(2, "do")
	dyn.Set(5, "re")
	dyn.Set(10, "mi")
	assert.Equal(
		t,
		map[int]string{2: "do", 5: "re", 10: "mi"},
		iterToMap(iter))
	dyn.Set(1, "fa")
	assert.Equal(
		t,
		map[int]string{1: "fa", 2: "do", 5: "re", 10: "mi"},
		iterToMap(iter))
	assert.Equal(
		t,
		map[int]string{1: "fa", 2: "do"},
		iterToMapTo(iter, 3))
	assert.Equal(
		t,
		map[int]string{5: "re", 10: "mi"},
		iterToMap(dyn.AllFrom(5)))
	assert.Equal(
		t,
		map[int]string{10: "mi"},
		iterToMap(dyn.AllFrom(10)))
	assert.Empty(t, iterToMap(dyn.AllFrom(12)))
	assert.Equal(
		t,
		map[int]string{1: "fa", 2: "do", 5: "re", 10: "mi"},
		iterToMap(dyn.AllFrom(0)))
}

func verify[T any](c *Cursor[T]) string {
	c.checkBasic()
	if result := verifyCursorPiecesLen(c); result != "" {
		return result
	}
	if result := verifyIndexesSoFar(c); result != "" {
		return result
	}
	return verifyCursorPieces(c)
}

func verifyCursorPiecesLen[T any](c *Cursor[T]) string {
	cpLen := len(c.cursorPieces)
	expected := maxLevels(c.array.length)
	if cpLen != expected {
		return fmt.Sprintf("cpLen = %d expected %d", cpLen, expected)
	}
	if cpLen != len(c.array.nodeLists) {
		return fmt.Sprintf("cpLen = %d should match dynarray", cpLen)
	}
	return ""
}

func verifyIndexesSoFar[T any](c *Cursor[T]) string {
	for idx, cp := range c.cursorPieces {
		if cp.indexSoFar > c.index {
			return "indexSoFar greater than index"
		}
		if cp.ptr.next != nil && cp.nextIndex() <= c.index {
			return "next indexSoFar not greater than index"
		}
		if idx > 0 && cp.indexSoFar > c.cursorPieces[idx-1].indexSoFar {
			return "indexSoFar should decrease as level increases"
		}
	}
	return ""
}

func verifyCursorPieces[T any](c *Cursor[T]) string {
	for idx := range c.cursorPieces {
		if result := verifyCursorPiece(c, idx); result != "" {
			return result
		}
	}
	return ""
}

func verifyCursorPiece[T any](c *Cursor[T], idx int) string {
	ptr := c.array.nodeLists[idx]
	indexSoFar := 0
	for ptr != c.cursorPieces[idx].ptr {
		ptr = ptr.next
		indexSoFar += ptr.indexIncr
	}
	if c.cursorPieces[idx].indexSoFar != indexSoFar {
		return fmt.Sprintf("indexSoFar incorrect on index %d. expected %d got %d", idx, c.cursorPieces[idx].indexSoFar, indexSoFar)
	}
	return ""
}

func toMap[T any](d *DynArray[T]) map[int]T {
	return iterToMap(d.All())
}

func iterToMap[T any](it iter.Seq2[int, T]) map[int]T {
	result := make(map[int]T)
	for index, value := range it {
		result[index] = value
	}
	return result
}

func iterToMapTo[T any](it iter.Seq2[int, T], to int) map[int]T {
	result := make(map[int]T)
	for index, value := range it {
		if index >= to {
			break
		}
		result[index] = value
	}
	return result
}
