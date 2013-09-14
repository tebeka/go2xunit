package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"bitbucket.org/tebeka/go2xunit/types"
	"bitbucket.org/tebeka/go2xunit/gocheck"
	"bitbucket.org/tebeka/go2xunit/gotest"
)

const (
	version = "0.2.1"
)


type TestResults struct {
	Suites []*types.Suite
	Bamboo bool
}


func hasFailures(suites []*types.Suite) bool {
	for _, suite := range suites {
		if suite.NumFailed() > 0 {
			return true
		}
	}
	return false
}

var xmlTemplate string = `<?xml version="1.0" encoding="utf-8"?>
{{if .Bamboo}}<testsuites>{{end}}
{{range $suite := .Suites}}  <testsuite name="{{.Name}}" tests="{{.Count}}" errors="0" failures="{{.NumFailed}}" skip="0">
{{range  $test := $suite.Tests}}    <testcase classname="{{$suite.Name}}" name="{{$test.Name}}" time="{{$test.Time}}">
{{if $test.Failed }}      <failure type="go.error" message="error">
        <![CDATA[{{$test.Message}}]]>
      </failure>{{end}}    </testcase>
{{end}}  </testsuite>
{{end}}{{if .Bamboo}}</testsuites>{{end}}
`

// writeXML exits xunit XML of tests to out
func writeXML(suites []*types.Suite, out io.Writer, bamboo bool) {
	testsResult := TestResults{Suites: suites, Bamboo: bamboo}
	t := template.New("test template")
	t, err := t.Parse(xmlTemplate)
	if err != nil {
		fmt.Println("Error en parse %v", err)
		return
	}
	err = t.Execute(out, testsResult)
	if err != nil {
		fmt.Println("Error en execute %v", err)
		return
	}
}

// getInput return input io.Reader from file name, if file name is - it will
// return os.Stdin
func getInput(filename string) (io.Reader, error) {
	if filename == "-" || filename == "" {
		return os.Stdin, nil
	}

	return os.Open(filename)
}

// getInput return output io.Writer from file name, if file name is - it will
// return os.Stdout
func getOutput(filename string) (io.Writer, error) {
	if filename == "-" || filename == "" {
		return os.Stdout, nil
	}

	return os.Create(filename)
}

// getIO returns input and output streams from file names
func getIO(inputFile, outputFile string) (io.Reader, io.Writer, error) {
	input, err := getInput(inputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %s for reading: %s", inputFile, err)
	}

	output, err := getOutput(outputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %s for writing: %s", outputFile, err)
	}

	return input, output, nil
}

func main() {
	inputFile := flag.String("input", "", "input file (default to stdin)")
	outputFile := flag.String("output", "", "output file (default to stdout)")
	fail := flag.Bool("fail", false, "fail (non zero exit) if any test failed")
	showVersion := flag.Bool("version", false, "print version and exit")
	bamboo := flag.Bool("bamboo", false, "xml compatible with Atlassian's Bamboo")
	is_gocheck := flag.Bool("gocheck", false, "parse gocheck output")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// No time ... prefix for error messages
	log.SetFlags(0)

	if flag.NArg() > 0 {
		log.Fatalf("error: %s does not take parameters (did you mean -input?)", os.Args[0])
	}

	input, output, err := getIO(*inputFile, *outputFile)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	var parse func(rd io.Reader) ([]*types.Suite, error)

	if *is_gocheck {
		parse = gocheck.Parse
	} else {
		parse = gotest.Parse
	}

	suites, err := parse(input)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	if len(suites) == 0 {
		log.Fatalf("error: no tests found")
		os.Exit(1)
	}

	writeXML(suites, output, *bamboo)
	if *fail && hasFailures(suites) {
		os.Exit(1)
	}
}
