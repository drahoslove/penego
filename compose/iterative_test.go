package compose

import (
	"testing"

	// "git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)

func getNet() net.Net {
	net, _ := net.Parse(`
g(1)
f(0)"fronta"
k(5)"kuchařky"
v(0)"výdej"
s(0)"stravování"
o(0)
z(0)
c(0)
i(0)"karanténa"
x(0)"restartů"
----
g -> t1[exp(3m)]"příchod studentů" -> g, f
f, k -> t2[] -> v
v -> t3[exp(1m)] -> s, k
s -> t4[10m..15m] -> o
t5[exp(240h)] -> z
z, g -> t6[p=1] -> c
c, f -> t7[p=3] -> c, o
c, v -> t8[p=2] -> c, o, k
c, s -> t9[p=1] -> c, o
c -> t10[] -> i
i -> t11[24h] -> g, x
	`)
	return net
}

func TestLoadGraph(test *testing.T) {
	graph := loadGraph(getNet())

	test.Logf("%v", graph)
}

// as image 5.13 in https://knihy.nic.cz/files/edice/pruvodce_labyrintem_algoritmu.pdf
func getGraphWithComponets() graph {

	c10 := &node{Composable: 10}
	c11 := &node{Composable: 11}

	c20 := &node{Composable: 20}

	c30 := &node{Composable: 30}
	c31 := &node{Composable: 31}
	c32 := &node{Composable: 32}

	c10c11 := &edge{from: c10, to: c11}
	c10c30 := &edge{from: c10, to: c30}
	c10c20 := &edge{from: c10, to: c20}
	c10c32 := &edge{from: c10, to: c32}

	c11c10 := &edge{from: c11, to: c10}

	c20c32 := &edge{from: c20, to: c32}

	c30c32 := &edge{from: c30, to: c32}

	c31c30 := &edge{from: c31, to: c30}

	c32c21 := &edge{from: c32, to: c31}

	return graph{
		[]*node{c10, c11, c20, c30, c31, c32},
		[]*edge{c10c11, c10c30, c10c20, c10c32, c11c10, c20c32, c30c32, c31c30, c32c21},
	}
}

func TestGraphTranspose(test *testing.T) {
	graph := getGraphWithComponets()
	graphT := graph.transpose()

	for _, e := range graph.edges {
		test.Logf("%p->%p", e.from, e.to)
	}
	test.Logf("")
	for _, e := range graphT.edges {
		test.Logf("%p->%p", e.from, e.to)
	}
}

func TestDFS(test *testing.T) {
	graph := getGraphWithComponets()
	v0 := graph.nodes[0]

	log := func(str string) func(*node) {
		return func(v *node) {
			test.Log(str, v.Composable)
		}
	}

	dfs(graph, v0, func(*node) bool { return true }, log("in"), log("out"))
	test.Log("")
	dfs(graph, v0, func(v *node) bool { return v.Composable.(int) <= 20 }, log("in"), log("out"))
}

func TestGraphComponents(test *testing.T) {
	graph := getGraphWithComponets()

	comps, nontriv := graph.components()
	test.Log("nontriv", nontriv)

	for i, v := range comps {
		test.Log(i.Composable, v.Composable)
	}

	if len(nontriv) != 2 {
		test.Error("graph should have 2 nontrivial strongly connected components")
	}

	/////

	acyclic := graph.acyclic()
	comps, nontriv = acyclic.components()
	test.Log("nontriv", nontriv)

	for i, v := range comps {
		test.Log(i.Composable, v.Composable)
	}

	if len(nontriv) != 0 {
		test.Error("acyclic graph should have no nontrivial strongly connected components")
	}

}

func TestInOutEdges(test *testing.T) {
	graph := loadGraph(getNet())

	for _, e := range graph.edges {
		test.Log(e, e.reversed())
	}
}

func TestTighTree(test *testing.T) {
	graph := loadGraph(getNet())

	tightTree, size := tightTree(&graph)
	if !isTightTree(&tightTree.graph) {
		test.Errorf("Tree %v is no tight", tightTree)
	}

	test.Log("graph size", len(graph.nodes))
	test.Log("tightTree size", size)
	test.Log("nodes", len(tightTree.nodes), "edeges", len(tightTree.edges))
	for _, e := range tightTree.edges {
		test.Log("edge", e)
	}
}

func TestFeasibleTree(test *testing.T) {
	graph := loadGraph(getNet()).acyclic()
	feasibleTree := feasibleTree(&graph)

	if !isTightTree(&feasibleTree.graph) {
		test.Errorf("Tree %v is not tight", feasibleTree)
	}

	if len(feasibleTree.nodes) != len(graph.nodes) {
		test.Errorf("fessible tre contains only %v nodes of original %v graphs nodes", len(feasibleTree.nodes), len(graph.nodes))
	}

	test.Logf("Feasible tree %v", feasibleTree)
}

func TestFillPathNodes(test *testing.T) {
	graph := loadGraph(getNet()).acyclic()
	rank(&graph)

	for _, e := range graph.edges {
		if e.from.rank-e.to.rank > 1 || e.to.rank-e.from.rank > 1 {
			test.Log("ranks apart", e.from.rank, e.to.rank)
		}
	}
	fillPathNodes(&graph)

	for _, e := range graph.edges {
		if e.from.rank-e.to.rank > 1 || e.to.rank-e.from.rank > 1 {
			test.Error("ranks apart", e.from.rank, e.to.rank)
		}
	}

}

func TestIterative(test *testing.T) {
	comp := GetIterative(getNet())

	test.Log(comp)
}
