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
	root      uint16
	nodes     [uint16(19 * 19)]node
	solution  []uint8
	edgeCount uint16
}

func New(vals [19][19]uint8) *Graph {
	g := &Graph{}
	idMap := make(map[[2]int]uint16)
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	curRow, curCol := 9, 9
	curDir := 0
	steps := 1

	for i := 0; i < 19*19; {
		for j := 0; j < 2; j++ {
			for k := 0; k < steps && i < 19*19; k++ {
				g.nodes[i] = newNode(uint16(i), vals[curRow][curCol])
				idMap[[2]int{curRow, curCol}] = uint16(i)
				i++
				curRow += directions[curDir][0]
				curCol += directions[curDir][1]
			}
			curDir = (curDir + 1) % 4
		}
		steps++
	}

	for y, rows := range vals {
		for x := range rows {
			idx := idMap[[2]int{y, x}]
			if y > 0 {
				g.connect(idx, idMap[[2]int{y - 1, x}])
			}
			if x > 0 {
				g.connect(idx, idMap[[2]int{y, x - 1}])
			}
		}
	}

	return g
}

func (g *Graph) Normalize() {
	for _, node := range g.nodes {
		if node.mag == 0 {
			continue
		}

		g.normalizeNode(node)
	}
}

func (g *Graph) normalizeNode(n node) (count uint16) {
START:
	for _, eId := range n.edges {
		if newId, merged := g.merge(n.id, eId); merged {
			count++
			if newId != n.id {
				n = g.nodes[newId]
			}
			goto START
		}
	}
	return count
}

func (g *Graph) RootEdgeValues() (values map[uint8]NodeStats) {
	values = make(map[uint8]NodeStats)
	for _, eId := range g.nodes[g.root].edges {
		node := g.nodes[eId]
		stats := values[node.v]
		stats.Mag += node.mag
		stats.Count++
		values[node.v] = stats
	}
	return values
}

func (g *Graph) Flood(v uint8) (merged bool) {
	if len(g.solution) > 0 && g.solution[len(g.solution)-1] == v {
		return false
	}

	root := g.nodes[g.root]
	if root.mag == 0 || root.mag == 19*19 {
		return false
	}

	root.v = v
	g.nodes[g.root] = root
	if mergeCount := g.normalizeNode(root); mergeCount > 0 {
		g.solution = append(g.solution, v)
		merged = true
	}

	return merged
}

func (g *Graph) connect(id0, id1 uint16) {
	if id0 == id1 {
		panic("connecting same node")
	}
	node0 := g.nodes[id0]
	node1 := g.nodes[id1]
	node0.edges = append(node0.edges, id1)
	node1.edges = append(node1.edges, id0)
	g.nodes[id0] = node0
	g.nodes[id1] = node1
	g.edgeCount++
}

func (g *Graph) merge(id0, id1 uint16) (id uint16, merged bool) {
	if id0 == id1 {
		panic(fmt.Sprintf("merging same node:\nid0:\t%v\nid1:\t%v\nroot:\t%v", g.nodes[id0], g.nodes[id1], g.root))
	} else if id0 > id1 {
		id0, id1 = id1, id0
	}
	node0 := g.nodes[id0]
	node1 := g.nodes[id1]
	if node0.mag == 0 || node1.mag == 0 {
		return 65535, false
	} else if node0.v != node1.v {
		return 65535, false
	}

	if g.root == id1 {
		g.root = id0
	}

	node0.mag += node1.mag
	node1.mag = 0

	// delete node1 edge on node0
	for i, id := range node0.edges {
		if id == id1 {
			lastIdx := len(node0.edges) - 1
			node0.edges[i] = node0.edges[lastIdx]
			node0.edges = node0.edges[:lastIdx]
			break
		}
	}

	// merge edges from node1 to node0 and remove node1 from its former edges
	for _, id := range node1.edges {
		if id == id0 {
			// skip edge to node0
			continue
		}

		exists := false
		for _, eId := range node0.edges {
			if eId == id {
				exists = true
				break
			}
		}

		node := g.nodes[id]
		for i, eId := range node.edges {
			if eId == id1 {
				lastIdx := len(node.edges) - 1
				node.edges[i] = node.edges[lastIdx]
				node.edges = node.edges[:lastIdx]
				break
			}
		}

		if !exists {
			node0.edges = append(node0.edges, id)
			node.edges = append(node.edges, id0)
		} else {
			g.edgeCount--
		}

		g.nodes[id] = node
	}

	g.edgeCount--

	node1.edges = nil
	g.nodes[id0] = node0
	g.nodes[id1] = node1

	return id0, true
}

