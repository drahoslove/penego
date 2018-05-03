package compose

import (
	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
	"log"
)

const maxInt = int(^uint(0) >> 1)
const min_len = 1

type node struct {
	Composable
	rank     int
	position int
	priority int
	tree     *tree // for finding tight tree
	// related to cut value computing:
	parent *edge
	lim    int
	low    int
}
type edge struct {
	from     *node
	to       *node
	cutValue int
}

func (e *edge) weight() int {
	return 1
}

func (e *edge) len() int {
	return e.to.rank - e.from.rank
}

func (e *edge) slack() int {
	return e.len() - min_len
}

type graph struct {
	nodes []*node
	edges []*edge
}

func newGraph() graph {
	return graph{
		make([]*node, 0, 16),
		make([]*edge, 0, 16),
	}
}

func (g *graph) addEdge(e *edge) {
	g.edges = append(g.edges, e)
	if !g.includesNode(e.from) {
		g.nodes = append(g.nodes, e.from)
	}
	if !g.includesNode(e.to) {
		g.nodes = append(g.nodes, e.to)
	}
}

func (g *graph) includesNode(n *node) bool {
	for _, no := range g.nodes {
		if n == no {
			return true
		}
	}
	return false
}

func (g *graph) includesEdge(e *edge) bool {
	for _, ed := range g.edges {
		if e == ed {
			return true
		}
	}
	return false
}

func (g *graph) minIncidentEdge(edges []*edge) (min_e *edge, min_n *node) {
	min_slack := maxInt
	for _, e := range edges {
		if e.slack() < min_slack {
			from_is_in := g.includesNode(e.from)
			to_is_in := g.includesNode(e.to)
			if from_is_in && !to_is_in {
				min_slack = e.slack()
				min_e = e
				min_n = e.to
			}
			if to_is_in && !from_is_in {
				min_slack = e.slack()
				min_e = e
				min_n = e.from
			}
		}
	}
	return
}

func (g *graph) normalizeRanks() {
	minRank := maxInt
	for _, n := range g.nodes {
		if n.rank < minRank {
			minRank = n.rank
		}
	}
	if minRank != 0 {
		for _, n := range g.nodes {
			n.rank -= minRank
		}
	}
}

func (e *edge) reversed() *edge {
	re := *e
	re.from, re.to = e.to, e.from
	return &re
}

func (n *node) inEdges(edges []*edge) []*edge {
	inEdges := []*edge{}
	for _, e := range edges {
		if e.to == n {
			inEdges = append(inEdges, e)
		}
	}
	return inEdges
}

func (n *node) outEdges(edges []*edge) []*edge {
	outEdges := []*edge{}
	for _, e := range edges {
		if e.from == n {
			outEdges = append(outEdges, e)
		}
	}
	return outEdges
}

// Returns graph sturcure suited for graph related operations
// created from net.Net structure (which is more suitable for simulation and rendering purposes)
func loadGraph(network net.Net) graph {
	g := newGraph()
	nodeByComposable := make(map[Composable]*node)

	for _, place := range network.Places() {
		n := &node{Composable: place}
		g.nodes = append(g.nodes, n)
		nodeByComposable[place] = n
	}
	for _, tran := range network.Transitions() {
		n := &node{Composable: tran}
		g.nodes = append(g.nodes, n)
		nodeByComposable[tran] = n
	}
	for _, tran := range network.Transitions() {
		for _, arc := range tran.Origins {
			if arc.IsDumb() {
				continue
			}
			place := arc.Place
			g.edges = append(g.edges, &edge{
				from: nodeByComposable[place],
				to:   nodeByComposable[tran],
			})
		}
		for _, arc := range tran.Targets {
			if arc.IsDumb() {
				continue
			}
			place := arc.Place
			g.edges = append(g.edges, &edge{
				from: nodeByComposable[tran],
				to:   nodeByComposable[place],
			})
		}
	}
	return g
}

// returns new graph with inverted edges
func (g graph) transpose() graph {
	edges := make([]*edge, len(g.edges))
	copy(edges, g.edges)
	g.edges = edges
	for i, e := range g.edges {
		g.edges[i] = e.reversed()
	}
	return g
}

// depth first search algorithm
// only enter to node if cond is fulfilled
// calls onopen for each node when opened
// calls onclose for each node when closed
func dfs(g graph, v0 *node,
	cond func(v *node) bool,
	onopen func(v *node),
	onclose func(v *node),
) (in, out map[*node]int) {
	const (
		notfound int = iota
		open
		closed
	)
	state := map[*node]int{}
	in = map[*node]int{}
	out = map[*node]int{}
	for _, node := range g.nodes {
		state[node] = notfound
		in[node] = 0
		out[node] = 0
	}
	step := 0

	var dfs2 func(*node)
	dfs2 = func(v *node) {
		if !cond(v) {
			return
		}
		if onopen != nil {
			onopen(v)
		}
		state[v] = open
		step++
		in[v] = step

		for _, edge := range g.edges {
			if edge.from == v {
				w := edge.to
				if state[w] == notfound {
					dfs2(w)
				}
			}
		}

		state[v] = closed
		step++
		out[v] = step
		if onclose != nil {
			onclose(v)
		}
	}

	dfs2(v0)
	return
}

