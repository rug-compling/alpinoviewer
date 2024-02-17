package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pebbe/util"
)

var (
	IDs  = make(map[string][]int)
	x    = util.CheckErr
	optN = flag.String("n", "", "gemarkeerde nodes in boom")
	optI = flag.Bool("i", false, "filenames and id's from stdin")
	// optU = flag.String("u", "", "gemarkeerde nodes in UD")
	// optE = flag.String("e", "", "gemarkeerde nodes in extended UD")
)

func usage() {
	fmt.Printf(`
Syntax:

    %s [opties] file.xml
    %s [opties] < file.xml
    find . -name '*.xml' | %s -i

Opties:

    -i             : bestandsnamen en id's via stdin, één per regel
                     bestandsnaam gevolgd door tab, gevolgd door id's gescheiden door spaties
    -n ID1,ID2,... : IDs van nodes voor markering, gescheiden door komma

Gebruik:

    Ctrl - : zoom uit
    Ctrl = : zoom in
    Ctrl 0 : reset zoom
    Ctrl Q : exit

`, os.Args[0], os.Args[0], os.Args[0])
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

	var filenames []string

	if *optI {
		filenames = make([]string, 0)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			aa := strings.SplitN(scanner.Text(), "\t", 2)
			aa[0] = strings.TrimSpace(aa[0])
			if aa[0] != "" {
				filenames = append(filenames, aa[0])
				if len(aa) == 2 {
					ii := strings.Fields(aa[1])
					ids := make([]int, len(ii))
					for i, v := range ii {
						ids[i], err = strconv.Atoi(v)
						x(err)
					}
					IDs[aa[0]] = ids
				}
			}
		}
		x(scanner.Err())
		b, err = os.ReadFile(filenames[0])
		x(err)
	} else if flag.NArg() > 0 {
		filenames = flag.Args()
		b, err = os.ReadFile(filenames[0])
		x(err)
	} else {
		filenames = []string{""}
		b, err = io.ReadAll(os.Stdin)
		x(err)
	}

	if *optN != "" {
		aa := strings.Split(*optN, ",")
		id1 := make([]int, len(aa))
		for i, a := range aa {
			id1[i], err = strconv.Atoi(strings.TrimSpace(a))
			x(err)
		}
		IDs[filenames[0]] = id1
	}

	var buf bytes.Buffer

	tree(b, &buf, filenames[0])

	run(buf.String(), filenames[0], filenames)
}
