package main

//go:generate go run ./scripts/gentmpl.go
//go:generate go fmt templates.go

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"text/template"
)

const (
	// Version is the current version
	Version = "2.0.0"
)

func inFile() (*os.File, error) {
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		return os.Stdin, nil
	}

	return os.Open(flag.Arg(0))
}

func outFile(path string) (*os.File, error) {
	if path == "" || path == "-" {
		return os.Stdout, nil
	}

	return os.Create(flag.Arg(1))
}

func xmlEscape(in string) (string, error) {
	w := &bytes.Buffer{}
	if err := xml.EscapeText(w, []byte(in)); err != nil {
		return "", fmt.Errorf("error escaping text: %s", err)
	}
	return w.String(), nil
}

// getInput return input io.File from file name, if file name is - it will
// return os.Stdin
func main() {
	if args.version {
		fmt.Printf("%s %s\n", Version, path.Base(os.Args[0]))
		os.Exit(0)
	}

	// No time ... prefix for error messages
	log.SetFlags(0)

	if err := parseArgs(); err != nil {
		log.Fatalf("error: %s", err)
	}

	input, err := inFile()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	out, err := outFile(args.output)
	if err != nil {
		log.Fatalf("error:  %s", err)
	}

	root, err := Parse(input)
	if err != nil {
		log.Fatalf("error: can't parse - %s", err)
	}

	tmplData := Templates[args.format]
	if tmplData == "" {
		log.Fatalf("error: can't find tempalte for %q", args.format)
	}

	funcs := template.FuncMap{
		"escape": xmlEscape,
	}
	tmpl, err := template.New(args.format).Funcs(funcs).Parse(tmplData)
	if err != nil {
		log.Fatalf("error: can't compile template %s - %s", args.format, err)
	}

	if err = tmpl.Execute(out, root); err != nil {
		log.Fatalf("error: can't execute template - %s", err)
	}
}
