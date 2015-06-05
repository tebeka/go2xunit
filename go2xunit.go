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
	version = "1.1.1"

	// gotest regular expressions

	// === RUN TestAdd
	gt_startRE = "^=== RUN:? ([a-zA-Z_][^[:space:]]*)"

	// --- PASS: TestSub (0.00 seconds)
	// --- FAIL: TestSubFail (0.00 seconds)
	// --- SKIP: TestSubSkip (0.00 seconds)
	gt_endRE = "^--- (PASS|FAIL|SKIP): ([a-zA-Z_][^[:space:]]*) \\((\\d+(.\\d+)?)"

	// FAIL	_/home/miki/Projects/goroot/src/xunit	0.004s
	// ok  	_/home/miki/Projects/goroot/src/anotherTest	0.000s
	gt_suiteRE = "^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(\\d+.\\d+)"

	// ?       alipay  [no test files]
	gt_noFiles = "^\\?.*\\[no test files\\]$"
	// FAIL    node/config [build failed]
	gt_buildFailed = `^FAIL.*\[(build|setup) failed\]$`

	// gocheck regular expressions

	// START: mmath_test.go:16: MySuite.TestAdd
	gc_startRE = "START: [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)"
	// PASS: mmath_test.go:16: MySuite.TestAdd	0.000s
	// FAIL: mmath_test.go:35: MySuite.TestDiv
	gc_endRE = "(PASS|FAIL): [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)([[:space:]]+([0-9]+.[0-9]+))?"
)

var (
	failOnRace = false
)

type Test struct {
	Name, Time, Message string
	Failed              bool
	Skipped             bool
	Passed              bool
}

type Suite struct {
	Name   string
	Time   string
	Status string
	Tests  []*Test
}

type SuiteStack struct {
	nodes []*Suite
	count int
}

// Push adds a node to the stack.
func (s *SuiteStack) Push(n *Suite) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *SuiteStack) Pop() *Suite {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.nodes[s.count]
}

type TestResults struct {
	Suites []*Suite
	Multi  bool
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

func (suite *Suite) NumSkipped() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Skipped {
			count++
		}
	}
	return count
}

func (suite *Suite) NumPassed() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Passed {
			count++
		}
	}
	return count
}

func (suite *Suite) Count() int {
	return len(suite.Tests)
}

func hasDatarace(lines []string) bool {
	has_datarace := regexp.MustCompile("^WARNING: DATA RACE$").MatchString
	for _, line := range lines {
		if has_datarace(line) {
			return true
		}
	}
	return false
}

