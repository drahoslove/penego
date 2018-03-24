package pnml

import (
	"bytes"
	"testing"

	"git.yo2.cz/drahoslav/penego/net"
)

func TestParse(test *testing.T) {
	pnml := bytes.NewReader([]byte(`
		<pnml>
		  <net>
		    <place id="p1">
		      <initialMarking>
		        <text>3</text>
		      </initialMarking>
		    </place>
		    <place id="p2" />
		    <transition id="t1" />
		    <arc id="a1" source="p1" target="t1" />
		    <arc id="a2" source="t1" target="p1" />
		    <arc id="a2" source="t1" target="p2">
		      <inscription>
		        <text>2</text>
		      </inscription>
		    </arc>
		  </net>
		</pnml>
	`))
	resNet := Parse(pnml)

	//  TODO make equal
	p1 := &net.Place{Id: "p1", Tokens: 3}
	p2 := &net.Place{Id: "p2"}
	t := &net.Transition{
		Id:      "t1",
		Origins: net.Arcs{{1, p1}},
		Targets: net.Arcs{{1, p1}, {2, p2}},
	}
	refNet := net.New(
		net.Places{p1, p2},
		net.Transitions{t},
	)

	// TODO compare function
	if eq, err := resNet.Equals(&refNet); !eq {
		test.Errorf("Parser failed, because %s \n%s\nshould be\n%s\n", err, refNet, resNet)
	}
}
