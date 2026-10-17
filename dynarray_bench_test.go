package dynarray

import (
	"testing"
)

func BenchmarkAdvanceForward100(b *testing.B) {
	dyn := New[int](100)
	initDynArray(dyn)
	for b.Loop() {
		advanceForward(dyn)
	}
}

func BenchmarkAdvanceForward1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArray(dyn)
	for b.Loop() {
		advanceForward(dyn)
	}
}

func BenchmarkAdvanceBackward100(b *testing.B) {
	dyn := New[int](100)
	initDynArray(dyn)
	for b.Loop() {
		advanceBackward(dyn)
	}
}

func BenchmarkAdvanceBackward1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArray(dyn)
	for b.Loop() {
		advanceBackward(dyn)
	}
}

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

func BenchmarkIterateSnapshot100(b *testing.B) {
	dyn := New[int](100)
	initDynArray(dyn)
	s := dyn.MakeSnapshot()
	for b.Loop() {
		iterateSnapshot(s)
	}
}

func BenchmarkIterateSnapshot1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArray(dyn)
	s := dyn.MakeSnapshot()
	for b.Loop() {
		iterateSnapshot(s)
	}
}

func BenchmarkIterateSnapshotByIndex100(b *testing.B) {
	dyn := New[int](100)
	initDynArray(dyn)
	s := dyn.MakeSnapshot()
	for b.Loop() {
		iterateSnapshotByIndex(s)
	}
}

func BenchmarkIterateSnapshotByIndex1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArray(dyn)
	s := dyn.MakeSnapshot()
	for b.Loop() {
		iterateSnapshotByIndex(s)
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

func BenchmarkSort100(b *testing.B) {
	dyn := New[int](100)
	initDynArrayRandom(dyn, 77)
	for b.Loop() {
		sortDyn(dyn)
	}
}

func BenchmarkSort1000(b *testing.B) {
	dyn := New[int](1000)
	initDynArrayRandom(dyn, 377)
	for b.Loop() {
		sortDyn(dyn)
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
		dyn.Set((skip*i-1)%dyn.Len()+1, i)
	}
}

func advanceForward(dyn *DynArray[int]) {
	c := dyn.NewCursor()
	defer c.Close()
	for c.Index() <= c.Len() {
		c.Get()
		c.Advance(1)
	}
}

func advanceBackward(dyn *DynArray[int]) {
	c := dyn.NewCursorAt(dyn.Len() + 1)
	defer c.Close()
	for c.Index() > 1 {
		c.Advance(-1)
		c.Get()
	}
}

func iterateByCursor(dyn *DynArray[int]) {
	for range dyn.All() {
	}
}

func iterateByIndex(dyn *DynArray[int]) {
	for i := 1; i <= dyn.Len(); i++ {
		dyn.GetValue(i)
	}
}

func iterateSnapshot(s Snapshot[int]) {
	for range s.All() {
	}
}

func iterateSnapshotByIndex(s Snapshot[int]) {
	for i := 1; i <= s.Len(); i++ {
		s.GetOrElse(i, 0)
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

func sortDyn(dyn *DynArray[int]) {
	dyn.Sort(1, dyn.Len()+1, func(x, y int) bool { return x < y })
}
