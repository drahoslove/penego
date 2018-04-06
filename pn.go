package main

import (
	"strings"
	"os"
	"fmt"
	"git.yo2.cz/drahoslav/penego/compose"
	"git.yo2.cz/drahoslav/penego/net"
)

const (
	netDelim  = "# NET"
	compDelim = "# COMPOSITION"
)

func Stringify(network net.Net, composition compose.Composition) string {
	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s", netDelim, network, compDelim, composition)
}

func Parse(str string) (network net.Net, composition compose.Composition) {

	parts := splitBy(str, []string{netDelim, compDelim})
	netStr := parts[netDelim]
	if netStr == "" {
		netStr = parts[""]
	}

	network, err := net.Parse(netStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	compoStr := parts[compDelim]
	if compoStr != "" {
		// TODO composition = compose.Parse(compoStr)
		composition = compose.GetSimple(network)
	} else {
		composition = compose.GetSimple(network)
	}

	return
}

func splitBy(str string, delims []string) (map[string]string) {
	sections := map[string]string{"": str}

	for _, delim := range delims {
		for i, section := range sections {
			subparts := strings.Split(section, delim)
			sections[i] = subparts[0]
			if len(subparts) > 1 {
				sections[delim] = subparts[1]
			}
		}
	}
	return sections
}