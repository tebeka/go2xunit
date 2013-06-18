package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	startPrefix = "=== RUN "
	passPrefix  = "--- PASS: "
	failPrefix  = "--- FAIL: "

	version = "0.1.2"
)

// "end of test" regexp for name and time, examples:
// --- PASS: TestSub (0.00 seconds)
// --- FAIL: TestSubFail (0.00 seconds)
var endRegexp *regexp.Regexp = regexp.MustCompile(`([^ ]+) \((\d+\.\d+)`)

type Test struct {
	Name, Time, Message string
	Failed              bool
}

type TestResults struct {
	Tests  []*Test
	Count  int
	Failed int
	Bamboo bool
}

// parseEnd parses "end of test" line and returns (name, time, error)
func parseEnd(prefix, line string) (string, string, error) {
	matches := endRegexp.FindStringSubmatch(line[len(prefix):])

	if len(matches) == 0 {
		return "", "", fmt.Errorf("can't parse %s", line)
	}

	return matches[1], matches[2], nil
}

// parseOutput parses output of "go test -v", returns a list of tests
func parseOutput(rd io.Reader) ([]*Test, error) {
	tests := []*Test{}
	var test *Test = nil

	nextTest := func() {
		// We are switching to the next test, store the current one.
		if test == nil {
			return
		}

		tests = append(tests, test)
		test = nil
	}

	reader := bufio.NewReader(rd)
	for {
		buf, _, err := reader.ReadLine()

		switch err {
		case io.EOF:
			nextTest()
			return tests, nil
		case nil:
			// nil is OK

		default: // Error other than io.EOF
			return nil, err
		}

		line := string(buf)

		switch {
		case strings.HasPrefix(line, startPrefix):
		case strings.HasPrefix(line, failPrefix):
			nextTest()

			// Extract the test name and the duration:
			name, time, err := parseEnd(passPrefix, line)
			if err != nil {
				return nil, err
			}

			test = &Test{
				Name:   name,
				Time:   time,
				Failed: true,
			}

		case strings.HasPrefix(line, passPrefix):
			nextTest()
			// Extract the test name and the duration:
			name, time, err := parseEnd(passPrefix, line)
			if err != nil {
				return nil, err
			}

			// Create the test structure and store it.
			tests = append(tests, &Test{
				Name:   name,
				Time:   time,
				Failed: false,
			})
			test = nil
		case line == "FAIL":
			nextTest()
		default:
			if test != nil { // test != nil marks we're in the middle of a test
				test.Message += line + "\n"
			}
		}
	}

	// If we're here, it's an error
	return nil, fmt.Errorf("Error parsing")
}

// numFailures count how man tests failed
func numFailures(tests []*Test) int {
	count := 0
	for _, test := range tests {
		if test.Failed {
			count++
		}
	}

	return count
}

var xmlTemplate string = `<?xml version="1.0" encoding="utf-8"?>
{{if .Bamboo}}<testsuites>{{end}}
  <testsuite name="go2xunit" tests="{{.Count}}" errors="0" failures="{{.Failed}}" skip="0">
{{range $test := .Tests}}    <testcase classname="go2xunit" name="{{$test.Name}}" time="{{$test.Time}}">
{{if $test.Failed }}      <failure type="go.error" message="error">
        <![CDATA[{{$test.Message}}]]></failure>
{{end}}    </testcase>
{{end}}  </testsuite>
{{if .Bamboo}}</testsuites>{{end}}
	`

// writeXML exits xunit XML of tests to out
func writeXML(tests []*Test, out io.Writer, bamboo bool) {
	count := len(tests)
	failed := numFailures(tests)

	testsResult := TestResults{Tests: tests, Count: count, Failed: failed, Bamboo: bamboo}
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

	tests, err := parseOutput(input)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	if len(tests) == 0 {
		log.Fatalf("error: no tests found")
		os.Exit(1)
	}

	writeXML(tests, output, *bamboo)
	if *fail && numFailures(tests) > 0 {
		os.Exit(1)
	}
}
