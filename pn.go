package main

import (
	"fmt"
	"os"
	"strings"

	"git.yo2.cz/drahoslav/penego/compose"
	"git.yo2.cz/drahoslav/penego/net"

	"git.yo2.cz/drahoslav/penego/storage"
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
		composition = compose.Parse(compoStr, network)
	} else {
		composition = Compose(network)
	}

	return
}

// returns composition based on settings
func Compose(network net.Net) (composition compose.Composition) {
	composer := storage.Of("settings").String("composer")
	if composer == "simple" {
		composition = compose.GetSimple(network)
	} else {
		composition = compose.GetIterative(network)
	}
	return
}

func splitBy(str string, delims []string) map[string]string {
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
