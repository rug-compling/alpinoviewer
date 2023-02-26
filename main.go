package main

import (
	"github.com/pebbe/util"

	"bytes"
	"flag"
	"os"
)

var (
	x    = util.CheckErr
	optN = flag.String("n", "", "gemarkeerde nodes in boom")
	// optU = flag.String("u", "", "gemarkeerde nodes in UD")
	// optE = flag.String("e", "", "gemarkeerde nodes in extended UD")
)

func main() {
	flag.Parse()

	b, err := os.ReadFile(flag.Arg(0))
	x(err)

	var buf bytes.Buffer

	tree(b, &buf)

	run(buf.String(), flag.Arg(0))

}
