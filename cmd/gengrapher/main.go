package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeName = flag.String("type", "", "type name; must be set")
)

const importStmt = `import (
	"context"

	"github.com/faceair/jio"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/standalone/graph"
)
`

const newGrapherMethod = `// NewGrapher ...
func NewGrapher(ctx context.Context, rawRequestData []byte) (*graph.Grapher, error) {
	req, err := newRequest(ctx, rawRequestData)
	if err != nil {
		return nil, err
	}
	p, err := newPlotter(ctx, req)
	if err != nil {
		return nil, err
	}
	g := graph.NewGrapher(req, p.DAG, p.Chain, p)
	return g, nil
}

func newRequest(ctx context.Context, rawRequestData []byte) (*request.Request, error) {
	requestArgs, err := jio.ValidateJSON(&rawRequestData, schema)
	if err != nil {
		return nil, err
	}
	req := request.NewRequestWithContext(ctx)
	req.RequestArgs = requestArgs["requestArgs"].(map[string]interface{})
	for _, v := range requestArgs["requestTags"].([]interface{}) {
		v := v.(map[string]interface{})
		req.RequestTags = append(req.RequestTags, request.Tag{Name: v["name"].(string), Value: v["value"].(string)})
	}
	return req, nil
}

// newPlotter ...
func newPlotter(ctx context.Context, req *request.Request) (*plotter, error) {
	p := &plotter{Plotter: graph.NewPlotter(NAME, VERSION)}
	if err := p.Begin(ctx, req); err != nil {
		return nil, err
	}
	return p, nil
}
`

const fileSuffix string = "_grapher.go"

func sourceFilter(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasSuffix(fi.Name(), fileSuffix)
}


// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of gengrapher:\n")
	fmt.Fprintf(os.Stderr, "\tgengrapher\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("gengrapher: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	g := Generator{}

	g.parsePackage()

	// Print the header and package clause.
	g.Printf("// Code generated by \"gengrapher %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s\n", g.pkg.name)
	g.Printf("\n")
	g.Printf(importStmt)
	g.Printf("\n")

	// Run generate.
	g.Printf(newGrapherMethod)

	// Format the output.
	src := g.format()

	// Write to file.
	baseName := fmt.Sprintf("%s%s", *typeName, fileSuffix)
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
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", sourceFilter, 0)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs)
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkgs map[string]*ast.Package) {
	for pkgName := range pkgs {
		g.pkg = &Package{
			name: pkgName,
		}
		break
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