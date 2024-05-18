package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Key   Key
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	count int
	first *ListItem
	last  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
	}

	if l.count == 0 {
		l.first = newItem
		l.last = newItem
	} else {
		newItem.Next = l.first
		l.first.Prev = newItem
		l.first = newItem
	}
	l.count++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
	}

	if l.count == 0 {
		l.first = newItem
		l.last = newItem
	} else {
		newItem.Prev = l.last
		l.last.Next = newItem
		l.last = newItem
	}
	l.count++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	next := i.Next
	prev := i.Prev

	if next != nil {
		if prev != nil {
			next.Prev = prev
		} else {
			next.Prev = nil
		}
	}

	if prev != nil {
		if next != nil {
			prev.Next = next
		} else {
			prev.Next = nil
		}
	}

	if l.first == i {
		if l.first.Next != nil {
			l.first = l.first.Next
		} else {
			l.first = nil
		}
	}

	if l.last == i {
		if l.last.Prev != nil {
			l.last = l.last.Prev
		} else {
			l.last = nil
		}
	}

	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.first == i {
		return
	}
	l.Remove(i)
	i.Prev = nil
	i.Next = l.first
	l.first.Prev = i
	l.first = i
	l.count++
}
