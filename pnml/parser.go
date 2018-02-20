package pnml

import (
	"encoding/xml"
	"io"
	"git.yo2.cz/drahoslav/penego/net"
)

func Parse(pnml io.Reader) net.Net {
	decoder := xml.NewDecoder(pnml)
	// start := xml.StartElement{xml.Name{"pnml"}, []xml.Attr{}}
	decoder.Decode(nil)
	return net.Net{}
}