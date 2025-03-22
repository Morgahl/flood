package graph

import (
	"fmt"
	"slices"
	"strings"
)

const (
	keepCount = 250
)

type Graph struct {
	nodes     [uint16(19 * 19)]node
	solution  []uint8
	nodeCount uint16
	edgeCount uint16
}

func New(vals [19][19]uint8) *Graph {
	g := &Graph{
		solution:  make([]uint8, 0, 32),
		nodeCount: 19 * 19,
	}
	idMap := make(map[[2]int]uint16)
	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	curRow, curCol := 9, 9
	curDir := 0
	steps := 1

	for i := 0; i < 19*19; {
		for j := 0; j < 2; j++ {
			for k := 0; k < steps && i < 19*19; k++ {
				g.nodes[i] = newNode(vals[curRow][curCol])
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
	for i, node := range g.nodes {
		if node.v == 0 {
			continue
		}

		g.normalizeNode((uint16(i)), node)
	}
}

func (g *Graph) normalizeNode(id uint16, n node) (count uint16) {
START:
	for _, eId := range n.edges {
		if newId, merged := g.merge(id, eId); merged {
			count++
			if newId != id {
				id = newId
				n = g.nodes[id]
			}
			goto START
		}
	}
	return count
}

func (g *Graph) RootEdgeValues() (values map[uint8]NodeStats) {
	values = make(map[uint8]NodeStats)
	for _, eId := range g.nodes[0].edges {
		node := g.nodes[eId]
		stats := values[node.v]
		stats.Count++
		values[node.v] = stats
	}
	return values
}

func (g *Graph) Flood(v uint8) (merged bool) {
	if len(g.solution) > 0 && g.solution[len(g.solution)-1] == v {
		return false
	} else if g.edgeCount == 0 {
		return false
	}

	root := g.nodes[0]

	root.v = v
	g.nodes[0] = root
	if mergeCount := g.normalizeNode(0, root); mergeCount > 0 {
		g.solution = append(g.solution, v)
		merged = true
	}

	return merged
}

func (g *Graph) connect(id0, id1 uint16) {
	node0 := g.nodes[id0]
	node1 := g.nodes[id1]
	node0.edges = append(node0.edges, id1)
	node1.edges = append(node1.edges, id0)
	g.nodes[id0] = node0
	g.nodes[id1] = node1
	g.edgeCount++
}

func (g *Graph) merge(id0, id1 uint16) (id uint16, merged bool) {
	if id0 > id1 {
		id0, id1 = id1, id0
	}
	node0 := g.nodes[id0]
	node1 := g.nodes[id1]
	if node0.v == 0 || node1.v == 0 {
		return 65535, false
	} else if node0.v != node1.v {
		return 65535, false
	}

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
	g.nodeCount--

	g.nodes[id0] = node0
	g.nodes[id1] = node{}

	return id0, true
}

func (g *Graph) Solution() []uint8 {
	return g.solution
}

func (g *Graph) Solve() *Graph {
	if g.edgeCount == 0 {
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
				pool.PutAll(theSet[1:])
				return currentSet[0]
			}
			newSet = append(newSet, currentSet...)
			pool.Put(g)
		}

		if len(newSet) > keepCount {
			slices.SortFunc(newSet, func(g0, g1 *Graph) int {
				if edgDelta := int(g0.edgeCount) - int(g1.edgeCount); edgDelta != 0 {
					return edgDelta
				}
				return int(g0.nodeCount) - int(g1.nodeCount)
			})
			pool.PutAll(newSet[keepCount:])
			newSet = newSet[:keepCount]
		}
		theSet = newSet
	}

	return nil
}

func (g *Graph) stepEach() ([]*Graph, bool) {
	nextSet := make([]*Graph, 0, 5)
	rootValue := g.nodes[0].v
	for i := uint8(1); i <= 6; i++ {
		if i == rootValue {
			continue
		}

		ng := g.copy(true)
		if ng.Flood(i) {
			if ng.edgeCount == 0 {
				pool.PutAll(nextSet)
				return []*Graph{ng}, true
			}
			nextSet = append(nextSet, ng)
		}
	}

	return nextSet, false
}

func (g *Graph) Copy() *Graph {
	return g.copy(false)
}

func (g *Graph) copy(pooled bool) (ng *Graph) {
	if pooled {
		ng = pool.Get()
		if cap(ng.solution) < len(g.solution)+1 {
			ng.solution = make([]uint8, 0, len(g.solution)+1)
		} else {
			ng.solution = ng.solution[:0]
		}
	} else {
		ng = &Graph{
			solution: make([]uint8, 0, 32),
		}
	}
	ng.edgeCount = g.edgeCount
	ng.nodeCount = g.nodeCount
	ng.solution = append(ng.solution, g.solution...)

	i := uint16(0)
	var idMap [361]uint16
	for j, node := range g.nodes {
		if node.v == 0 {
			idMap[j] = 65535
			continue
		}
		newNode := node.Copy()
		ng.nodes[i] = newNode
		idMap[j] = i
		i++
	}

	// remap node edge ids
	for idx := uint16(0); idx < i; idx++ {
		node := ng.nodes[idx]
		for j, eId := range node.edges {
			node.edges[j] = idMap[eId]
		}
		ng.nodes[idx] = node
	}

	if ng.nodeCount != i {
		panic(fmt.Sprintf("node count mismatch: %d != %d", ng.nodeCount, i))
	}

	return ng
}

func (g *Graph) String() string {
	buf := &strings.Builder{}

	fmt.Fprintf(buf, "Nodes:\n")
	for id, n := range g.nodes {
		if n.v > 0 {
			fmt.Fprintf(buf, "%d: %v\n", id, n)
		}
	}
	fmt.Fprintf(buf, "Nodes: %d\n", g.nodeCount)
	fmt.Fprintf(buf, "Edge Count: %d\n", g.edgeCount)
	fmt.Fprintf(buf, "Root: %+v\n", g.nodes[0])

	return buf.String()
}

type node struct {
	v     uint8
	edges []uint16
}

func newNode(v uint8) node {
	return node{
		v:     v,
		edges: make([]uint16, 0, 4),
	}
}

func (n node) Copy() node {
	nn := node{
		v:     n.v,
		edges: make([]uint16, len(n.edges), cap(n.edges)),
	}

	copy(nn.edges, n.edges)

	return nn
}

type NodeStats struct {
	Count uint16
}

func (ns NodeStats) String() string {
	return fmt.Sprintf("count: %d", ns.Count)
}
