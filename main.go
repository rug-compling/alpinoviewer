package main

import (
	"github.com/pebbe/util"

	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	x    = util.CheckErr
	optN = flag.String("n", "", "gemarkeerde nodes in boom")
	// optU = flag.String("u", "", "gemarkeerde nodes in UD")
	// optE = flag.String("e", "", "gemarkeerde nodes in extended UD")
)

func usage() {
	fmt.Printf(`
Syntax:

    %s [opties] file.xml
    %s [opties] < file.xml

Opties:

    -n ID1,ID2,... : IDs van nodes voor markering, gescheiden door komma

Gebruik:

    Ctrl - : zoom uit
    Ctrl = : zoom in
    Ctrl 0 : reset zoom
    Ctrl Q : exit

`, os.Args[0], os.Args[0])
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var b []byte
	var err error

	if flag.NArg() == 0 && util.IsTerminal(os.Stdin) {
		usage()
		return
	}

	if flag.NArg() > 0 {
		b, err = os.ReadFile(flag.Arg(0))
		x(err)
	} else {
		b, err = io.ReadAll(os.Stdin)
		x(err)
	}

	var buf bytes.Buffer

	tree(b, &buf)

	run(buf.String(), flag.Arg(0), flag.Args())

}
