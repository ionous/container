// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Copyright 2017 - ionous. Modified to create intrusive linked lists.

package dlist

import "testing"

func checkListLen(t *testing.T, l *List, len int) (okay bool) {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
	} else {
		okay = true
	}
	return
}

func checkListPointers(t *testing.T, l *List, es []Intrusive) {
	root := &l.root

	if !checkListLen(t, l, len(es)) {
		return
	}

	// zero length lists must be the zero value or properly initialized (sentinel circle)
	if len(es) == 0 {
		if l.root.next != nil && l.root.next != root || l.root.prev != nil && l.root.prev != root {
			t.Errorf("l.root.next = %p, l.root.prev = %p; both should both be nil or %p", &l.root.next, &l.root.prev, root)
		}
		return
	}
	// len(es) > 0

	// check internal and external prev/next connections
	for i, e := range es {
		prev := Intrusive(root)
		var PrevElement Intrusive
		if i > 0 {
			prev = es[i-1]
			PrevElement = prev
		}
		if p := e.Predecessor(); p != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p", i, &e, &p, &prev)
		}
		if p := Prev(e); p != PrevElement {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, &e, &p, &PrevElement)
		}

		next := Intrusive(root)
		var NextElement Intrusive
		if i < len(es)-1 {
			next = es[i+1]
			NextElement = next
		}
		if n := e.Successor(); n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, &e, &n, &next)
		}
		if n := Next(e); n != NextElement {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, &e, &n, &NextElement)
		}
	}
}

