package graph

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

const (
	keepCount = 250
)

type Graph struct {
	root     uint16
	nodes    [uint16(19 * 19)]node
	edges    map[edge]struct{}
	solution []uint8
}

func New(vals [19][19]uint8) *Graph {
	g := &Graph{
		edges: map[edge]struct{}{},
	}

	idx := uint16(0)
	for y, rows := range vals {
		for x, v := range rows {
			node := newNode(idx, v)
			g.nodes[idx] = node
			if y > 0 {
				g.connect(idx, idx-19)
			}
			if x > 0 {
				g.connect(idx, idx-1)
			}
			if x == 9 && y == 9 {
				g.root = idx
			}
			idx++
		}
	}

	return g
}

func (g *Graph) Normalize() {
	for _, node := range g.nodes {
		if node.mag == 0 {
			continue
		}

		g.normalizeNode(node.id)
	}
}

func (g *Graph) normalizeNode(id uint16) (count uint16) {
START:
	for _, e := range g.edgesForNode(id) {
		if newId, merged := g.merge(e); merged {
			count++
			if newId != id {
				id = newId
				goto START
			}
		}
	}
	return count
}

func (g *Graph) RootEdgeValues() (values map[uint8]NodeStats) {
	values = make(map[uint8]NodeStats)
	for _, n := range g.edgesForNode(g.root) {
		id := n.id0
		if id == g.root {
			id = n.id1
		}
		v := g.nodes[id].v
		stats := values[v]
		stats.Count++
		stats.Mag += g.nodes[id].mag
		values[v] = stats
	}
	return values
}

func (g *Graph) Flood(v uint8) (merged bool) {
	if len(g.solution) > 0 && g.solution[len(g.solution)-1] == v {
		return false
	} else if len(g.edges) == 0 {
		return false
	}

	root := g.nodes[g.root]
	if root.v == v {
		return false
	}

	root.v = v
	g.nodes[g.root] = root
	if mergeCount := g.normalizeNode(g.root); mergeCount > 0 {
		g.solution = append(g.solution, v)
		merged = true
	}

	return merged
}

func (g *Graph) connect(id0, id1 uint16) {
	if id0 == id1 {
		panic("connecting same node")
	} else if id0 > id1 {
		id0, id1 = id1, id0
	}
	g.edges[edge{id0, id1}] = struct{}{}
}

func (g *Graph) edgesForNode(id uint16) (edges []edge) {
	for e := range g.edges {
		if e.id0 == id || e.id1 == id {
			edges = append(edges, e)
		}
	}
	return edges
}

func (g *Graph) merge(edge edge) (id uint16, merged bool) {
	node0 := g.nodes[edge.id0]
	node1 := g.nodes[edge.id1]
	if node0.v != node1.v {
		return 65535, false
	}
	delete(g.edges, edge)

	if g.root == edge.id1 {
		g.root = edge.id0
	}
	node0.mag += node1.mag
	node1.mag = 0
	g.nodes[edge.id0] = node0
	g.nodes[edge.id1] = node1

	for _, e := range g.edgesForNode(edge.id1) {
		delete(g.edges, e)
		if e.id0 == edge.id1 {
			e.id0 = edge.id0
		} else {
			e.id1 = edge.id0
		}
		g.connect(e.id0, e.id1)
	}

	return edge.id0, true
}

func (g *Graph) Solution() []uint8 {
	return g.solution
}

func (g *Graph) Solve() *Graph {
	theSet, solved := g.stepEach()
	if solved {
		return theSet[0]
	}

	for i := 1; ; i++ {
		newSet := make([]*Graph, 0, len(theSet)*5)
		for _, g := range theSet {
			currentSet, solved := g.stepEach()
			if solved {
				return currentSet[0]
			}
			newSet = append(newSet, currentSet...)
		}

		if len(newSet) > keepCount {
			slices.SortStableFunc(newSet, func(g0, g1 *Graph) int {
				return len(g0.edges) - len(g1.edges)
			})
			theSet = newSet[:keepCount]
		} else {
			theSet = newSet
		}
	}
}

func (g *Graph) stepEach() ([]*Graph, bool) {
	nextSet := make([]*Graph, 0, 5)
	rootValue := g.nodes[g.root].v
	for i := uint8(1); i <= 6; i++ {
		if i == rootValue {
			continue
		}

		ng := g.Copy()
		if ng.Flood(i) {
			if len(ng.edges) == 0 {
				return []*Graph{ng}, true
			}

			nextSet = append(nextSet, ng)
		}
	}

	return nextSet, false
}

func (g *Graph) Copy() *Graph {
	ng := &Graph{
		root:     g.root,
		nodes:    g.nodes,
		edges:    make(map[edge]struct{}, len(g.edges)),
		solution: make([]uint8, len(g.solution), len(g.solution)+1),
	}

	for e := range g.edges {
		ng.edges[e] = struct{}{}
	}

	copy(ng.solution, g.solution)

	return ng
}

func (g *Graph) String() string {
	buf := &strings.Builder{}

	fmt.Fprintf(buf, "Nodes:\n")
	nodeCount := 0
	for _, n := range g.nodes {
		if n.mag > 0 {
			nodeCount++
			fmt.Fprintf(buf, "%d: %d\n", n.id, n.v)
		}
	}
	fmt.Fprintf(buf, "Edges:\n")
	for e := range g.edges {
		fmt.Fprintf(buf, "%d %d\n", e.id0, e.id1)
	}
	fmt.Fprintf(buf, "Nodes: %d\n", nodeCount)
	fmt.Fprintf(buf, "Edges: %d\n", len(g.edges))
	fmt.Fprintf(buf, "Root: %+v\n", g.nodes[g.root])

	return buf.String()
}

type node struct {
	id  uint16
	v   uint8
	mag uint16
}

func newNode(id uint16, v uint8) node {
	return node{
		id:  id,
		v:   v,
		mag: 1,
	}
}

type edge struct {
	id0, id1 uint16
}

type NodeStats struct {
	Count uint16
	Mag   uint16
}

func (ns NodeStats) String() string {
	return fmt.Sprintf("count: %d, mag: %d, score: %f", ns.Count, ns.Mag, float64(ns.Count)*math.Sqrt(float64(ns.Mag)))
}
