package main

import (
	"bytes"
	"os"

	"github.com/pebbe/util"
)

var (
	x = util.CheckErr
)

func main() {
	b, err := os.ReadFile(os.Args[1])
	x(err)

	var buf bytes.Buffer

	tree(b, &buf)

	run(buf.String(), os.Args[1])

}
