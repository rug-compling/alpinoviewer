package main

// TODO: UD gebruiken uit xml-bestand, indien aanwezig
// TODO: optie: save svg

import (
	"archive/zip"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pebbe/compactcorpus"
	"github.com/pebbe/util"
)

var (
	reNumber  = regexp.MustCompile("[0-9]+")
	IDs       = make(map[string][]int)
	filenames = make([]string, 0)
	seen      = make(map[string]bool)
	stdin     []byte
	x         = util.CheckErr
	optI      = flag.Bool("i", false, "filenames and id's from stdin")
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

    -i  : bestandsnamen en id's via stdin, één per regel
          bestandsnaam gevolgd door tab, gevolgd door id's gescheiden door spaties

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

	if *optI {
		filenames = make([]string, 0)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			doItem(scanner.Text())
			// TODO: open display met eerste boom zodra eerst regel is ingelezen
		}
		x(scanner.Err())
	} else if flag.NArg() > 0 {
		for _, name := range flag.Args() {
			doItem(name)
		}
	} else {
		filenames = []string{"<stdin>"}
	}

	b, err = getFile(filenames[0])
	x(err)

	var buf bytes.Buffer

	tree(b, &buf, filenames[0])

	run(buf.String(), filenames[0], filenames)
}

func doItem(item string) {
	i := strings.IndexAny(item, "\t|")
	idlist := ""
	if i > 0 {
		idlist = item[i+1:]
		item = item[:i]
	}

	item = strings.TrimSpace(item)
	if len(item) == 0 {
		return
	}

	ids := make([]int, 0)
	for _, num := range reNumber.FindAllString(idlist, -1) {
		id, err := strconv.Atoi(num)
		x(err)
		ids = append(ids, id)
	}

	if strings.HasSuffix(item, ".xml") {
		filenames = append(filenames, item)
		IDs[item] = ids
		return
	}

	if strings.HasSuffix(item, ".dact") || strings.HasSuffix(item, ".dbxml") {
		doDact(item)
		return
	}

	if strings.HasSuffix(item, ".data.dz") || strings.HasSuffix(item, ".index") {
		doCompact(item)
		return
	}

	if strings.HasSuffix(item, ".zip") {
		doZip(item)
		return
	}

	fi, err := os.Stat(item)
	x(err)
	if fi.IsDir() {
		doDir(item)
	}
}

func doCompact(item string) {
	var s string
	if strings.HasSuffix(item, ".data.dz") {
		s = item[:len(item)-8]
	} else {
		s = item[:len(item)-6]
	}
	if seen[s] {
		return
	}
	seen[s] = true

	cc, err := compactcorpus.Open(item)
	x(err)
	r, err := cc.NewRange()
	x(err)
	for r.HasNext() {
		name, _ := r.Next()
		if strings.HasSuffix(name, ".xml") {
			filenames = append(filenames, item+"::"+name)
		}
	}
}

func doZip(item string) {
	zr, err := zip.OpenReader(item)
	x(err)
	for _, file := range zr.File {
		if file.FileInfo().IsDir() {
			continue
		}
		name := file.Name
		if strings.HasSuffix(name, ".xml") {
			filenames = append(filenames, item+"::"+name)
		}
	}
}

func doDir(item string) {
	files, err := ioutil.ReadDir(item)
	x(err)
	for _, file := range files {
		filename := file.Name()
		if filename != "." && filename != ".." {
			doItem(filepath.Join(item, filename))
		}
	}
}

func getFile(name string) ([]byte, error) {
	if name == "<stdin>" {
		if stdin == nil {
			var err error
			stdin, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return []byte{}, err
			}
		}
		return stdin, nil
	}

	aa := strings.SplitN(name, "::", 2)
	if len(aa) == 1 {
		return ioutil.ReadFile(name)
	}

	if strings.HasSuffix(aa[0], ".dact") || strings.HasSuffix(aa[0], ".dbxml") {
		return getDact(aa[0], aa[1])
	}

	if strings.HasSuffix(aa[0], ".data.dz") || strings.HasSuffix(aa[0], ".index") {
		return getCompact(aa[0], aa[1])
	}

	if strings.HasSuffix(aa[0], ".zip") {
		return getZip(aa[0], aa[1])
	}

	return []byte{}, fmt.Errorf("Unknown corpus type for %s", aa[0])
}

func getCompact(corpus, name string) ([]byte, error) {
	cc, err := compactcorpus.RaOpen(corpus)
	if err != nil {
		return []byte{}, err
	}
	b, err := cc.Get(name)
	cc.Close()
	return b, err
}

func getZip(corpus, name string) ([]byte, error) {
	zr, err := zip.OpenReader(corpus)
	if err != nil {
		return []byte{}, err
	}
	fp, err := zr.Open(name)
	if err != nil {
		return []byte{}, err
	}
	data, err := io.ReadAll(fp)
	fp.Close()
	return data, err
}
