package board

import (
	"fmt"
	"sort"
)

const (
	keepCount = 250
)

type Board struct {
	cells       [19][19]*cell
	solveFor    []uint8
	solution    []uint8
	solvedCount int
}

func New(vals [][]uint8) *Board {
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

func (b *Board) Solution() []uint8 {
	return b.solution
}

func (b *Board) copy() *Board {
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

		nb := b.copy()
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
	}
	b.solvedCount = b.cells[9][9].flood(b.cells[9][9].v, t)
	b.solution = append(b.solution, b.cells[9][9].v)
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

func (b *Board) String() string {
	var display string
	for _, r := range b.cells {
		for _, c := range r {
			display += c.String() + " "
		}
		display += "\n"
	}
	display += fmt.Sprint(b.solution, b.solvedCount)

	return display
}
