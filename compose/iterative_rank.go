/*
 * Network Simplex Algorithm for Ranking Nodes of a directed acyclic graph (DAG)
 * this is strongly based on algorithm described in here https://graphviz.gitlab.io/_pages/Documentation/TSE93.pdf
 */
package compose

import (
	"log"
)

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

func rank(g *graph) {
	balance := func(tr *tree) {
		for _, e := range tr.edges {
			if e.cutValue == 0 {
				enter_e := enterEdge(e, tr.of, tr)
				if enter_e == nil {
					continue
				}
				delta := enter_e.slack()
				if delta <= 1 {
					continue
				}
				if e.from.lim < e.to.lim {
					rerank(e.from, delta/2, tr)
				} else {
					rerank(e.to, -delta/2, tr)
				}
			}
		}
	}
	tr := feasibleTree(g)
	// checkCutval(&graph, tr)
	for {
		e := leaveEdge(tr)
		if e == nil {
			break
		}
		f := enterEdge(e, g, tr)
		update(e, f, tr) // exchange e, f and recompute ranks and cutvals
		// TODO add max iteration limit
	}
	tr.normalizeRanks()
	balance(tr)
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

func checkRanks(g *graph) {
	cost := 0
	for _, n := range g.nodes {
		for _, e := range n.outEdges(g.edges) {
			l := e.len()
			if l < 0 {
				l = -l
			}
			cost += e.weight() * l
			if e.to.rank-e.from.rank-min_len < 0 {
				log.Fatalln("ranks of nodes of edge are wrong")
			}
		}
	}
	log.Println("rank cost is", cost)
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
