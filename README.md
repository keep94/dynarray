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

Documentation and examples can be found [here](https://pkg.go.dev/github.com/keep94/dynarray).
