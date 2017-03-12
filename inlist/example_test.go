// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Copyright 2017 - ionous. Modified to create intrusive linked lists.

package inlist_test

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

func ExampleHook() {
	// Create a custom list element
	type MyElement struct {
		inlist.Hook     // anonymous support for the intrusive list; alternatively, users can implement the inlist.Intrusive interface for more control.
		MyData      int // some example data.
	}
	l := inlist.New()
	l.PushBack(&MyElement{MyData: 17})
	l.PushFront(&MyElement{MyData: 20})
	//
	for e := l.Front(); e != nil; e = inlist.Next(e) {
		fmt.Println(e.(*MyElement).MyData)
	}
	// Output:
	// 20
	// 17
}
