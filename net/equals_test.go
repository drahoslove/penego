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
	v := &Transition{
		Origins: Arcs{{1, g}},
		Targets: Arcs{{2, e}, {1, g}},
	}
	u := &Transition{
		Origins: Arcs{{1, g}},
		Targets: Arcs{{1, e}, {2, g}},
	}
	netA := New(
		Places{g, e},
		Transitions{t},
	)

	netB := New(
		Places{e, g},
		Transitions{t},
	)

	netC := New(
		Places{e, g},
		Transitions{v},
	)

	netD := New(
		Places{e, g},
		Transitions{u},
	)


	if eq, err := netA.Equals(&netB); !eq {
		test.Errorf("Nets shoud be equal but, %s", err)
	}

	if eq, err := netA.Equals(&netC); !eq {
		test.Errorf("Nets should also be equal but, %s", err)
	}

	if eq, _ := netA.Equals(&netD); eq {
		test.Errorf("Nets should no equal but they are")
	}

}