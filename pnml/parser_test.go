package pnml

import (
	"testing"
	"bytes"
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
		    />
		  </net>
		</pnml>
	`))
	resNet := Parse(pnml)

 	//  TODO make equal
	g := &net.Place{Tokens:1}
	e := &net.Place{Description: "exit"}
	t := &net.Transition{
		Origins: net.Arcs{{1, g}},
		Targets: net.Arcs{{1, g}, {2, e}},
	}
	refNet := net.New(
		net.Places{g, e},
		net.Transitions{t},
	)

	// TODO compare function
	if eq, err := resNet.Equals(&refNet); !eq {
		test.Errorf("Parser failed, because %s", err)
	}
}