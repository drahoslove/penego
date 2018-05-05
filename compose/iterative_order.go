package compose

import (
	"sort"
)

const MAX_ORDERING_ITER = 32

// returns ordered orders of all incident nodes in given rank
func (n *node) adjOrders(g *graph, order order) []int {
	orders := []int{}
	for _, e := range n.inEdges(g.edges) {
		for i, rank_n := range order {
			if e.from == rank_n {
				orders = append(orders, i)
			}
		}
	}
	for _, e := range n.outEdges(g.edges) {
		for i, rank_n := range order {
			if e.from == rank_n {
				orders = append(orders, i)
			}
		}
	}

	sort.Ints(orders)
	return orders
}

type order []*node
type orders map[int]order

func (o order) Len() int {
	return len(o)
}
func (o order) Less(i, j int) bool {
	return o[i].weight < o[j].weight
}

func (o order) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (os orders) clone() orders {
	clone := make(orders, len(os))
	for i, _ := range os {
		clone[i] = make(order, len(os[i]))
		copy(clone[i], os[i])
	}
	return clone
}

// Orders nodes withing same rank using median wighting based algorithm described here:
// http://www.graphviz.org/Documentation/TSE93.pdf
func ordering(g *graph) {
	fillPathNodes(g)

	orders := initOrder(g)
	bestOrders := orders.clone()
	for i := 0; i < MAX_ORDERING_ITER; i++ {
		weightMedian(g, orders, i)
		transpose(g, orders)
		if crossings(g, orders) < crossings(g, bestOrders) {
			bestOrders = orders.clone()
		}
	}

	// apply best orders to nodes
	for _, nodes := range bestOrders {
		for i, n := range nodes {
			n.order = i
		}
	}
}

// re-assign weight value to each node
func weightMedian(g *graph, orders orders, iter int) {
	medianValue := func(n *node, rank int) float64 {
		P := n.adjOrders(g, orders[rank])
		mid := len(P) / 2
		if len(P) == 0 {
			return -1.0
		}
		if len(P)%2 == 1 {
			return float64(P[mid])
		}
		if len(P) == 2 {
			return float64(P[0]+P[1]) / 2
		} else {
			left := float64(P[mid-1] - P[0])
			right := float64(P[len(P)-1] - P[mid])
			return (float64(P[mid-1])*right + float64(P[mid])*left) / (left + right)
		}
	}

	if iter%2 == 0 {
		for rank, nodes := range orders {
			for _, n := range nodes {
				n.weight = medianValue(n, rank-1)
			}
			sort.Sort(nodes)
		}
	} else {
		for rank, nodes := range orders {
			defer func() {
				for _, n := range nodes {
					n.weight = medianValue(n, rank+1)
				}
				sort.Sort(nodes)
			}()
		}
	}
}

func transpose(g *graph, orders orders) {
	improved := true
	for improved {
		improved = false
		for rank, nodes := range orders {
			crossing := func(nodes order) int {
				if rank+1 == len(orders) { // last rank
					return 0
				}
				count := countCrossings(g, nodes, orders[rank+1])
				return count
			}
			for i := 0; i < len(nodes)-2; i++ {
				n1 := nodes[i]
				n2 := nodes[i+1]
				exchNodes := make(order, len(nodes))
				copy(exchNodes, nodes)
				for i := 0; i < len(exchNodes)-2; i++ {
					if exchNodes[i] == n1 && exchNodes[i+1] == n2 {
						exchNodes[i], exchNodes[i+1] = n2, n1
					}
				}
				if crossing(exchNodes) < crossing(nodes) {
					improved = true
					orders[rank] = exchNodes // swap
				}
			}
		}
	}
}

// counts crossings between two ranks
// algorithm by Barth, Mutzel and Junger 
func countCrossings(g *graph, northRank, southRank order) int {
	firstIndex := 1
	for firstIndex < len(southRank) { // 10
		firstIndex *= 2 // 1 2 4 8 16
	}
	treeSize := 2*firstIndex - 1 // 31
	tree := make([]int, treeSize)
	firstIndex -= 1 // 15 // index of leftmost leaf

	// cache orders of south nodes
	orderOfSn := make(map[*node]int, len(southRank))
	for i, n := range southRank {
		orderOfSn[n] = i
	}

	count := 0
	for _, nn := range northRank {
		for _, e := range nn.outEdges(g.edges) {
			if i, ok := orderOfSn[e.to]; ok {
				index := i + firstIndex
				for index > 0 {
					if index%2 == 1 {
						count += tree[index+1]
					}
					index = (index - 1) / 2
					tree[index]++
				}
			}
		}
		for _, e := range nn.inEdges(g.edges) {
			if i, ok := orderOfSn[e.from]; ok {
				index := i + firstIndex
				for index > 0 {
					if index%2 == 1 {
						count += tree[index+1]
					}
					index = (index - 1) / 2
					tree[index]++
				}
			}
		}
	}

	return count
}

// counts all crossings in graph
func crossings(g *graph, orders orders) int {
	count := 0
	for rank := 0; rank < len(orders)-2; rank++ {
		count += countCrossings(g, orders[rank], orders[rank+1])
	}
	return count
}

// create inital ordering using dfs
func initOrder(g *graph) orders {
	ordersByRank := orders{}
	visited := map[*node]bool{}

	notVisited := func(n *node) bool {
		return !visited[n]
	}
	asignOrder := func(n *node) {
		if _, ok := ordersByRank[n.rank]; !ok {
			ordersByRank[n.rank] = []*node{}
		}
		ordersByRank[n.rank] = append(ordersByRank[n.rank], n)
		visited[n] = true
	}
	for _, n := range g.nodes {
		dfs(*g, n, notVisited, asignOrder, nil)
	}
	return ordersByRank
}

// adds intermediate "path nodes" to edges whose ranks are not adjencent
func fillPathNodes(g *graph) {
	for _, e := range g.edges {
		path := &path{e.from.Composable, e.to.Composable}
		addNodeEdge := func(rank int) {
			inter_n := &node{Composable: path, rank: rank}
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
