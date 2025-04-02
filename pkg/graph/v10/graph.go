package graph

import (
	"fmt"
	"slices"
	"strings"
)

const (
	keepCount = 160_000
)

type Graph struct {
	root          uint16
	rootMag       uint8
	nodeCount     uint16
	nodeCountNorm uint16
	edgeCount     uint16
	edgeCountNorm uint16
	nodes         [uint16(19 * 19)]node
	solution      []uint8
}

func New(vals [19][19]uint8) *Graph {
	g := &Graph{
		solution:  make([]uint8, 0, 32),
		root:      9*19 + 9,
		nodeCount: 19 * 19,
	}
	idx := uint16(0)
	for y, rows := range vals {
		for x, v := range rows {
			node := newNode(v, x, y, 9, 9)
			g.nodes[idx] = node
			if y > 0 {
				g.connect(idx, idx-19)
			}
			if x > 0 {
				g.connect(idx, idx-1)
			}
			idx++
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

	g.nodeCountNorm = g.nodeCount
	g.edgeCountNorm = g.edgeCount
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
	for _, eId := range g.nodes[g.root].edges {
		node := g.nodes[eId]
		stats := values[node.v]
		stats.add(node)
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

	root := g.nodes[g.root]

	root.v = v
	g.nodes[g.root] = root
	if mergeCount := g.normalizeNode(g.root, root); mergeCount > 0 {
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

	node0.mag += node1.mag
	if g.root == id1 {
		g.root = id0
		g.rootMag = node0.mag
	} else if g.root == id0 {
		g.rootMag = node0.mag
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
				if edgeDelta := int(g0.edgeCount) - int(g1.edgeCount); edgeDelta != 0 {
					// optimize for least edges first
					return edgeDelta
				} else if rootMagDelta := int(g1.rootMag) - int(g0.rootMag); rootMagDelta != 0 {
					// optimize for largest root mag second
					return rootMagDelta
				}
				// OK let's at least fall back to farthest progress....
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
	rootValue := g.nodes[g.root].v
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
	ng.rootMag = g.rootMag
	ng.edgeCount = g.edgeCount
	ng.edgeCountNorm = g.edgeCountNorm
	ng.nodeCount = g.nodeCount
	ng.nodeCountNorm = g.nodeCountNorm
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

	ng.root = idMap[g.root]

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
	fmt.Fprintf(buf, "Node Count: %d\n", g.nodeCount)
	fmt.Fprintf(buf, "Edge Count: %d\n", g.edgeCount)
	fmt.Fprintf(buf, "Root: %+v\n", g.nodes[g.root])

	return buf.String()
}

type node struct {
	edges []uint16
	v     uint8
	mag   uint8
}

func newNode(v uint8, x, y, rx, ry int) node {
	return node{
		edges: make([]uint16, 0, 4),
		v:     v,
		mag:   1,
	}
}

func (n node) Copy() node {
	nn := node{
		edges: make([]uint16, len(n.edges), cap(n.edges)),
		v:     n.v,
		mag:   n.mag,
	}

	copy(nn.edges, n.edges)

	return nn
}

type NodeStats struct {
	Count  uint16
	Mag    uint8
	edges2 map[uint16]struct{}
}

func (ns *NodeStats) add(n node) {
	ns.Count++
	ns.Mag += n.mag
	if ns.edges2 == nil && len(n.edges) > 0 {
		ns.edges2 = make(map[uint16]struct{}, len(n.edges))
	}

	for _, eId := range n.edges {
		ns.edges2[eId] = struct{}{}
	}
}

func (ns NodeStats) String() string {
	return fmt.Sprintf("count: %3d, mag: %3d, count_2nd: %3d", ns.Count, ns.Mag, len(ns.edges2))
}
