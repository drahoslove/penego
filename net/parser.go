package net

const FORMAT_EXAMPLE = `
// definice míst:

g (1) // generování studentů
f (0) "fronta"
k (4) "kuchařky"
v ( ) "výdej"
s ( ) "stravování"
o ( ) // odchod

z ( ) // ohlášena žloutenka
c ( ) // vyprazdňovací cyklus
k ( ) "karanténa"

// definice přechodů:

g	-> [exp(30s)] "příchod studentů" -> g,f
f,k	-> [] -> v
v	-> [exp(1m)] -> s,k
s 	-> [10m-15m] -> o
o	-> [] "odchod"

[exp(100d)]  -> z
z,g -> [p=1] -> v
c,f	-> [p=3] -> c,o
c,v	-> [p=2] -> c,o,k
c,s	-> [p=1] -> c,o
c	-> [p=0] -> k
k	-> [10d] -> x
`

func Parse(input string) (Transitions, Places, error) {
	return Transitions{}, Places{}, nil
}

