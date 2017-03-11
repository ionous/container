container provides intrusive linked lists for go. container/inlist modifies golang's container/list to provide an intrusive version of that package.


# Examples 

This first example mimics the container/list [example](https://golang.org/pkg/container/list/#example_) as closely as possible.
It creates a linked list of numbers, and prints them.

```go
import (
	"fmt"
	"github.com/ionous/container/inlist"
)

func Example() {
	// Create a new list and put some numbers in it.
	l := inlist.New()
	e4 := l.PushBack(inlist.NewElement(4))
	e1 := l.PushFront(inlist.NewElement(1))
	l.InsertBefore(inlist.NewElement(3), e4)
	l.InsertAfter(inlist.NewElement(2), e1)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = inlist.Next(e) {
		fmt.Println(inlist.Value(e))
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
	// Create a custom list element using inlist.Hook.
 	// Hook provides built-in support for the intrusive list:
	// you simply include it as an anonymous member of a struct. 
 	// ( Alternatively, you can implement the inlist.Intrusive interface for more control. )
	type MyElement struct {
		inlist.Hook     
		MyData     int // some example data.
	}
  
	l := inlist.New()
	l.PushBack(&MyElement{MyData: 17})
	l.PushFront(&MyElement{MyData: 20})
  
	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = inlist.Next(e) {
		fmt.Println(e.(*MyElement).MyData)
	}
  
	// Output:
	// 20
	// 17
}
```
