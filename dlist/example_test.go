// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Copyright 2017 - ionous. Modified to create intrusive linked lists.

package dlist_test

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

func ExampleCustomElement() {
	// Create a custom list element
	type MyElement struct {
		dlist.Hook     // anonymous support for the intrusive list; alternatively, users can implement the dlist.Intrusive interface for more control.
		MyData     int // some example data.
	}
	l := dlist.New()
	l.PushBack(&MyElement{MyData: 17})
	l.PushFront(&MyElement{MyData: 20})
	//
	for e := l.Front(); e != nil; e = dlist.Next(e) {
		fmt.Println(e.(*MyElement).MyData)
	}
	// Output:
	// 20
	// 17
}
