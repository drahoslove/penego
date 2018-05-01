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

type tree struct {
	graph
	of *graph
}

// computes low and lim values for nodes
func (tr *tree) dfsRange(n *node, parent *edge, low int) int {
	lim := low
	n.parent = parent
	n.low = low
	for _, e := range n.outEdges(tr.edges) {
		if e != parent {
			lim = tr.dfsRange(e.to, e, lim)
		}
	}
	for _, e := range n.inEdges(tr.edges) {
		if e != parent {
			lim = tr.dfsRange(e.from, e, lim)
		}
	}
	n.lim = lim
	return lim + 1
}

func (tr *tree) dfsCutval(n *node, parent *edge) {
	for _, e := range n.outEdges(tr.edges) {
		if e != parent {
			tr.dfsCutval(e.to, e)
		}
	}
	for _, e := range n.inEdges(tr.edges) {
		if e != parent {
			tr.dfsCutval(e.from, e)
		}
	}
	if parent != nil {
		x_cutval(parent, tr)
	}
}

func initCutValues(g *graph, tr *tree) {
	// computes the cut values of the tree edges. For each tree edge,this is computed by marking the nodes as belonging to the head or tail component, and thenperforming the sum of the signed weights of all edges whose head and tail are in differentcomponents, the sign being negative for those edges going from the head to the tail component
	if len(tr.nodes) == 0 {
		log.Fatalln("No nodes to init cutvalues")
	}
	tr.dfsRange(tr.nodes[0], nil, 1)
	tr.dfsCutval(tr.nodes[0], nil)
}

// set cut value of f, assuming values of edges on one side were already set
func x_cutval(f *edge, tr *tree) {
	var g = tr.of
	var n *node
	var dir int
	if f.from.parent == f {
		n = f.from
		dir = 1
	} else {
		n = f.to
		dir = -1
	}

	x_val := func(e *edge, n *node, dir int) int {
		var (
			other_n *node
			flip    bool
			d       int
			rv      int
		)

		if e.from == n {
			other_n = e.to
		} else {
			other_n = e.from
		}
		if !(n.low <= other_n.lim && other_n.lim <= n.lim) {
			flip = true
			rv = e.weight()
		} else {
			flip = false
			if tr.includesEdge(e) {
				rv = e.cutValue
			} else {
				rv = 0
			}
			rv -= e.weight()
		}
		if dir > 0 && e.to == n || dir < 0 && e.from == n {
			d = 1
		} else {
			d = -1
		}
		if flip {
			d = -d
		}
		if d < 0 {
			rv = -rv
		}
		return rv
	}

	sum := 0
	for _, e := range n.outEdges(g.edges) {
		sum += x_val(e, n, dir)
	}
	for _, e := range n.inEdges(g.edges) {
		sum += x_val(e, n, dir)
	}
	f.cutValue = sum
}

func checkCutval(tr *tree) {
	for _, n := range tr.nodes {
		for _, e := range n.outEdges(tr.edges) {
			save := e.cutValue
			x_cutval(e, tr)
			if save != e.cutValue {
				log.Fatalln("Edge cutvalues not computed")
			}
		}
	}
}

func isTightTree(g *graph) bool {
	for _, e := range g.edges {
		if e.slack() > 0 {
			return false
		}
	}
	return true
}

func tightTree(g *graph) (*tree, int) {
	tr := &tree{newGraph(), g}
	size := 0
	var findSubtree func(v *node, subtree *tree) int

	findSubtree = func(v *node, subtree *tree) int {
		size := 1
		v.tree = subtree
		for _, e := range v.inEdges(g.edges) {
			if subtree.includesEdge(e) {
				continue
			}
			if e.from.tree == nil && e.slack() == 0 {
				subtree.addEdge(e)
				size += findSubtree(e.from, subtree)
			}
		}
		for _, e := range v.outEdges(g.edges) {
			if subtree.includesEdge(e) {
				continue
			}
			if e.to.tree == nil && e.slack() == 0 {
				subtree.addEdge(e)
				size += findSubtree(e.to, subtree)
			}
		}
		return size
	}
	for _, v := range g.nodes {
		v.tree = nil
	}
	if len(g.nodes) > 0 {
		size += findSubtree(g.nodes[0], tr)
	}
	return tr, size
}