func TestList(t *testing.T) {
	l := New()
	checkListPointers(t, l, []Intrusive{})

	// Single element list
	e := l.PushFront(NewElement("a"))
	checkListPointers(t, l, []Intrusive{e})
	l.MoveToFront(e)
	checkListPointers(t, l, []Intrusive{e})
	l.MoveToBack(e)
	checkListPointers(t, l, []Intrusive{e})
	l.Remove(e)
	checkListPointers(t, l, []Intrusive{})

	// Bigger list
	e2 := l.PushFront(NewElement(2))
	e1 := l.PushFront(NewElement(1))
	e3 := l.PushBack(NewElement(3))
	e4 := l.PushBack(NewElement("banana"))
	checkListPointers(t, l, []Intrusive{e1, e2, e3, e4})

	l.Remove(e2)
	checkListPointers(t, l, []Intrusive{e1, e3, e4})

	l.MoveToFront(e3) // move from middle
	checkListPointers(t, l, []Intrusive{e3, e1, e4})

	l.MoveToFront(e1)
	l.MoveToBack(e3) // move from middle
	checkListPointers(t, l, []Intrusive{e1, e4, e3})

	l.MoveToFront(e3) // move from back
	checkListPointers(t, l, []Intrusive{e3, e1, e4})
	l.MoveToFront(e3) // should be no-op
	checkListPointers(t, l, []Intrusive{e3, e1, e4})

	l.MoveToBack(e3) // move from front
	checkListPointers(t, l, []Intrusive{e1, e4, e3})
	l.MoveToBack(e3) // should be no-op
	checkListPointers(t, l, []Intrusive{e1, e4, e3})

	e2 = l.InsertBefore(NewElement(2), e1) // insert before front
	checkListPointers(t, l, []Intrusive{e2, e1, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore(NewElement(2), e4) // insert before middle
	checkListPointers(t, l, []Intrusive{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore(NewElement(2), e3) // insert before back
	checkListPointers(t, l, []Intrusive{e1, e4, e2, e3})
	l.Remove(e2)

	e2 = l.InsertAfter(NewElement(2), e1) // insert after front
	checkListPointers(t, l, []Intrusive{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertAfter(NewElement(2), e4) // insert after middle
	checkListPointers(t, l, []Intrusive{e1, e4, e2, e3})
	l.Remove(e2)
	e2 = l.InsertAfter(NewElement(2), e3) // insert after back
	checkListPointers(t, l, []Intrusive{e1, e4, e3, e2})
	l.Remove(e2)

	// Check standard iteration.
	sum := 0
	for e := l.Front(); e != nil; e = Next(e) {
		if i, ok := Value(e).(int); ok {
			sum += i
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all elements by iterating
	var next Intrusive
	for e := l.Front(); e != nil; e = next {
		next = Next(e)
		l.Remove(e)
	}
	checkListPointers(t, l, []Intrusive{})
}

func checkList(t *testing.T, l *List, es []interface{}) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = Next(e) {
		le := Value(e).(int)
		if le != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, es[i])
		}
		i++
	}

	// NOTE: this is a change from the original test;
	// it didn't verify that the elements were actually visited.
	if i != len(es) {
		t.Errorf("didnt visit every element; have %v, what %v", i, len(es))
	}
}

func newList(ns ...int) *List {
	l := New()
	for _, n := range ns {
		l.PushBack(NewElement(n))
	}
	return l
}

// NOTE: this is a change from the original tests:
// there is no list copy, only list move.
func TestExtending(t *testing.T) {
	// move back: lists into lists
	{
		l1, l2, l3 := newList(1, 2, 3), newList(4, 5), newList()
		l3.MoveBackList(l1)
		checkList(t, l3, []interface{}{1, 2, 3})
		checkListLen(t, l1, 0)

		l3.MoveBackList(l2)
		checkList(t, l3, []interface{}{1, 2, 3, 4, 5})
		checkListLen(t, l2, 0)
	}

	// move front: lists into lists
	{
		l1, l2, l3 := newList(1, 2, 3), newList(4, 5), newList()
		l3.MoveFrontList(l2)
		checkList(t, l3, []interface{}{4, 5})
		checkListLen(t, l2, 0)

		l3.MoveFrontList(l1)
		checkList(t, l3, []interface{}{1, 2, 3, 4, 5})
		checkListLen(t, l1, 0)
	}

	// move back: lists into empty lists
	{
		l1, l3 := newList(1, 2, 3), newList()
		l3.MoveBackList(l1)
		checkList(t, l3, []interface{}{1, 2, 3})
		checkListLen(t, l1, 0)
		// no op:
		l3.MoveBackList(l3)
		checkList(t, l3, []interface{}{1, 2, 3})
	}

	// move front: lists into empty lists
	{
		l1, l3 := newList(1, 2, 3), newList()
		l3.MoveFrontList(l1)
		checkList(t, l3, []interface{}{1, 2, 3})
		checkListLen(t, l1, 0)
		// no op:
		l3.MoveFrontList(l3)
		checkList(t, l3, []interface{}{1, 2, 3})
	}

	// move back: empty lists into filled lists
	{
		l1, l3 := newList(1, 2, 3), newList()
		l1.MoveBackList(l3)
		checkList(t, l1, []interface{}{1, 2, 3})
		checkListLen(t, l3, 0)
		// no op:
		l1.MoveBackList(l3)
		checkList(t, l1, []interface{}{1, 2, 3})
		checkListLen(t, l3, 0)
	}

	// move front: empty lists into filled lists
	{
		l1, l3 := newList(1, 2, 3), newList()
		l1.MoveFrontList(l3)
		checkList(t, l1, []interface{}{1, 2, 3})
		checkListLen(t, l3, 0)
		// no op:
		l1.MoveFrontList(l3)
		checkList(t, l1, []interface{}{1, 2, 3})
		checkListLen(t, l3, 0)
	}
}

func TestRemove(t *testing.T) {
	l := New()
	e1 := l.PushBack(NewElement(1))
	e2 := l.PushBack(NewElement(2))
	checkListPointers(t, l, []Intrusive{e1, e2})
	e := l.Front()
	l.Remove(e)
	checkListPointers(t, l, []Intrusive{e2})
	l.Remove(e)
	checkListPointers(t, l, []Intrusive{e2})
}

func TestIssue4103(t *testing.T) {
	l1 := newList(1, 2)
	l2 := newList(3, 4)

	e := l1.Front()
	l2.Remove(e) // l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(NewElement(8), e)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := newList(1, 2)

	e := l.Front()
	l.Remove(e)
	if Value(e) != 1 {
		t.Errorf("e.value = %d, want 1", Value(e))
	}
	if Next(e) != nil {
		t.Errorf("Next(e) != nil")
	}
	if Prev(e) != nil {
		t.Errorf("Prev(e) != nil")
	}
}

func TestMove(t *testing.T) {
	l := New()
	e1 := l.PushBack(NewElement(1))
	e2 := l.PushBack(NewElement(2))
	e3 := l.PushBack(NewElement(3))
	e4 := l.PushBack(NewElement(4))

	l.MoveAfter(e3, e3)
	checkListPointers(t, l, []Intrusive{e1, e2, e3, e4})
	l.MoveBefore(e2, e2)
	checkListPointers(t, l, []Intrusive{e1, e2, e3, e4})

	l.MoveAfter(e3, e2)
	checkListPointers(t, l, []Intrusive{e1, e2, e3, e4})
	l.MoveBefore(e2, e3)
	checkListPointers(t, l, []Intrusive{e1, e2, e3, e4})

	l.MoveBefore(e2, e4)
	checkListPointers(t, l, []Intrusive{e1, e3, e2, e4})
	e2, e3 = e3, e2

	l.MoveBefore(e4, e1)
	checkListPointers(t, l, []Intrusive{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.MoveAfter(e4, e1)
	checkListPointers(t, l, []Intrusive{e1, e4, e2, e3})
	e2, e3, e4 = e4, e2, e3

	l.MoveAfter(e2, e3)
	checkListPointers(t, l, []Intrusive{e1, e3, e2, e4})
	e2, e3 = e3, e2
}

// Test PushFront, PushBack, MoveFrontList, MoveBackList with uninitialized List
func TestZeroList(t *testing.T) {
	var l1 = new(List)
	l1.PushFront(NewElement(1))
	checkList(t, l1, []interface{}{1})

	var l2 = new(List)
	l2.PushBack(NewElement(2))
	checkList(t, l2, []interface{}{2})

	var l3 = new(List)
	l3.MoveFrontList(l1)
	checkList(t, l3, []interface{}{1})

	var l4 = new(List)
	l4.MoveBackList(l2)
	checkList(t, l4, []interface{}{2})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	var l List
	l.PushBack(NewElement(1))
	l.PushBack(NewElement(2))
	l.PushBack(NewElement(3))
	l.InsertBefore(NewElement(1), new(Element))
	checkList(t, &l, []interface{}{1, 2, 3})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	var l List
	l.PushBack(NewElement(1))
	l.PushBack(NewElement(2))
	l.PushBack(NewElement(3))
	l.InsertAfter(NewElement(1), new(Element))
	checkList(t, &l, []interface{}{1, 2, 3})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestMoveUnknownMark(t *testing.T) {
	var l1 List
	e1 := l1.PushBack(NewElement(1))

	var l2 List
	e2 := l2.PushBack(NewElement(2))

	l1.MoveAfter(e1, e2)
	checkList(t, &l1, []interface{}{1})
	checkList(t, &l2, []interface{}{2})

	l1.MoveBefore(e1, e2)
	checkList(t, &l1, []interface{}{1})
	checkList(t, &l2, []interface{}{2})
}
