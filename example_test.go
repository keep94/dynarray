package dynarray_test

import (
	"fmt"

	"github.com/keep94/dynarray"
)

func ExampleDynArray_Sort() {
	var dyn dynarray.DynArray[string]
	dyn.Set(1, "do")
	dyn.Set(2, "re")
	dyn.Set(3, "mi")
	dyn.Set(5, "fa")
	dyn.Set(6, "so")
	dyn.Set(7, "la")
	dyn.Set(8, "ti")

	dyn.Sort(1, dyn.Len()+1, func(x, y string) bool { return x < y })

	for index, value := range dyn.All() {
		fmt.Println(index, value)
	}
	// Output:
	// 1 do
	// 2 fa
	// 3 la
	// 5 mi
	// 6 re
	// 7 so
	// 8 ti
}

func ExampleDynArray_GetOrElse() {
	var dyn dynarray.DynArray[string]
	dyn.Set(1, "do")
	dyn.Set(4, "re")
	dyn.Set(9, "mi")
	fmt.Println(dyn.Len())
	fmt.Println(dyn.GetOrElse(1, "novalue"))
	fmt.Println(dyn.GetOrElse(4, "novalue"))
	fmt.Println(dyn.GetOrElse(9, "novalue"))
	fmt.Println(dyn.GetOrElse(8, "novalue"))
	fmt.Println(dyn.GetOrElse(101, "novalue"))
	// Output:
	// 9
	// do
	// re
	// mi
	// novalue
	// novalue
}

func ExampleDynArray_Insert() {
	var dyn dynarray.DynArray[string]
	dyn.Set(1, "do")
	dyn.Set(4, "re")
	dyn.Set(9, "mi")
	dyn.Insert(3, 2)
	fmt.Println(dyn.Len())
	fmt.Println(dyn.GetOrElse(1, "novalue"))
	fmt.Println(dyn.GetOrElse(6, "novalue"))
	fmt.Println(dyn.GetOrElse(11, "novalue"))
	// Output:
	// 11
	// do
	// re
	// mi
}

func ExampleDynArray_Delete() {
	var dyn dynarray.DynArray[string]
	dyn.Set(1, "do")
	dyn.Set(4, "re")
	dyn.Set(9, "mi")
	dyn.Delete(3, 2)
	fmt.Println(dyn.Len())
	for i := 1; i <= dyn.Len(); i++ {
		fmt.Println(dyn.GetOrElse(i, "novalue"))
	}
	// Output:
	// 7
	// do
	// novalue
	// novalue
	// novalue
	// novalue
	// novalue
	// mi
}

func ExampleDynArray_All() {
	var dyn dynarray.DynArray[string]
	dyn.Set(2, "do")
	dyn.Set(5, "re")
	dyn.Set(10, "mi")
	for index, value := range dyn.All() {
		fmt.Println(index, value)
	}
	// Output:
	// 2 do
	// 5 re
	// 10 mi
}

func ExampleDynArray_Clear() {
	var dyn dynarray.DynArray[string]
	dyn.Set(7, "do")
	dyn.Clear(7)
	fmt.Println(dyn.Len())
	fmt.Println(dyn.GetOrElse(7, "novalue"))
	// Output:
	// 7
	// novalue
}

func ExampleCursor() {
	var dyn dynarray.DynArray[string]
	c := dyn.NewCursor()
	defer c.Close()
	c.Set("the")
	c.Advance(1)
	c.Set("quick")
	c.Advance(1)
	c.Set("brown")
	c.Advance(1)
	c.Set("fox")
	c.Advance(1)
	c.Set("jumps")
	c.Advance(1)
	c.Set("over")
	c.Advance(1)
	c.Set("the")
	c.Advance(1)
	c.Set("lazy")
	c.Advance(1)
	c.Set("dog")

	c.Goto(1)

	for c.GetValue() != "fox" {
		c.Advance(1)
	}
	c.Insert(1)
	c.Advance(-1)
	c.Set("inserted_value")
	for i := 1; i <= dyn.Len(); i++ {
		fmt.Println(dyn.GetValue(i))
	}
	// Output:
	// the
	// quick
	// brown
	// inserted_value
	// fox
	// jumps
	// over
	// the
	// lazy
	// dog
}
