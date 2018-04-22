package net

import (
	"testing"
)

func TestNetString(test *testing.T) {

	te := &Transition{Id:"te"}

	origins := Arcs{}
	targets := Arcs{}

	p1 := &Place{Id: "p1"}
	p2 := &Place{Id: "p2"}
	origins.Push(1, p1)
	targets.Push(1, p2)

	tl := &Transition{Id:"tl", Origins:origins}
	tr := &Transition{Id:"tr", Targets:targets}

	n := Net{Places{p1, p2}, Transitions{te, tl, tr}}

	netStr := n.String()

	test.Log(netStr)

	n, err := Parse(netStr)
	if err != nil {
		test.Errorf("Stringified net should be parsable %v", err)
	}
}

func TestNetTransEmpty(test *testing.T) {
	n, _ := Parse(`
		t[]
	`)

	if !n.transitions[0].Origins.IsEmpty() {
		test.Error("Origins of tran should be empty")
	}

	if !n.transitions[0].Targets.IsEmpty() {
		test.Error("Targets of tran should be empty")
	}

}