func initRank(g *graph) {
	// An initial feasible ranking is computed. For brevity,init_rankis not given here.  Our versionkeeps nodes in a queue.  Nodes are placed in the queue when they have no unscanned in-edges.As nodes are taken off the queue, they are assigned the least rank that satisfies their in-edges, andtheir out-edges are marked as scanned. In the simplest case, whereδ =1 for all edges, thiscorresponds to viewing the graph as a poset and assigning the minimal elements to rank 0.  Thesenodes are removed from the poset and the new set of minimal elements are assigned rank 1, etc.
	isFeasible := true
	for _, n := range g.nodes {
		n.priority = 0
		for _, e := range n.inEdges(g.edges) {
			n.priority++
			e.cutValue = 0
			// e.treeIndex = -1
			if isFeasible && e.len() < min_len { // head-tail
				isFeasible = false
			}
		}
	}
	if isFeasible {
		return
	}

	queue := make(chan *node, len(g.nodes))
	cntr := 0 // counter of nodes passed by queue

	for _, n := range g.nodes {
		if n.priority == 0 {
			queue <- n
		}
	}

	for len(queue) > 0 {
		n := <-queue
		n.rank = 0
		cntr++
		for _, e := range n.inEdges(g.edges) {
			rank := e.from.rank + min_len
			if rank > n.rank {
				n.rank = rank
			}
		}
		for _, e := range n.outEdges(g.edges) {
			e.to.priority--
			if e.to.priority == 0 {
				queue <- e.to
			}
		}
	}
	close(queue)
	if cntr != len(g.nodes) {
		for _, n := range g.nodes {
			log.Printf("node %p has priority %v and rank %v", n, n.priority, n.rank)
			// if n.priority != 0 {
			// }
		}
		log.Fatalln("trouble creating feasible tree", cntr, len(g.nodes))
	}
}

// return smallest feasible tight tree with initialized cut values
func feasibleTree(g *graph) (tr *tree) {

	initRank(g)
	size := 0

	for {
		tr, size = tightTree(g)
		if size == len(g.nodes) {
			break
		}
		incid_e, incid_n := tr.minIncidentEdge(g.edges) // a non-tree edge incident on the tree with a minimal amount of slack
		if incid_e == nil {
			break
		}
		delta := incid_e.slack() // slack of an edge is the difference of its length and its minimum length.
		if incid_n == incid_e.from {
			delta = -delta
		}
		for _, v := range tr.nodes {
			v.rank += delta
		}
	}
	initCutValues(g, tr)
	return
}

// Returns tree edge with most negative cutValue to be replaced
func leaveEdge(tr *tree) *edge {
	var leave_e *edge
	for _, e := range tr.edges {
		if e.cutValue < 0 {
			if leave_e == nil {
				leave_e = e
			} else {
				if e.cutValue < leave_e.cutValue {
					leave_e = e
				}
			}
		}
	}
	return leave_e
}

// returns non-tree edge withc smallest slack
// which connects head component with tree component after breaking tree by l
func enterEdge(leave_e *edge, g *graph, tr *tree) *edge {
	var (
		enter_e         *edge
		n               *node
		outsearch       bool
		dfsEnterOutEdge func(*node)
		dfsEnterInEdge  func(*node)
	)

	if leave_e.from.lim < leave_e.to.lim {
		n = leave_e.from
		outsearch = false
	} else {
		n = leave_e.to
		outsearch = true
	}

	min_slack := maxInt
	low := n.low
	lim := n.lim

	dfsEnterOutEdge = func(n *node) {
		for _, e := range n.outEdges(g.edges) {
			if !tr.includesEdge(e) {
				if !(low <= e.to.lim && e.to.lim <= lim) {
					slack := e.slack()
					if slack < min_slack || enter_e == nil {
						enter_e = e
						min_slack = slack
					}
				}
			} else {
				if e.to.lim < n.lim {
					dfsEnterOutEdge(e.to)
				}
			}
		}
		for _, e := range n.inEdges(tr.edges) {
			if min_slack <= 0 {
				break
			}
			if e.from.lim < n.lim {
				dfsEnterOutEdge(e.from)
			}
		}
	}

	dfsEnterInEdge = func(n *node) {
		for _, e := range n.inEdges(g.edges) {
			if !tr.includesEdge(e) {
				if !(low <= e.from.lim && e.from.lim <= lim) {
					slack := e.slack()
					if slack < min_slack || enter_e == nil {
						enter_e = e
						min_slack = slack
					}
				}
			} else {
				if e.from.lim < n.lim {
					dfsEnterInEdge(e.from)
				}
			}
		}
		for _, e := range n.outEdges(tr.edges) {
			if min_slack <= 0 {
				break
			}
			if e.to.lim < n.lim {
				dfsEnterInEdge(e.to)
			}
		}
	}

	if outsearch {
		dfsEnterOutEdge(n)
	} else {
		dfsEnterInEdge(n)
	}
	return enter_e
}

