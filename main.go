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
	"sort"
	"strings"
	"text/template"
)

const (
	// Version is the current version
	Version = "2.0.0"
)

var args struct {
	failRace    bool
	format      string
	input       string
	noFail      bool
	output      string
	suitePrefix string
	version     bool
}

func templateNames() string {
	names := make([]string, 0, len(internalTemplates))
	for name := range internalTemplates {
		names = append(names, name)
	}

	sort.Strings(names)
	return strings.Join(names, ", ")
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]\nOptions:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	formatHelp := fmt.Sprintf("output format: %s", templateNames())

	flag.BoolVar(&args.failRace, "fail-race", false, "mark test as failing if it exposes a data race")
	flag.BoolVar(&args.noFail, "no-fail", false, "don't fail if tests failed")
	flag.BoolVar(&args.version, "version", false, "print version and exit")
	flag.StringVar(&args.format, "format", "junit", formatHelp)
	flag.StringVar(&args.input, "input", "", "input file")
	flag.StringVar(&args.output, "output", "", "output file")
	flag.StringVar(&args.suitePrefix, "suite-prefix", "", "prefix to include before all suite names")
}

func inFile(name string) (*os.File, error) {
	if name == "" || name == "-" {
		return os.Stdin, nil
	}

	return os.Open(name)
}

func outFile(name string) (*os.File, error) {
	if name == "" || name == "-" {
		return os.Stdout, nil
	}

	return os.Create(name)
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

	flag.Parse()
	if flag.NArg() > 0 {
		log.Fatalf("error: %s takes no arguments", os.Args[0])
	}

	input, err := inFile(args.input)
	if err != nil {
		log.Fatalf("error: input: %s", err)
	}

	out, err := outFile(args.output)
	if err != nil {
		log.Fatalf("error: output: %s", err)
	}

	tmplData := internalTemplates[args.format]
	if tmplData == "" {
		log.Fatalf("error: can't find %q template", args.format)
	}

	root, err := Parse(input)
	if err != nil {
		log.Fatalf("error: can't parse - %s", err)
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
