// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// genaccessor generates accessor methods for structs fields.

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
	"sort"
	"strings"
	"text/template"
)

const (
	fileSuffix = "_accessor.go"
)

var (
	verbose = flag.Bool("v", false, "Print verbose log messages")
	typeNames   = flag.String("type", "", "comma-separated list of type names; must be set")

	sourceTmpl = template.Must(template.New("source").Parse(source))
)

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}

func main() {
	flag.Parse()
	fset := token.NewFileSet()

	types := strings.Split(*typeNames, ",")
	pkgs, err := parser.ParseDir(fset, ".", sourceFilter, 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	for pkgName, pkg := range pkgs {
		t := &templateData{
			filename: pkgName + fileSuffix,
			Package:  pkgName,
			Imports:  map[string]string{},
		}
		for filename, f := range pkg.Files {
			logf("Processing %v...", filename)
			if err := t.processAST(f, types); err != nil {
				log.Fatal(err)
			}
		}
		if err := t.dump(); err != nil {
			log.Fatal(err)
		}
	}
	logf("Done.")
}

func (t *templateData) processAST(f *ast.File, types []string) error {
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			ok = false
			for _, typeName := range(types) {
				if ts.Name.String() == typeName {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range st.Fields.List {
				if len(field.Names) == 0 || !ok {
					continue
				}
				fieldName := field.Names[0]
				// Skip unexported identifiers.
				if !fieldName.IsExported() {
					logf("Field %v is unexported; skipping.", fieldName)
					continue
				}

				se, ok := field.Type.(*ast.StarExpr)
				if !ok {
					switch x := field.Type.(type) {
					case *ast.ArrayType:
						t.addArrayType(x, ts.Name.String(), fieldName.String(), false)
					case *ast.Ident:
						t.addIdent(x, ts.Name.String(), fieldName.String(), false)
					case *ast.MapType:
						t.addMapType(x, ts.Name.String(), fieldName.String(), false)
					case *ast.ChanType:
						t.addChanType(x, ts.Name.String(), fieldName.String(), false)
					case *ast.SelectorExpr:
						t.addSelectorExpr(x, ts.Name.String(), fieldName.String(), false)
					default:
						logf("processAST: type %q, field %q, unknown %T: %+v", ts.Name, fieldName, x, x)
					}
				} else {
					switch x := se.X.(type) {
					case *ast.ArrayType:
						t.addArrayType(x, ts.Name.String(), fieldName.String(), true)
					case *ast.Ident:
						t.addIdent(x, ts.Name.String(), fieldName.String(), true)
					case *ast.MapType:
						t.addMapType(x, ts.Name.String(), fieldName.String(), true)
					case *ast.ChanType:
						t.addChanType(x, ts.Name.String(), fieldName.String(), true)
					case *ast.SelectorExpr:
						t.addSelectorExpr(x, ts.Name.String(), fieldName.String(), true)
					default:
						logf("processAST: type %q, field %q, unknown %T: %+v", ts.Name, fieldName, x, x)
					}
				}
			}
		}
	}
	return nil
}

func sourceFilter(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasSuffix(fi.Name(), fileSuffix)
}

func (t *templateData) dump() error {
	if len(t.Getters) == 0 {
		logf("No getters for %v; skipping.", t.filename)
		return nil
	}

	// Sort getters by ReceiverType.FieldName.
	sort.Sort(byName(t.Getters))

	var buf bytes.Buffer
	if err := sourceTmpl.Execute(&buf, t); err != nil {
		return err
	}
	clean, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	logf("Writing %v...", t.filename)
	return ioutil.WriteFile(t.filename, clean, 0644)
}

func newGetter(receiverType, fieldName, fieldType string, isPointer bool) *getter {
	return &getter{
		sortVal:      strings.ToLower(receiverType) + "." + strings.ToLower(fieldName),
		ReceiverVar:  strings.ToLower(receiverType[:1]),
		ReceiverType: receiverType,
		FieldName:    fieldName,
		FieldType:    fieldType,
		IsPointer:    isPointer,
	}
}

func (t *templateData) addArrayType(x *ast.ArrayType, receiverType, fieldName string, isPointer bool) {
	var eltType string
	switch elt := x.Elt.(type) {
	case *ast.Ident:
		eltType = elt.String()
	default:
		logf("addArrayType: type %q, field %q: unknown elt type: %T %+v; skipping.", receiverType, fieldName, elt, elt)
		return
	}

	t.Getters = append(t.Getters, newGetter(receiverType, fieldName, "[]"+eltType, isPointer))
}

func (t *templateData) addIdent(x *ast.Ident, receiverType, fieldName string, isPointer bool) {
	t.Getters = append(t.Getters, newGetter(receiverType, fieldName, x.String(), isPointer))
}

func (t *templateData) addMapType(x *ast.MapType, receiverType, fieldName string, isPointer bool) {
	var keyType string
	switch key := x.Key.(type) {
	case *ast.Ident:
		keyType = key.String()
	default:
		logf("addMapType: type %q, field %q: unknown key type: %T %+v; skipping.", receiverType, fieldName, key, key)
		return
	}

	var valueType string
	switch value := x.Value.(type) {
	case *ast.Ident:
		valueType = value.String()
	default:
		logf("addMapType: type %q, field %q: unknown value type: %T %+v; skipping.", receiverType, fieldName, value, value)
		return
	}

	fieldType := fmt.Sprintf("map[%v]%v", keyType, valueType)
	t.Getters = append(t.Getters, newGetter(receiverType, fieldName, fieldType, isPointer))
}

func (t *templateData) addChanType(x *ast.ChanType, receiverType, fieldName string, isPointer bool) {
	var valueType string
	switch value := x.Value.(type) {
	case *ast.Ident:
		valueType = value.String()
	default:
		logf("addChanType: type %q, field %q: unknown value type: %T %+v; skipping.", receiverType, fieldName, value, value)
		return
	}
	var fieldType string
	if x.Dir&ast.SEND == ast.SEND && x.Dir&ast.RECV == ast.RECV {
		fieldType = fmt.Sprintf("chan %v", valueType)
	} else if x.Dir&ast.SEND == 0 {
		fieldType = fmt.Sprintf("<-chan %v", valueType)
	} else {
		fieldType = fmt.Sprintf("chan<- %v", valueType)
	}
	t.Getters = append(t.Getters, newGetter(receiverType, fieldName, fieldType, isPointer))
}

func (t *templateData) addSelectorExpr(x *ast.SelectorExpr, receiverType, fieldName string, isPointer bool) {
	if strings.ToLower(fieldName[:1]) == fieldName[:1] { // Non-exported field.
		return
	}

	var xX string
	if xx, ok := x.X.(*ast.Ident); ok {
		xX = xx.String()
	}

	switch xX {
	case "time", "json":
		if xX == "json" {
			t.Imports["encoding/json"] = "encoding/json"
		} else {
			t.Imports[xX] = xX
		}
		fieldType := fmt.Sprintf("%v.%v", xX, x.Sel.Name)
		t.Getters = append(t.Getters, newGetter(receiverType, fieldName, fieldType, isPointer))
	default:
		logf("addSelectorExpr: xX %q, type %q, field %q: unknown x=%+v; skipping.", xX, receiverType, fieldName, x)
	}
}

type templateData struct {
	filename string
	Package  string
	Imports  map[string]string
	Getters  []*getter
}

type getter struct {
	sortVal      string // Lower-case version of "ReceiverType.FieldName".
	ReceiverVar  string // The one-letter variable name to match the ReceiverType.
	ReceiverType string
	FieldName    string
	FieldType    string
	IsPointer    bool
}

type byName []*getter

func (b byName) Len() int           { return len(b) }
func (b byName) Less(i, j int) bool { return b[i].sortVal < b[j].sortVal }
func (b byName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

const source = `// Code generated by genaccessor; DO NOT EDIT.

package {{.Package}}
{{with .Imports}}
import (
  {{- range . -}}
  "{{.}}"
  {{end -}}
)
{{end}}
{{range .Getters}}
{{if .IsPointer}}
// Get{{.FieldName}} returns the {{.FieldName}} field.
func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() *{{.FieldType}} {
  return {{.ReceiverVar}}.{{.FieldName}}
}

// MustSet{{.FieldName}} set the {{.FieldName}} field.
func ({{.ReceiverVar}} *{{.ReceiverType}}) MustSet{{.FieldName}}({{.FieldName}} *{{.FieldType}}) {
  {{.ReceiverVar}}.{{.FieldName}} = {{.FieldName}}
}
{{else}}
// Get{{.FieldName}} returns the {{.FieldName}} field.
func ({{.ReceiverVar}} *{{.ReceiverType}}) Get{{.FieldName}}() {{.FieldType}} {
  return {{.ReceiverVar}}.{{.FieldName}}
}

// MustSet{{.FieldName}} set the {{.FieldName}} field.
func ({{.ReceiverVar}} *{{.ReceiverType}}) MustSet{{.FieldName}}({{.FieldName}} {{.FieldType}}) {
  {{.ReceiverVar}}.{{.FieldName}} = {{.FieldName}}
}
{{end}}
{{end}}
`
