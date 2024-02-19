//go:build !nodbxml
// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"
)

func doDact(item string) {
	db, err := dbxml.OpenRead(item)
	x(err)
	docs, err := db.All()
	x(err)
	for docs.Next() {
		filenames = append(filenames, item+"::"+docs.Name())
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
