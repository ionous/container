inlist provides intrusive linked lists for go. inlist/dlist modifies golang's container/list to provide an intrusive version of that package.


# Examples 

This first example mimics the container/list [example](https://golang.org/pkg/container/list/#example_) as closely as possible.
It creates a linked list of numbers, and prints them.

```go
import (
	"fmt"
	"github.com/ionous/inlist/dlist"
)

func Example() {
	// Create a new list and put some numbers in it.
	l := dlist.New()
	e4 := l.PushBack(dlist.NewElement(4))
	e1 := l.PushFront(dlist.NewElement(1))
	l.InsertBefore(dlist.NewElement(3), e4)
	l.InsertAfter(dlist.NewElement(2), e1)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = dlist.Next(e) {
		fmt.Println(dlist.Value(e))
	}
  
	// Output:
	// 1
	// 2
	// 3
	// 4
}
```

This second example creates custom list elements. It uses the package's "Hook" to change a simple user struct into a linked list element.

```go
func ExampleCustomElement() {
	// Create a custom list element using dlist.Hook.
  // Hook provides built-in support for the intrusive list:
  // you simply include it as an anonymous member of a struct. 
  // ( Alternatively, you can implement the dlist.Intrusive interface for more control. )
	type MyElement struct {
		dlist.Hook     
		MyData     int // some example data.
	}
  
	l := dlist.New()
	l.PushBack(&MyElement{MyData: 17})
	l.PushFront(&MyElement{MyData: 20})
  
	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = dlist.Next(e) {
		fmt.Println(e.(*MyElement).MyData)
	}
  
	// Output:
	// 20
	// 17
}
```
