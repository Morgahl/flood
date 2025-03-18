package board

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Morgahl/flood/pkg/color"
)

const (
	keepCount = 1000
)

type Board struct {
	cells       [19][19]*cell
	solveFor    []uint8
	solution    []uint8
	solvedCount int
}

func New(vals [19][19]uint8) *Board {
	b := &Board{}

	// setup cells
	for yi := 0; yi < 19; yi++ {
		for xi := 0; xi < 19; xi++ {
			v := vals[yi][xi]
			n := newCell(v)
			b.cells[yi][xi] = n
			b.addSolveFor(v)
			n.link(up, b.getNodeAt(yi-1, xi))
			n.link(left, b.getNodeAt(yi, xi-1))
		}
	}

	return b
}

func (b *Board) Normalize() {
	b.Flood(b.cells[9][9].v)
}

func (b *Board) Solution() []uint8 {
	return b.solution
}

func (b *Board) Copy() *Board {
	nb := &Board{
		solveFor:    b.solveFor,
		solution:    make([]uint8, len(b.solution), len(b.solution)+1),
		solvedCount: b.solvedCount,
	}

	copy(nb.solution, b.solution)

	for yi := 0; yi < 19; yi++ {
		for xi := 0; xi < 19; xi++ {
			n := newCell(b.cells[yi][xi].v)
			nb.cells[yi][xi] = n
			n.link(up, nb.getNodeAt(yi-1, xi))
			n.link(left, nb.getNodeAt(yi, xi-1))
		}
	}

	return nb
}

func (b *Board) addSolveFor(val uint8) {
	for _, sf := range b.solveFor {
		if sf == val {
			return
		}
	}

	b.solveFor = append(b.solveFor, val)
}

func (b *Board) Solve() *Board {
	theSet, solved := b.solveEach()
	if solved {
		return theSet[0]
	}

	for i := 0; ; i++ {
		newSet := make([]*Board, 0, len(theSet)*5)
		for _, nb := range theSet {
			currentSet, solved := nb.solveEach()
			if solved {
				return currentSet[0]
			}
			newSet = append(newSet, currentSet...)
		}

		if len(newSet) > keepCount {
			sort.Slice(newSet, func(i, j int) bool {
				return newSet[i].solvedCount > newSet[j].solvedCount
			})
			theSet = newSet[:keepCount]
		} else {
			theSet = newSet
		}
	}
}

func (b *Board) solveEach() ([]*Board, bool) {
	nextSet := make([]*Board, 0, len(b.solveFor))
	for _, sf := range b.solveFor {
		if sf == b.cells[9][9].v {
			continue
		}

		nb := b.Copy()
		nb.Flood(sf)
		if nb.solvedCount == 19*19 {
			return []*Board{nb}, true
		}

		nextSet = append(nextSet, nb)
	}

	return nextSet, false
}

func (b *Board) Flood(t uint8) {
	if len(b.solution) > 0 && b.solution[len(b.solution)-1] == t {
		return
	} else if b.solvedCount == 19*19 {
		return
	}

	if solvedCount := b.cells[9][9].flood(b.cells[9][9].v, t); b.solvedCount < solvedCount {
		if b.solvedCount > 0 {
			b.solution = append(b.solution, t)
		}
		b.solvedCount = solvedCount
	}
	b.clear()
}

func (b *Board) getNodeAt(y, x int) *cell {
	if x < 0 ||
		19 <= x ||
		y < 0 ||
		19 <= y {
		return nil
	}

	return b.cells[y][x]
}

func (b *Board) clear() {
	for yi := 0; yi < 19; yi++ {
		for xi := 0; xi < 19; xi++ {
			b.cells[yi][xi].mark = false
		}
	}
}

func (b *Board) ANSIString() string {
	builder := &strings.Builder{}
	for i, r := range b.cells {
		if i == 9 {
			for _, c := range r[0:8] {
				builder.WriteString(color.Colorize(c.v))
				builder.WriteByte(' ')
			}
			builder.WriteString(color.Colorize(r[8].v))
			builder.WriteByte('[')
			builder.WriteString(color.Colorize(r[9].v))
			builder.WriteByte(']')
			for _, c := range r[10:] {
				builder.WriteString(color.Colorize(c.v))
				builder.WriteByte(' ')
			}
		} else {
			for _, c := range r {
				builder.WriteString(color.Colorize(c.v))
				builder.WriteByte(' ')
			}
		}
		builder.WriteByte('\n')
	}

	builder.WriteByte('[')
	for i, v := range b.solution {
		builder.WriteString(color.Colorize(v))
		if i < len(b.solution)-1 {
			builder.WriteByte(' ')
		}
	}
	builder.WriteString("] ")

	builder.WriteString(strconv.Itoa(len(b.solution)))
	builder.WriteByte(' ')
	builder.WriteString(strconv.Itoa(b.solvedCount))

	return builder.String()
}

func (b *Board) String() string {
	builder := &strings.Builder{}
	for _, r := range b.cells {
		for _, c := range r {
			builder.WriteString(c.String())
			builder.WriteByte(' ')
		}
		builder.WriteByte('\n')
	}

	builder.WriteByte('[')
	for _, v := range b.solution {
		builder.WriteString(strconv.Itoa(int(v)))
		builder.WriteByte(' ')
	}
	builder.WriteString("] ")

	builder.WriteString(strconv.Itoa(len(b.solution)))
	builder.WriteByte(' ')
	builder.WriteString(strconv.Itoa(b.solvedCount))

	return builder.String()
}
