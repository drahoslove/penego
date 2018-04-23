package compose

import (
	// "git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)


type node struct {
	Composable
	rank int
	position int
}
type edge struct {
	from *node
	to *node
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
	onclose func(v *node),
	onopen func(v *node),
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
	for i := len(stack)-1; i >= 0; i-- {
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



// Iterative method for graph drawing
// based on dot algorithm and work of Warfield, Sugiyamaet at al.  
func GetIterative(network net.Net) Composition {
	graph := loadGraph(network)

	_ = graph
	// assign rank λ(v) to each node v
	// edge e = (v, w)
	// lenght of edge l(e) ≥ δ(e)
	// l(e) = λ(w) − λ(v)
	rank := func() {
		min_len := 1 // δ(e)

		_ = min_len
		// make graph acyclic by reversing edges

		// tree edges
		// nontree edges
		//    cross
		//    forward
		//    back - these must be reversed to be forward

		// projit vsechny netriviální silně souvislé komponenty
		// depth-first 
		// count how many cycles each edge makes in each komponent

	}

	ordering := func() {

	}

	position := func() {

	}

	make_splines := func() {

	}

	rank()
	ordering()
	position()
	make_splines()
	return Composition{}
}
