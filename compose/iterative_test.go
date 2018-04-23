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

func TestIterativeComposition(test *testing.T) {
	graph := loadGraph(getNet())

	test.Logf("%v", graph)
}

// as image 5.13 in https://knihy.nic.cz/files/edice/pruvodce_labyrintem_algoritmu.pdf
func getGraphWithComponets() graph {

	c10 := &node{10, 1, 0}
	c11 := &node{11, 1, 1}

	c20 := &node{20, 2, 0}

	c30 := &node{30, 3, 0}
	c31 := &node{31, 3, 1}
	c32 := &node{32, 3, 2}

	c10c11 := &edge{c10, c11, 1}
	c10c30 := &edge{c10, c30, 1}
	c10c20 := &edge{c10, c20, 1}
	c10c32 := &edge{c10, c32, 1}

	c11c10 := &edge{c11, c10, 1}
	
	c20c32 := &edge{c20, c32, 1}

	c30c32 := &edge{c30, c32, 1}

	c31c30 := &edge{c31, c30, 1}

	c32c21 := &edge{c32, c31, 1}

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

	log := func(str string) (func (*node)) {
		return func (v *node) {
			test.Log(str, v.Composable)
		}
	}

	dfs(graph, v0, func(*node) bool {return true}, log("in"), log("out"))
	test.Log("")
	dfs(graph, v0, func(v *node) bool {return v.Composable.(int) <= 20}, log("in"), log("out"))
}

func TestGraphComponents(test *testing.T) {
	graph := getGraphWithComponets()

	components, n := graph.components()
	test.Log("count of comp", n)

	for i, v := range components {
		test.Log(i.Composable, v.Composable)
	}

}