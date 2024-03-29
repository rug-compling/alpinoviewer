package main

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.0
#cgo CFLAGS: -DGDK_DISABLE_DEPRECATED -D_doesnt_work_with_webkit_GTK_DISABLE_DEPRECATED
#include <stdlib.h>
#include "view_my.h"
*/
import "C"

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"unsafe"
)

type msg struct {
	id int
	ms string
}

var (
	chGtkDone   = make(chan bool)
	chGoDone    = make(chan bool)
	chMessage   = make(chan msg, 100)
	tmpfilename string
)

//export go_message
func go_message(id int, cContent *C.char) {
	content := C.GoString(cContent)
	chMessage <- msg{id: id, ms: content}
}

func run(content, title string, filenames []string) {
	C.setnfiles(C.int(len(filenames)))
	for _, filename := range filenames {
		cs := C.CString(filename)
		C.addfile(cs)
	}

	tmpfile, err := ioutil.TempFile("/tmp", "alpinoview*.html")
	x(err)
	tmpfilename = tmpfile.Name()
	defer os.Remove(tmpfilename)
	_, err = tmpfile.Write([]byte(content))
	x(err)
	x(tmpfile.Close())

	cs := C.CString("file://" + tmpfile.Name())
	defer C.free(unsafe.Pointer(cs))

	ct := C.CString(title)
	defer C.free(unsafe.Pointer(ct))

	go doStuff()

	runtime.LockOSThread()

	C.run(cs, ct)
	// log.Println("Gtk done")
	close(chGtkDone)

	<-chGoDone
	// log.Println("All done")
}

func doStuff() {
	defer close(chGoDone)

	// Wait for Gtk to get ready
LOOP:
	for {
		select {
		case <-chGtkDone:
			return
		case m := <-chMessage:
			doMessage(m)
			if m.id == C.idREADY {
				break LOOP
			}
		}
	}

	for {
		select {
		case <-chGtkDone:
			return
		case m := <-chMessage:
			doMessage(m)
		}
	}
}

func doMessage(m msg) {
	switch m.id {
	default:
		log.Printf("-- unknown message id %d: %q\n", m.id, m.ms)
	case C.idERROR:
		log.Printf("-- error: %s\n", m.ms)
	case C.idREADY:
		//		log.Printf("-- ready: %s\n", m.ms)
	case C.idDELETE:
		//		log.Printf("-- delete event: %s\n", m.ms)
	case C.idDESTROY:
		//		log.Printf("-- destroy event: %s\n", m.ms)
	case C.idSELECT:
		doFile(m.ms)
	}
}

func doFile(filename string) {
	b, err := getFile(filename)
	x(err)

	var buf bytes.Buffer

	tree(b, &buf, filename)

	fp, err := os.Create(tmpfilename)
	x(err)
	fp.WriteString(buf.String())
	fp.Close()

	ct := C.CString(filename)
	defer C.free(unsafe.Pointer(ct))
	C.reload(ct)
}
