package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	typeName = flag.String("type", "", "type name; must be set")
)

const importStmt = `import (
	"reflect"

	"github.com/longsolong/flow/pkg/workflow/atom"
)
`

const atomIDMethod = `// AtomID ...
func (s *%[1]s) AtomID() atom.AtomID {
	return atom.AtomID{
		Type: reflect.TypeOf(s).Elem().String(),
		ID: s.ID,
		ExpansionDigest: s.ExpansionDigest,
	}
}
`

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of genatom:\n")
	fmt.Fprintf(os.Stderr, "\tgenatom -type T\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("genatom: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	g := Generator{}

	g.parsePackage()

	// Print the header and package clause.
	g.Printf("// Code generated by \"genatom %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s\n", g.pkg.name)
	g.Printf("\n")
	g.Printf(importStmt)
	g.Printf("\n")

	// Run generate for each type.
	g.Printf(atomIDMethod, *typeName)

	// Format the output.
	src := g.format()

	// Write to file.
	baseName := fmt.Sprintf("%s_atom.go", *typeName)
	outputName := filepath.Join(".", strings.ToLower(baseName))
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

type Package struct {
	name string
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage() {
	cfg := &packages.Config{
		Mode:  packages.NeedName,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name: pkg.Name,
	}
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf(string(g.buf.Bytes()))
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
