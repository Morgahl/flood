package graph

import "sync"

var pool = newPool()

type Pool struct {
	pool sync.Pool
}

func newPool() *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() any {
				return &Graph{
					solution: make([]uint8, 0, 32),
				}
			},
		},
	}
}

func (p *Pool) Get() *Graph {
	g := p.pool.Get().(*Graph)
	g.solution = g.solution[:0]
	g.edgeCount = 0
	for i := range g.nodes {
		g.nodes[i] = node{} // zero out nodes
	}
	return g
}

func (p *Pool) Put(g *Graph) {
	p.pool.Put(g)
}

func (p *Pool) PutAll(gs []*Graph) {
	for _, g := range gs {
		p.pool.Put(g)
	}
}
