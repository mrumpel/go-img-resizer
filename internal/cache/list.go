package cache

type listInterface interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(string, string) *listItem
	Remove(*listItem)
	MoveToFront(*listItem)
}

type listItem struct {
	Key   string
	Value string
	Next  *listItem
	Prev  *listItem
}

type list struct {
	len   int
	front *listItem
	back  *listItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *listItem {
	return l.front
}

func (l list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(k, v string) *listItem {
	defer func() { l.len++ }()
	item := &listItem{
		Key:   k,
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	if l.Len() == 0 {
		l.addFirstItem(item)
		return item
	}

	l.front.Prev = item
	l.front = item

	return item
}

func (l *list) Remove(i *listItem) {
	defer func() { l.len-- }()

	if i.Next == nil && i.Prev == nil {
		l.front = nil
		l.back = nil
		return
	}

	if i.Next == nil {
		i.Prev.Next = nil
		l.back = i.Prev
		return
	}

	if i.Prev == nil {
		i.Next.Prev = nil
		l.front = i.Next
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
}

func (l *list) MoveToFront(i *listItem) {
	if i.Prev == nil {
		return
	}

	l.Remove(i)

	l.front.Prev = i

	i.Next = l.front
	i.Prev = nil

	l.front = i
	l.len++
}

func (l *list) addFirstItem(i *listItem) {
	l.front = i
	l.back = i
}

func newList() listInterface {
	return new(list)
}
