package graph

import (
	"fmt"
	"strings"
)

type node struct {
	id int
	v  uint8
}

func (n node) copy() node {
	// If the node ever contains a reference to another node, this will need to be
	// changed to a deep copy.
	return n
}

func newNode(id int, v uint8) node {
	return node{
		id,
		v,
	}
}

type edge struct {
	id0, id1 int
}

type Graph struct {
	nodes map[int]node
	edges map[edge]struct{}
}

func (g *Graph) Copy() *Graph {
	ng := &Graph{
		nodes: make(map[int]node, len(g.nodes)),
		edges: make(map[edge]struct{}, len(g.edges)),
	}
	for id, n := range g.nodes {
		ng.nodes[id] = n.copy()
	}
	for e := range g.edges {
		ng.edges[e] = struct{}{}
	}
	return ng
}

func New(vals [19][19]uint8) *Graph {
	g := &Graph{
		nodes: make(map[int]node, 19*19),
		edges: map[edge]struct{}{},
	}

	for y, rows := range vals {
		for x, v := range rows {
			idx := y*19 + x
			g.nodes[idx] = newNode(idx, v)
			if y > 0 {
				g.connect(idx, idx-19)
			}
			if x > 0 {
				g.connect(idx, idx-1)
			}
		}
	}

	g.normalize()

	return g
}

func (g *Graph) String() string {
	buf := &strings.Builder{}

	fmt.Fprintf(buf, "Nodes: %d\n", len(g.nodes))
	for _, n := range g.nodes {
		fmt.Fprintf(buf, "%d: %d\n", n.id, n.v)
	}
	fmt.Fprintf(buf, "Edges: %d\n", len(g.edges))
	for e := range g.edges {
		fmt.Fprintf(buf, "%d %d\n", e.id0, e.id1)
	}

	return buf.String()
}

func (g *Graph) normalize() {
START:
	for e := range g.edges {
		if g.merge(e.id0, e.id1) {
			goto START
		}
	}
}

func (g *Graph) connect(id0, id1 int) {
	if id0 > id1 {
		id0, id1 = id1, id0
	}
	g.edges[edge{id0, id1}] = struct{}{}
}

func (g *Graph) edgesForNode(id int) (edges []edge) {
	for e := range g.edges {
		if e.id0 == id || e.id1 == id {
			edges = append(edges, e)
		}
	}
	return edges
}

func (g *Graph) merge(id0, id1 int) (merged bool) {
	if id0 > id1 {
		id0, id1 = id1, id0
	}

	node0 := g.nodes[id0]
	node1 := g.nodes[id1]
	if node0.v != node1.v {
		return false
	}

	delete(g.edges, edge{id0, id1})
	delete(g.nodes, id1)

	for _, edge := range g.edgesForNode(id1) {
		delete(g.edges, edge)
		if edge.id0 == id1 {
			g.connect(id0, edge.id1)
		} else {
			g.connect(id0, edge.id0)
		}
	}

	return true
}
