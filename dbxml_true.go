//go:build !nodbxml
// +build !nodbxml

package main

import (
	"strings"

	"github.com/pebbe/dbxml"
)

const optDact = " *.dact,"

func doDact(item string) {
	db, err := dbxml.OpenRead(item)
	x(err)
	docs, err := db.All()
	x(err)
	for docs.Next() {
		name := docs.Name()
		if strings.HasSuffix(name, ".xml") {
			filenames = append(filenames, item+"::"+name)
		}
	}
	x(docs.Error())
	db.Close()
}

func getDact(corpus, file string) ([]byte, error) {
	db, err := dbxml.OpenRead(corpus)
	if err != nil {
		return []byte{}, err
	}
	s, err := db.Get(file)
	db.Close()
	if err != nil {
		return []byte{}, err
	}
	return []byte(s), nil
}
