package main

//go:generate go run ./scripts/gentmpl.go
//go:generate go fmt templates.go

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	// Version is the current version
	Version = "1.4.8"
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

	output, err := outFile(args.output)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	root, err := Parse(input)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(output, "\n\n%+v\n", root)
	for _, t := range root.Children {
		fmt.Fprintf(output, "\t%s [%s]\n", t.Name, t.Status)
	}
}
