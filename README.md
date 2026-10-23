dynarray
========

Package dynarray contains data structures for Dynamic Arrays.

## About

DynArray works like regular arrays except for the following differences:

1. Indexing is one based.
1. An index can hold no value at all.
1. Memory usage depends on count of values stored not the index count.
1. When storing values beyond the highest index, length grows automatically instead of seeing index out of bounds.
1. Inserting into or deleting from the beginning or middle of the array takes O(log n) time instead of O(n) time.
1. Accessing a random array element takes O(log n) time instead of O(1) time.

Cursor instances can speed up array access if reading and writing to the same area of the DynArray. DynArray instances are not safe to use with multiple goroutines.

## History

In high school, I took my first real computer science course where I learned about big-O notation. I learned that arrays offer fast, constant time random access, but inserting or deleting elements in the middle of an array is slow -- O(n) time. I learned that inserting or deleting elements in the middle of a linked list is fast taking constant time, but accessing a random element in a linked list is slow taking O(n) time.

I wanted to invent a data structure that is fairly efficient at both random access as well as inserting and deleting. During the summer of 1989, I designed a data structure that offers O(log n) time for both random access and inserting and deleting. This data structure is similar to a B+ tree. However, I never got around to implementing it until now. I call it the dynamic array.

While the dynamic array offers O(log n) time for both random access and inserting and deleting, the overhead of the complex data structure slows it down overall. Therefore, it should only be used on very large datasets that require lots of inserting and deleting in the middle.

## Examples

```golang
package main

import (
    "fmt"

    "github.com/keep94/dynarray"
)

func main() {
    var dyn dynarray.DynArray[string]
    dyn.Set(1, "hello")
    dyn.Set(3, "goodbye")
    for i := 1; i <= dyn.Len(); i++ {
        fmt.Println(dyn.GetOrElse(i, "novalue"))
    }
    // Prints:
    //
    // hello
    // novalue
    // goodbye
}
```

More documentation and examples can be found [here](https://pkg.go.dev/github.com/keep94/dynarray).