func gt_Parse(rd io.Reader) ([]*Suite, error) {
	find_start := regexp.MustCompile(gt_startRE).FindStringSubmatch
	find_end := regexp.MustCompile(gt_endRE).FindStringSubmatch
	find_suite := regexp.MustCompile(gt_suiteRE).FindStringSubmatch
	is_nofiles := regexp.MustCompile(gt_noFiles).MatchString
	is_buildFailed := regexp.MustCompile(gt_buildFailed).MatchString
	is_exit := regexp.MustCompile("^exit status -?\\d+").MatchString

	suites := []*Suite{}
	var curTest *Test
	var curSuite *Suite
	var out []string
	suiteStack := SuiteStack{}
	// Handles a test that ended with a panic.
	handlePanic := func() {
		curTest.Failed = true
		curTest.Skipped = false
		curTest.Passed = false
		curTest.Time = "N/A"
		curSuite.Tests = append(curSuite.Tests, curTest)
		curTest = nil
	}

	// Appends output to the last test.
	appendError := func() error {
		if len(out) > 0 && curSuite != nil && len(curSuite.Tests) > 0 {
			message := strings.Join(out, "\n")
			if curSuite.Tests[len(curSuite.Tests)-1].Message == "" {
				curSuite.Tests[len(curSuite.Tests)-1].Message = message
			} else {
				curSuite.Tests[len(curSuite.Tests)-1].Message += "\n" + message
			}
		}
		out = []string{}
		return nil
	}

	scanner := bufio.NewScanner(rd)
	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()

		// TODO: Only outside a suite/test, report as empty suite?
		if is_nofiles(line) {
			continue
		}

		if is_buildFailed(line) {
			return nil, fmt.Errorf("%d: package build failed: %s", lnum, line)
		}

		if curSuite == nil {
			curSuite = &Suite{}
		}

		tokens := find_start(line)
		if tokens != nil {
			if curTest != nil {
				// This occurs when the last test ended with a panic.
				if suiteStack.count == 0 {
					suiteStack.Push(curSuite)
					curSuite = &Suite{Name: curTest.Name}
				} else {
					handlePanic()
				}
			}
			if e := appendError(); e != nil {
				return nil, e
			}
			curTest = &Test{
				Name: tokens[1],
			}
			continue
		}

		tokens = find_end(line)
		if tokens != nil {
			if curTest == nil {
				if suiteStack.count > 0 {
					prevSuite := suiteStack.Pop()
					suites = append(suites, curSuite)
					curSuite = prevSuite
					continue
				} else {
					return nil, fmt.Errorf("%d: orphan end test", lnum)
				}
			}
			if tokens[2] != curTest.Name {
				err := fmt.Errorf("%d: name mismatch (try disabling parallel mode)", lnum)
				return nil, err
			}
			curTest.Failed = (tokens[1] == "FAIL") || (failOnRace && hasDatarace(out))
			curTest.Skipped = (tokens[1] == "SKIP")
			curTest.Passed = (tokens[1] == "PASS")
			curTest.Time = tokens[3]
			curTest.Message = strings.Join(out, "\n")
			curSuite.Tests = append(curSuite.Tests, curTest)
			curTest = nil
			out = []string{}
			continue
		}

		tokens = find_suite(line)
		if tokens != nil {
			if curTest != nil {
				// This occurs when the last test ended with a panic.
				handlePanic()
			}
			if e := appendError(); e != nil {
				return nil, e
			}
			curSuite.Name = tokens[2]
			curSuite.Time = tokens[3]
			suites = append(suites, curSuite)
			curSuite = nil
			continue
		}

		if is_exit(line) || (line == "FAIL") || (line == "PASS") {
			continue
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
			test.Passed = (tokens[1] == "PASS")

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

const xmlDeclaration = `<?xml version="1.0" encoding="utf-8"?>`

const xunitTemplate string = `
{{range $suite := .Suites}}  <testsuite name="{{.Name}}" tests="{{.Count}}" errors="0" failures="{{.NumFailed}}" skip="{{.NumSkipped}}">
{{range  $test := $suite.Tests}}    <testcase classname="{{$suite.Name}}" name="{{$test.Name}}" time="{{$test.Time}}">
{{if $test.Skipped }}      <skipped/> {{end}}
{{if $test.Failed }}      <failure type="go.error" message="error">
        <![CDATA[{{$test.Message}}]]>
      </failure>{{end}}    </testcase>
{{end}}  </testsuite>
{{end}}`

const bambooTemplate string = `
<testsuites>` + xunitTemplate + `</testsuites>
`

// https://xunit.codeplex.com/wikipage?title=XmlFormat
const xunitNetTemplate string = `
<assembly name="go test"
          run-date="2014-09-22" run-time="0"
          configFile=""
          time="0"
          total="0"
          passed="0"
          failed="0"
          skipped="0"
          environment="n/a"
          test-framework="golang">
{{range $suite := .Suites}}
    <class time="{{.Time}}" name="{{.Name}}"
  	     total="{{.Count}}"
  	     passed="{{.NumPassed}}"
  	     failed="{{.NumFailed}}"
  	     skipped="{{.NumSkipped}}">
{{range  $test := $suite.Tests}}
        <test name="{{$test.Name}}"
          type="test"
          method="{{$test.Name}}"
          result={{if $test.Skipped }}"Skip"{{else if $test.Failed }}"Fail"{{else if $test.Passed }}"Pass"{{end}}
          time="{{$test.Time}}">
        {{if $test.Failed }}  <failure exception-type="go.error">
             <message><![CDATA[{{$test.Message}}]]></message>
      	  </failure>
      	{{end}}</test>
{{end}}
    </class>
{{end}}
</assembly>
`

// writeXML exits xunit XML of tests to out
func writeXML(suites []*Suite, out io.Writer, xmlTemplate string) {
	testsResult := TestResults{
		Suites: suites,
	}
	t := template.New("test template")
	t, err := t.Parse(xmlDeclaration + xmlTemplate)
	if err != nil {
		fmt.Printf("Error in parse %v\n", err)
		return
	}
	err = t.Execute(out, testsResult)
	if err != nil {
		fmt.Printf("Error in execute %v\n", err)
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
	xunitnet := flag.Bool("xunitnet", false, "xml compatible with xunit.net")
	is_gocheck := flag.Bool("gocheck", false, "parse gocheck output")
	flag.BoolVar(&failOnRace, "fail-on-race", false, "mark test as failing if it exposes a data race")
	flag.Parse()

	if *showVersion {
		fmt.Printf("go2xunit %s\n", version)
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
	}

	if *bamboo && *xunitnet {
		log.Fatalf("error: -bamboo and -xunitnet are mutually exclusive")
	}

	var xmlTemplate string
	if *xunitnet {
		xmlTemplate = xunitNetTemplate
	} else if *bamboo || (len(suites) > 1) {
		xmlTemplate = bambooTemplate
	} else {
		xmlTemplate = xunitTemplate
	}

	writeXML(suites, output, xmlTemplate)
	if *fail && hasFailures(suites) {
		os.Exit(1)
	}
}
