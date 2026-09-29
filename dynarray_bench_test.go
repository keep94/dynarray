package dynarray

import (
	"testing"
)

func BenchmarkIterateByCursor100(b *testing.B) {
	dyn := New[int](100)
	initDynArray(dyn)
	for b.Loop() {
		iterateByCursor(dyn)
	}
}

func BenchmarkIterateByCursor1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArray(dyn)
	for b.Loop() {
		iterateByCursor(dyn)
	}
}

func BenchmarkIterateByIndex100(b *testing.B) {
	dyn := New[int](100)
	initDynArray(dyn)
	for b.Loop() {
		iterateByIndex(dyn)
	}
}

func BenchmarkIterateByIndex1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArray(dyn)
	for b.Loop() {
		iterateByIndex(dyn)
	}
}

func BenchmarkIterateByIndexRandomInit1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArrayRandom(dyn, 377)
	for b.Loop() {
		iterateByIndex(dyn)
	}
}

func BenchmarkAddToEnd100(b *testing.B) {
	for b.Loop() {
		addToEnd(100)
	}
}

func BenchmarkAddToEnd1000(b *testing.B) {
	for b.Loop() {
		addToEnd(1000)
	}
}

func BenchmarkInsertDeleteOneIn100(b *testing.B) {
	dyn := New[int](100)
	initDynArray(dyn)
	for b.Loop() {
		insertDeleteOneIn(dyn, 63)
	}
}

func BenchmarkInsertDeleteOneIn1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArray(dyn)
	for b.Loop() {
		insertDeleteOneIn(dyn, 631)
	}
}

func BenchmarkInsertDeleteOneIn10000(b *testing.B) {
	dyn := New[int](10000)
	initDynArray(dyn)
	for b.Loop() {
		insertDeleteOneIn(dyn, 6317)
	}
}

func BenchmarkCreateCursor(b *testing.B) {
	dyn := New[int](1000)
	for b.Loop() {
		createCursor(dyn)
	}
}

func initDynArray(dyn *DynArray[int]) {
	for i := 1; i <= dyn.Len(); i++ {
		dyn.Set(i, 1)
	}
}

func initDynArrayRandom(dyn *DynArray[int], skip int) {
	for i := 1; i <= dyn.Len(); i++ {
		dyn.Set((skip*i-1)%dyn.Len()+1, 1)
	}
}

func iterateByCursor(dyn *DynArray[int]) {
	c, ok := dyn.FirstValue()
	if ok {
		defer c.Close()
	}
	for ; ok; ok = c.NextValue() {
		c.GetValue()
	}
}

func iterateByIndex(dyn *DynArray[int]) {
	for i := 1; i <= dyn.Len(); i++ {
		dyn.GetValue(i)
	}
}

func addToEnd(count int) {
	var dyn DynArray[int]
	c := dyn.NewCursor()
	defer c.Close()
	for range count {
		c.Set(1)
		c.Advance(1)
	}
}

func insertDeleteOneIn(dyn *DynArray[int], index int) {
	c := dyn.NewCursorAt(index)
	defer c.Close()
	c.Insert(1)
	c.Advance(-1)
	c.Delete(1)
}

func createCursor(dyn *DynArray[int]) {
	c := dyn.NewCursor()
	defer c.Close()
}
