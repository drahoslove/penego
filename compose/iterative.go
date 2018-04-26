package compose

import (
	// "git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)

type node struct {
	Composable
	rank     int
	position int
}
type edge struct {
	from   *node
	to     *node
	lenght int
}

type graph struct {
	nodes []*node
	edges []*edge
}

// Returns graph sturcure suited for graph related operations
// created from net.Net structure (which is more suitable for simulation and rendering purposes)
func loadGraph(network net.Net) graph {
	g := graph{
		make([]*node, 0, 16),
		make([]*edge, 0, 16),
	}
	nodeByComposable := make(map[Composable]*node)

	for _, place := range network.Places() {
		n := &node{place, 1, 1}
		g.nodes = append(g.nodes, n)
		nodeByComposable[place] = n
	}
	for _, tran := range network.Transitions() {
		n := &node{tran, 1, 1}
		g.nodes = append(g.nodes, n)
		nodeByComposable[tran] = n
	}
	for _, tran := range network.Transitions() {
		for _, arc := range tran.Origins {
			place := arc.Place
			g.edges = append(g.edges, &edge{
				nodeByComposable[place],
				nodeByComposable[tran],
				1,
			})
		}
		for _, arc := range tran.Targets {
			place := arc.Place
			g.edges = append(g.edges, &edge{
				nodeByComposable[tran],
				nodeByComposable[place],
				1,
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
		g.edges[i] = &edge{e.to, e.from, e.lenght}
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
					// reverse
					g.edges[i] = &edge{e.to, e.from, e.lenght}
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
	graph := loadGraph(network)

	_ = graph
	// assign rank λ(v) to each node v
	// edge e = (v, w)
	// lenght of edge l(e) ≥ δ(e)
	// l(e) = λ(w) − λ(v)

	// make graph acyclic by reversing edges
	graph = graph.acyclic()
	min_len := 1 // δ(e)

	// _ = min_len

	initRank := func() {
		// An initial feasible ranking is computed. For brevity,init_rankis not given here.  Our versionkeeps nodes in a queue.  Nodes are placed in the queue when they have no unscanned in-edges.As nodes are taken off the queue, they are assigned the least rank that satisfies their in-edges, andtheir out-edges are marked as scanned. In the simplest case, whereδ =1 for all edges, thiscorresponds to viewing the graph as a poset and assigning the minimal elements to rank 0.  Thesenodes are removed from the poset and the new set of minimal elements are assigned rank 1, etc.
	}
	tightTree := func() int {
		// finds a maximal tree of tight edges containing some fixed node andreturns the number of nodes in the tree. Note that such a maximal tree is just a spanning tree forthe subgraph induced by all nodes reachable from the fixed node in the underlying undirected
		// graph using only tight edges.  In particular, all such trees have the same number of nodes.
		return 0
	}
	initCutValues := func() {
		// computes the cut values of the tree edges. For each tree edge,this is computed by marking the nodes as belonging to the head or tail component, and thenperforming the sum of the signed weights of all edges whose head and tail are in differentcomponents, the sign being negative for those edges going from the head to the tail component
	}

	feasibleTree := func() {
		initRank()
		for tightTree() < len(graph.nodes) {
			e := (*edge)(nil)           // a non-tree edge incident on the tree with a minimal amount of slack
			delta := e.lenght - min_len // slack of an edge isthe difference of its length and its minimum length.
			if incident_node == e.from {
				delta = -delta
			}
			for _, v := range Tree {
				v.rank += delta
			}
		}
		initCutValues()
	}

	rank := func() {
		feasibleTree()
		for {
			e := leaveEdge()
			if e == nil {
				break
			}
			f := enterEdge(e)
			exchange(e, f)
		}
		normalize()
		balance()
	}

	ordering := func() {

	}

	position := func() {

	}

	makeSplines := func() {

	}

	rank()
	ordering()
	position()
	makeSplines()
	return Composition{}
}