// computes new ranks, cutvalues and swap leave and enter edges in tree
func update(leave_e, enter_e *edge, tr *tree) {
	delta := enter_e.slack()
	/* "for (v = in nodes in tail side of e) do ND_rank(v) -= delta;" */
	if delta > 0 {
		s := len(leave_e.from.inEdges(tr.edges)) + len(leave_e.from.outEdges(tr.edges))
		if s == 1 {
			rerank(leave_e.from, delta, tr)
		} else {
			s = len(leave_e.to.inEdges(tr.edges)) + len(leave_e.to.outEdges(tr.edges))
			if s == 1 {
				rerank(leave_e.to, -delta, tr)
			} else {
				if leave_e.from.lim < leave_e.to.lim {
					rerank(leave_e.from, delta, tr)
				} else {
					rerank(leave_e.to, -delta, tr)
				}
			}
		}
	}

	// set new cutvalues to all nodes between from_v and lca(from_n, to_n)
	treeUpdate := func(from_n, to_n *node, cutvalue int, dir int) *node {
		for !(from_n.low <= to_n.lim && to_n.lim <= from_n.lim) {
			e := from_n.parent
			if from_n != e.from {
				dir = -dir
			}
			e.cutValue += cutvalue * dir

			if e.from.lim > e.to.lim {
				from_n = e.from
			} else {
				from_n = e.to
			}
		}
		return from_n
	}

	cutvalue := leave_e.cutValue
	// Lowest common ancestor
	lca1 := treeUpdate(enter_e.from, enter_e.to, cutvalue, 1)
	lca2 := treeUpdate(enter_e.to, enter_e.from, cutvalue, -1)
	if lca1 != lca2 {
		log.Fatalln("mismatched lca in tree updates")
	}
	enter_e.cutValue = -cutvalue
	leave_e.cutValue = 0
	exchangeTreeEdges(leave_e, enter_e, tr)
	tr.dfsRange(lca1, lca1.parent, lca1.low)
}

func exchangeTreeEdges(leave_e, enter_e *edge, tr *tree) {
	for i, e := range tr.edges {
		if e == leave_e {
			tr.edges[i] = enter_e
		}
	}
}

func rerank(n *node, delta int, tr *tree) {
	n.rank -= delta
	for _, e := range n.outEdges(tr.edges) {
		if e != n.parent {
			rerank(e.to, delta, tr)
		}
	}
	for _, e := range n.inEdges(tr.edges) {
		if e != n.parent {
			rerank(e.from, delta, tr)
		}
	}
}

// Iterative method for graph drawing
// based on dot algorithm and work of Warfield, Sugiyamaet at al.
func GetIterative(network net.Net) Composition {
	graph := loadGraph(network)

	// assign rank λ(v) to each node v
	// edge e = (v, w)
	// lenght of edge l(e) ≥ δ(e)
	// l(e) = λ(w) − λ(v)

	// make graph acyclic by reversing edges
	graph = graph.acyclic()

	normalize := func() {

	}

	balance := func() {

	}

	rank := func() {
		tr := feasibleTree(&graph)
		// checkCutval(&graph, tr)
		for {
			e := leaveEdge(tr)
			if e == nil {
				break
			}
			f := enterEdge(e, &graph, tr)
			update(e, f, tr) // exchange e, f and recompute ranks and tuvals
		}
		normalize()
		balance()
	}

	ordering := func() {

	}

	position := func() {
		ranksToPos := map[int]int{}
		for _, n := range graph.nodes {
			if _, ok := ranksToPos[n.rank]; ok {
				ranksToPos[n.rank] = 0
			} else {
				ranksToPos[n.rank] += 1
			}
			n.position = ranksToPos[n.rank]
			log.Println("node rank and position", n.rank, n.position)
		}
	}

	makeSplines := func() {

	}

	rank()
	ordering()
	position()
	makeSplines()

	comp := New()

	for _, n := range graph.nodes {
		pos := draw.Pos{float64(n.position * 90), float64(n.rank) * 90}
		switch node := n.Composable.(type) {
		case *net.Transition:
			comp.transitions[node] = pos
		case *net.Place:
			comp.places[node] = pos
		}
	}

	return comp
}
