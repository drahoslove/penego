package net

import (
	"testing"
)

func TestNetEquals(test *testing.T) {

	g := &Place{Tokens:1}
	e := &Place{Description: "exit"}
	t := &Transition{
		Origins: Arcs{{1, g}},
		Targets: Arcs{{1, g}, {2, e}},
	}
	refNet := New(
		Places{g, e},
		Transitions{t},
	)


	if ok, err := refNet.Equals(&refNet); !ok {
		test.Errorf("Nets are not equal, %s", err)
	}

}