func (g *Graph) Solution() []uint8 {
	return g.solution
}

func (g *Graph) Solve() *Graph {
	if g.nodes[g.root].mag == 19*19 {
		return g
	}

	theSet, solved := g.stepEach()
	if solved {
		return theSet[0]
	}

	for i := 1; i < 108; i++ {
		newSet := make([]*Graph, 0, len(theSet)*5)
		for _, g := range theSet {
			currentSet, solved := g.stepEach()
			if solved {
				return currentSet[0]
			}
			newSet = append(newSet, currentSet...)
		}

		if len(newSet) > keepCount {
			slices.SortFunc(newSet, func(g0, g1 *Graph) int {
				if edgeCount := int(g0.edgeCount) - int(g1.edgeCount); edgeCount != 0 {
					return edgeCount
				}
				return int(g1.nodes[g1.root].mag) - int(g0.nodes[g0.root].mag)
			})
			theSet = newSet[:keepCount]
		} else {
			theSet = newSet
		}
	}

	return theSet[0]
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
			if afterMag := ng.nodes[ng.root].mag; afterMag == 19*19 {
				return []*Graph{ng}, true
			}
			nextSet = append(nextSet, ng)
		}
	}

	return nextSet, false
}

func (g *Graph) Copy() *Graph {
	ng := &Graph{
		solution:  make([]uint8, len(g.solution), len(g.solution)+1),
		edgeCount: g.edgeCount,
	}

	copy(ng.solution, g.solution)

	i := uint16(0)
	var idMap [361]uint16
	for j, node := range g.nodes {
		if node.mag == 0 {
			idMap[j] = 65535
			continue
		}
		newNode := node.Copy()
		newNode.id = i
		ng.nodes[i] = newNode
		idMap[node.id] = i
		i++
	}

	// remap node ids and edges
	// for idx, node := range view {
	for idx := uint16(0); idx < i; idx++ {
		node := ng.nodes[idx]
		node.id = uint16(idx)
		for j, eId := range node.edges {
			node.edges[j] = idMap[eId]
		}
		ng.nodes[idx] = node
	}

	ng.root = idMap[g.root]

	return ng
}

func (g *Graph) String() string {
	buf := &strings.Builder{}

	fmt.Fprintf(buf, "Nodes:\n")
	nodeCount := 0
	for _, n := range g.nodes {
		if n.mag > 0 {
			nodeCount++
			fmt.Fprintf(buf, "%d: %v\n", n.id, n)
		}
	}
	fmt.Fprintf(buf, "Nodes: %d\n", nodeCount)
	fmt.Fprintf(buf, "Edge Count: %d\n", g.edgeCount)
	fmt.Fprintf(buf, "Root: %+v\n", g.nodes[g.root])

	return buf.String()
}

type node struct {
	id    uint16
	v     uint8
	mag   uint16
	edges []uint16
}

func newNode(id uint16, v uint8) node {
	return node{
		id:    id,
		v:     v,
		mag:   1,
		edges: make([]uint16, 0, 16),
	}
}

func (n node) Copy() node {
	nn := node{
		id:    n.id,
		v:     n.v,
		mag:   n.mag,
		edges: make([]uint16, len(n.edges), cap(n.edges)),
	}

	copy(nn.edges, n.edges)

	return nn
}

type NodeStats struct {
	Count uint16
	Mag   uint16
}

func (ns NodeStats) String() string {
	return fmt.Sprintf("count: %d, mag: %d, score: %f", ns.Count, ns.Mag, float64(ns.Count)*math.Sqrt(float64(ns.Mag)))
}