// return all nontrivial strongly connected components of graph
// Kosaraju's algorithm
func (g graph) components() (map[*node]*node, []*node) {
	gt := g.transpose()
	stack := []*node{}
	visited := map[*node]bool{}
	components := map[*node]*node{}
	isCompTimes := map[*node]int{}

	notVisited := func(v *node) bool { return !visited[v] }
	markVisited := func(v *node) { visited[v] = true }
	addToStack := func(v *node) { stack = append(stack, v) }

	for _, v := range gt.nodes {
		dfs(gt, v, notVisited, markVisited, addToStack)
	}

	notAssigned := func(v *node) bool {
		_, isIn := components[v]
		return !isIn
	}
	for i := len(stack) - 1; i >= 0; i-- {
		v := stack[i]
		assignComp := func(w *node) { components[w] = v; isCompTimes[v]++ }
		dfs(g, v, notAssigned, assignComp, nil)
	}

	// count components which consists of more than one node
	nontrivials := []*node{}
	for v, n := range isCompTimes {
		if n > 1 {
			nontrivials = append(nontrivials, v)
		}
	}

	return components, nontrivials
}

// returns new graph without cycles
// this is done by reverting backwards edges
func (g graph) acyclic() graph {
	edges := make([]*edge, len(g.edges))
	copy(edges, g.edges)
	g.edges = edges

	visited := map[*node]bool{}
	stack := map[*node]bool{}

	var dfs func(v *node)
	dfs = func(v *node) {
		if visited[v] {
			return
		}
		visited[v] = true
		stack[v] = true
		for i, e := range g.edges {
			if e.from == v {
				u := e.to
				if stack[u] {
					g.edges[i] = e.reversed()
				} else {
					if !visited[u] {
						dfs(u)
					}
				}
			}
		}
		stack[v] = false
	}

	_, components := g.components()
	for _, comp := range components {
		dfs(comp)
	}

	return g
}

// Iterative method for graph drawing
// based on dot algorithm and work of Warfield, Sugiyamaet at al.
func GetIterative(network net.Net) Composition {
	position := func(g *graph) {
		// TODO implement better - this is just dumb way
		ranksToPos := map[int]int{}
		for _, n := range g.nodes {
			if _, ok := ranksToPos[n.rank]; ok {
				ranksToPos[n.rank] = 0
			} else {
				ranksToPos[n.rank] += 1
			}
			n.position = ranksToPos[n.rank]
			// log.Println("node rank and position", n.rank, n.position)
			_ = log.Println
		}
	}

	graph := loadGraph(network)

	// assign rank λ(v) to each node v
	// edge e = (v, w)
	// lenght of edge l(e) ≥ δ(e)
	// l(e) = λ(w) − λ(v)

	// make graph acyclic by reversing edges
	graph = graph.acyclic()

	makeSplines := func() {

	}

	rank(&graph)
	ordering(&graph)
	position(&graph)
	makeSplines()

	comp := New()

	for _, n := range graph.nodes {
		pos := draw.Pos{float64(n.rank * 90), float64(n.position) * 90}
		switch node := n.Composable.(type) {
		case *net.Transition:
			comp.transitions[node] = pos
		case *net.Place:
			comp.places[node] = pos
		}
	}
	comp.CenterTo(0, 0)

	return comp
}

func ordering(g *graph) {
	fillPathNodes(g)
}

// adds intermediate "path nodes" to edges whose ranks are not adjencent
func fillPathNodes(g *graph) {
	for _, e := range g.edges {
		addNodeEdge := func(rank int) {
			inter_n := &node{Composable: nil, rank: rank}
			inter_e := &edge{from: inter_n, to: e.to}
			e.to = inter_n
			g.nodes = append(g.nodes, inter_n)
			g.edges = append(g.edges, inter_e)
		}
		if e.from.rank < e.to.rank {
			// add new node and edge []++> from right end
			// (1)-->(4)
			// (1)-->[3]++>(4)
			// (1)-->[2]++>(3)-->(4)
			for rank := e.to.rank - 1; rank > e.from.rank; rank-- {
				addNodeEdge(rank)
			}
		} else {
			// add new node and edge []++> from right end
			// (4)-->(1)
			// (4)-->[2]++>(1)
			// (4)-->[3]++>(2)-->(1)
			for rank := e.to.rank + 1; rank < e.from.rank; rank++ {
				addNodeEdge(rank)
			}
		}
	}
}
