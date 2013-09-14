package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"bufio"
	"strings"
	"text/template"
)

const (
	version = "0.2.1"

	// gotest regular expressions

	// === RUN TestAdd
	gt_startRE = "^=== RUN ([a-zA-Z_][[:word:]]*)"

	// --- PASS: TestSub (0.00 seconds)
	// --- FAIL: TestSubFail (0.00 seconds)
	gt_endRE = "^--- (PASS|FAIL): ([a-zA-Z_][[:word:]]*) \\((\\d+.\\d+)"

	// FAIL	_/home/miki/Projects/goroot/src/xunit	0.004s
	// ok  	_/home/miki/Projects/goroot/src/anotherTest	0.000s
	gt_suiteRE = "^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(\\d+.\\d+)"

	// gocheck regular expressions

	// START: mmath_test.go:16: MySuite.TestAdd
	gc_startRE = "START: [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)"
	// PASS: mmath_test.go:16: MySuite.TestAdd	0.000s
	// FAIL: mmath_test.go:35: MySuite.TestDiv
	gc_endRE = "(PASS|FAIL): [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)([[:space:]]+([0-9]+.[0-9]+))?"
)

type Test struct {
	Name, Time, Message string
	Failed              bool
}

type Suite struct {
	Name   string
	Time   string
	Status string
	Tests  []*Test
}

type TestResults struct {
	Suites []*Suite
	Bamboo bool
}

func (suite *Suite) NumFailed() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Failed {
			count++
		}
	}

	return count
}

func (suite *Suite) Count() int {
	return len(suite.Tests)
}

func gt_Parse(rd io.Reader) ([]*Suite, error) {
	find_start := regexp.MustCompile(gt_startRE).FindStringSubmatch
	find_end := regexp.MustCompile(gt_endRE).FindStringSubmatch
	find_suite := regexp.MustCompile(gt_suiteRE).FindStringSubmatch
	is_exit := regexp.MustCompile("^exit status -?\\d+").MatchString

	suites := []*Suite{}
	var curTest *Test
	var curSuite *Suite
	var out []string

	scanner := bufio.NewScanner(rd)
	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()

		tokens := find_start(line)
		if tokens != nil {
			if curTest != nil {
				return nil, fmt.Errorf("%d: test in middle of other", lnum)
			}
			curTest = &Test{
				Name: tokens[1],
			}
			if len(out) > 0 {
				message := strings.Join(out, "\n")
				if curSuite == nil {
					return nil, fmt.Errorf("orphan output: %s", message)
				}
				curSuite.Tests[len(curSuite.Tests)-1].Message = message
			}
			out = []string{}
			continue
		}

		tokens = find_end(line)
		if tokens != nil {
			if curTest == nil {
				return nil, fmt.Errorf("%d: orphan end test", lnum)
			}
			if tokens[2] != curTest.Name {
				return nil, fmt.Errorf("%d: name mismatch", lnum)
			}

			curTest.Failed = (tokens[1] == "FAIL")
			curTest.Time = tokens[3]
			curTest.Message = strings.Join(out, "\n")
			if curSuite == nil {
				curSuite = &Suite{}
			}
			curSuite.Tests = append(curSuite.Tests, curTest)
			curTest = nil
			continue
		}

		tokens = find_suite(line)
		if tokens != nil {
			if curSuite == nil {
				return nil, fmt.Errorf("%d: orphan end suite", lnum)
			}
			curSuite.Name = tokens[2]
			curSuite.Time = tokens[3]
			suites = append(suites, curSuite)
			curSuite = nil

			continue
		}

		if is_exit(line) || (line == "FAIL") {
			continue
		}

		if curSuite == nil {
			return nil, fmt.Errorf("%d: orphan line", lnum)
		}

		out = append(out, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}

func map2arr(m map[string]*Suite) []*Suite {
	arr := make([]*Suite, 0, len(m))
	for _, suite := range m {
		/* FIXME:
		suite.Status =
		suite.Time =
		*/
		arr = append(arr, suite)
	}

	return arr
}

// gc_Parse parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
func gc_Parse(rd io.Reader) ([]*Suite, error) {
	find_start := regexp.MustCompile(gc_startRE).FindStringSubmatch
	find_end := regexp.MustCompile(gc_endRE).FindStringSubmatch

	scanner := bufio.NewScanner(rd)
	var test *Test
	var suites = make(map[string]*Suite)
	var suiteName string
	var out []string

	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()
		tokens := find_start(line)
		if len(tokens) > 0 {
			if test != nil {
				return nil, fmt.Errorf("%d: start in middle\n", lnum)
			}
			suiteName = tokens[1]
			test = &Test{Name: tokens[2]}
			out = []string{}
			continue
		}

		tokens = find_end(line)
		if len(tokens) > 0 {
			if test == nil {
				return nil, fmt.Errorf("%d: orphan end", lnum)
			}
			if (tokens[2] != suiteName) || (tokens[3] != test.Name) {
				return nil, fmt.Errorf("%d: suite/name mismatch", lnum)
			}
			test.Message = strings.Join(out, "\n")
			test.Time = tokens[4]
			test.Failed = (tokens[1] == "FAIL")

			suite, ok := suites[suiteName]
			if !ok {
				suite = &Suite{Name: suiteName}
			}
			suite.Tests = append(suite.Tests, test)
			suites[suiteName] = suite

			test = nil
			suiteName = ""
			out = []string{}

			continue
		}

		if test != nil {
			out = append(out, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return map2arr(suites), nil
}

func hasFailures(suites []*Suite) bool {
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
func writeXML(suites []*Suite, out io.Writer, bamboo bool) {
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

	var parse func(rd io.Reader) ([]*Suite, error)

	if *is_gocheck {
		parse = gc_Parse
	} else {
		parse = gt_Parse
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
