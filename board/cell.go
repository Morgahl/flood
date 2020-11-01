package board

import (
	"fmt"
)

type direction int

const (
	up    direction = iota //0
	right                  //1
	down                   //2
	left                   //3
)

type cell struct {
	v    uint8    // value
	nc   [4]*cell // neighbors in u, r, d, l order
	mark bool
}

func newCell(v uint8) *cell {
	return &cell{
		v: v,
	}
}

func (c *cell) link(d direction, nc *cell) {
	if c != nil {
		c.nc[d] = nc
	}

	if nc != nil {
		switch d {
		case up:
			nc.nc[down] = c
		case down:
			nc.nc[up] = c
		case left:
			nc.nc[right] = c
		case right:
			nc.nc[left] = c
		}
	}
}

func (c *cell) value() uint8 {
	if c == nil {
		return 0
	}

	return c.v
}

func (c *cell) flood(f, t uint8) (count int) {
	if c == nil {
		return 0
	} else if c.mark {
		return 0
	}

	c.mark = true

	if c.v != f && c.v != t {
		return 0
	}

	count++

	if c.v == f {
		c.v = t
		for _, nn := range c.nc {
			count += nn.flood(f, t)
		}
	}

	return count
}

func (c *cell) String() string {
	return fmt.Sprintf("%d", c.v)
}
