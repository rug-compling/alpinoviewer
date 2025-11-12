//go:build nodbxml
// +build nodbxml

package main

import (
	"fmt"
	"os"
)

const optDact = ""

func doDact(s string) {
	fmt.Fprintln(os.Stderr, "Support for Dact/DbXML not compiled in")
	os.Exit(1)
}

func getDact(corpus, file string) ([]byte, error) {
	return []byte{}, fmt.Errorf("Support for Dact/DbXML not compiled in")
}
