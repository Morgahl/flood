package buffer

// const (
// 	GROWTH_THRESHOLD = 1024
// 	MIN_SIZE         = 2
// )

// type Circular[T any] struct {
// 	getIdx int
// 	putIdx int
// 	elem   []*T
// 	new    func() *T
// }

// func NewCircular[C func() *T, T any](size uint C) *Circular[T] {
// 	if size < MIN_SIZE {
// 		size = MIN_SIZE
// 	}
// 	return &Circular[T]{
// 		new:    new,
// 		getIdx: -1,
// 		putIdx: 0,
// 		elem:   make([]*T, size),
// 	}
// }

// func (c *Circular[T]) Get() *T {
// 	if c.getIdx == c.putIdx || c.getIdx == -1 {
// 		return c.new()
// 	}
// 	elem := c.elem[c.getIdx]
// 	c.elem[c.getIdx] = nil
// 	c.getIdx++
// 	if c.getIdx == len(c.elem) {
// 		c.getIdx = 0
// 	}
// 	return elem
// }

// func (c *Circular[T]) GetN(n int) (elems []*T) {
// 	for i := 0; i < n; i++ {
// 		elems = append(elems, c.Get())
// 	}
// 	return elems
// }

// func (c *Circular[T]) Put(elem *T) {
// 	if c.putIdx == c.getIdx && c.elem[c.putIdx] != nil {
// 		c.grow()
// 	}
// 	c.elem[c.putIdx] = elem
// 	c.putIdx++
// 	if c.putIdx == len(c.elem) {
// 		c.putIdx = 0
// 	}
// }

// func (c *Circular[T]) PutAll(elems []*T) {
// 	for _, elem := range elems {
// 		c.Put(elem)
// 	}
// }

// func (c *Circular[T]) grow() {
// 	var n int
// 	l := len(c.elem)
// 	if l >= 1024 {
// 		n = l / 4
// 	}
// 	newElem := make([]*T, l+n)
// 	if c.putIdx > c.getIdx {
// 		copy(newElem, c.elem[c.getIdx:c.putIdx])
// 	} else {
// 		copy(newElem, c.elem[c.getIdx:])
// 		copy(newElem[l-c.getIdx:], c.elem[:c.putIdx])
// 	}
// 	c.putIdx = l
// 	c.getIdx = 0
// 	c.elem = newElem
// }

type Circular[T any] struct {
	isFull bool
	start  int
	end    int
	new    func() T
	data   []T
}

func NewCircular[T any](size int, new func() T) *Circular[T] {
	return &Circular[T]{
		data:  make([]T, size),
		start: 0,
		end:   0,
		new:   new,
	}
}

func (c *Circular[T]) Put(t T) {
	if c.isFull {
		c.grow()
	}

	c.data[c.end] = t
	c.end = (c.end + 1) % len(c.data)
	c.isFull = c.end == c.start
}

func (c *Circular[T]) PutAll(ts ...T) {
	var nt T
	for i, t := range ts {
		c.Put(t)
		ts[i] = nt
	}
}

func (c *Circular[T]) Get() (t T) {
	var nt T
	if !c.isFull && c.start == c.end {
		return c.new()
	}

	t = c.data[c.start]
	c.data[c.start] = nt
	c.start = (c.start + 1) % len(c.data)
	c.isFull = false
	return t
}

func (c *Circular[T]) grow() {
	n := len(c.data)
	if n >= 1024 {
		n = n / 4
	}
	newData := make([]T, len(c.data)+n)
	if c.start < c.end {
		copy(newData, c.data[c.start:c.end])
	} else {
		copy(newData, c.data[c.start:])
		copy(newData[len(c.data)-c.start:], c.data[:c.end])
	}
	c.start = 0
	c.end = len(c.data)
	c.data = newData
	c.isFull = false
}
