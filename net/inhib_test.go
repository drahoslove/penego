package net


import (
	"testing"
)

func TestInhibitorArc(test *testing.T) {
	g := &Place{Id: "g", Tokens: 1}
	e := &Place{Id: "e", Description: "exit"}
	t := &Transition{
		Id: "t",
		Origins: Arcs{{1, InhibitorArc, g}},
		Targets: Arcs{{1, NormalArc, g}, {2, NormalArc, e}},
	}

	netA := New(
		Places{g, e},
		Transitions{t},
	)

	test.Log(netA)

	expected := "!g -> t[] -> g, 2*e"

	if t.String() != expected {
		test.Error("Transition should be stringified as\n", expected)
		test.Error("but it is\n",  t)
	}
}

func TestInhibitorArcParse(test *testing.T) {
	netW, err := Parse(`
		g()
		e(0)"exit"
		----
		!g -> t[] -> g, 2*e
	`)

	test.Log(netW)

	if err != nil {
		test.Error("Error while parsing net with inhibitor edge\n", err)
	}
	if netW.transitions[0].Origins[0].Type != InhibitorArc {
		test.Error("Arc should be inhibitor")
	}
}

func TestInhibitorTarget(test *testing.T) {
	netW, err := Parse(`
		g()
		e(0)"exit"
		----
		g -> t[] -> !g, 2*e
	`)

	test.Log(netW)

	if err == nil {
		test.Error("Inhibitor in outgoing arcs is nonsense and it should not parse", netW)
	}
}