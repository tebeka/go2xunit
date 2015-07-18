package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Time when the test was run
var runTime time.Time

const (
	version = "1.2.2"

	// gotest regular expressions

	// === RUN TestAdd
	gt_startRE = "^=== RUN:?[[:space:]]+([a-zA-Z_][^[:space:]]*)"

	// --- PASS: TestSub (0.00 seconds)
	// --- FAIL: TestSubFail (0.00 seconds)
	// --- SKIP: TestSubSkip (0.00 seconds)
	gt_endRE = "^--- (PASS|FAIL|SKIP):[[:space:]]+([a-zA-Z_][^[:space:]]*) \\((\\d+(.\\d+)?)"

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
	gc_endRE = "(PASS|FAIL|SKIP|PANIC|MISS): [^:]+:[^:]+: ([A-Za-z_][[:word:]]*).([A-Za-z_][[:word:]]*)[[:space:]]?([0-9]+.[0-9]+)?"

	// FAIL	go2xunit/demo-gocheck	0.008s
	// ok  	go2xunit/demo-gocheck	0.008s
	gc_suiteRE = "^(ok|FAIL)[ \t]+([^ \t]+)[ \t]+(\\d+.\\d+)"
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
	Suites   []*Suite
	Assembly string
	RunDate  string
	RunTime  string
	Time     string
	Total    int
	Passed   int
	Failed   int
	Skipped  int
}

// calculate grand total for all suites
func (r *TestResults) calcTotals() {
	totalTime, _ := strconv.ParseFloat(r.Time, 64)
	for _, suite := range r.Suites {
		r.Passed += suite.NumPassed()
		r.Failed += suite.NumFailed()
		r.Skipped += suite.NumSkipped()

		suiteTime, _ := strconv.ParseFloat(suite.Time, 64)
		totalTime += suiteTime
		r.Time = fmt.Sprintf("%.3f", totalTime)
	}
	r.Total = r.Passed + r.Skipped + r.Failed
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

// Number of passed tests in the suite
func (suite *Suite) NumPassed() int {
	count := 0
	for _, test := range suite.Tests {
		if test.Passed {
			count++
		}
	}

	return count
}

// Number of tests in the suite
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

// gc_Parse parses output of "go test -gocheck.vv", returns a list of tests
// See data/gocheck.out for an example
func gc_Parse(rd io.Reader) ([]*Suite, error) {
	find_start := regexp.MustCompile(gc_startRE).FindStringSubmatch
	find_end := regexp.MustCompile(gc_endRE).FindStringSubmatch
	find_suite := regexp.MustCompile(gc_suiteRE).FindStringSubmatch

	scanner := bufio.NewScanner(rd)
	var suites = make([]*Suite, 0)
	var suiteName string
	var suite *Suite

	var testName string
	var out []string

	for lnum := 1; scanner.Scan(); lnum++ {
		line := scanner.Text()

		tokens := find_start(line)
		if len(tokens) > 0 {
			if tokens[2] == "SetUpTest" || tokens[2] == "TearDownTest" {
				continue
			}
			if testName != "" {
				return nil, fmt.Errorf("%d: start in middle\n", lnum)
			}
			suiteName = tokens[1]
			testName = tokens[2]
			out = []string{}
			continue
		}

		tokens = find_end(line)
		if len(tokens) > 0 {
			if tokens[3] == "SetUpTest" || tokens[3] == "TearDownTest" {
				continue
			}
			if testName == "" {
				return nil, fmt.Errorf("%d: orphan end", lnum)
			}
			if (tokens[2] != suiteName) || (tokens[3] != testName) {
				return nil, fmt.Errorf("%d: suite/name mismatch", lnum)
			}
			test := &Test{Name: testName}
			test.Message = strings.Join(out, "\n")
			test.Time = tokens[4]
			test.Failed = (tokens[1] == "FAIL") || (tokens[1] == "PANIC")
			test.Passed = (tokens[1] == "PASS")
			test.Skipped = (tokens[1] == "SKIP" || tokens[1] == "MISS")

			if suite == nil || suite.Name != suiteName {
				suite = &Suite{Name: suiteName}
				suites = append(suites, suite)
			}
			suite.Tests = append(suite.Tests, test)

			testName = ""
			suiteName = ""
			out = []string{}

			continue
		}

		// last "suite" is test summary
		tokens = find_suite(line)
		if tokens != nil {
			if suite == nil {
				suite = &Suite{Name: tokens[2], Status: tokens[1], Time: tokens[3]}
				suites = append(suites, suite)
			} else {
				suite.Status = tokens[1]
				suite.Time = tokens[3]
			}

			testName = ""
			suiteName = ""
			out = []string{}

			continue
		}

		if testName != "" {
			out = append(out, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}

func hasFailures(suites []*Suite) bool {
	for _, suite := range suites {
		if suite.NumFailed() > 0 {
			return true
		}
	}
	return false
}

const (
	xmlDeclaration = `<?xml version="1.0" encoding="utf-8"?>`

	xunitTemplate string = `
{{range $suite := .Suites}}  <testsuite name="{{.Name}}" tests="{{.Count}}" errors="0" failures="{{.NumFailed}}" skip="{{.NumSkipped}}">
{{range  $test := $suite.Tests}}    <testcase classname="{{$suite.Name}}" name="{{$test.Name}}" time="{{$test.Time}}">
{{if $test.Skipped }}      <skipped/> {{end}}
{{if $test.Failed }}      <failure type="go.error" message="error">
        <![CDATA[{{$test.Message}}]]>
      </failure>{{end}}    </testcase>
{{end}}  </testsuite>
{{end}}`

	multiTemplate string = `
<testsuites>` + xunitTemplate + `</testsuites>
`

	// https://xunit.codeplex.com/wikipage?title=XmlFormat
	xunitNetTemplate string = `
<assembly name="{{.Assembly}}"
          run-date="{{.RunDate}}" run-time="{{.RunTime}}"
          configFile="none"
          time="{{.Time}}"
          total="{{.Total}}"
          passed="{{.Passed}}"
          failed="{{.Failed}}"
          skipped="{{.Skipped}}"
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
)

// writeXML exits xunit XML of tests to out
func writeXML(suites []*Suite, out io.Writer, xmlTemplate string) {
	testsResult := TestResults{
		Suites:   suites,
		Assembly: suites[len(suites)-1].Name,
		RunDate:  runTime.Format("2006-01-02"),
		RunTime:  fmt.Sprintf("%02d:%02d:%02d", runTime.Hour(), runTime.Minute(), runTime.Second()),
	}
	testsResult.calcTotals()
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

// getInput return input io.File from file name, if file name is - it will
// return os.Stdin
func getInput(filename string) (*os.File, error) {
	if filename == "-" || filename == "" {
		return os.Stdin, nil
	}

	return os.Open(filename)
}

// getInput return output io.File from file name, if file name is - it will
// return os.Stdout
func getOutput(filename string) (*os.File, error) {
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

	setRunTimeFrom(input)

	return input, output, nil
}

// set test execution time from file date
func setRunTimeFrom(file *os.File) {
	statinfo, err := file.Stat()
	checkFatal(err)
	runTime = statinfo.ModTime()
}

func checkFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
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

	if *bamboo && *xunitnet {
		log.Fatalf("error: -bamboo and -xunitnet are mutually exclusive")
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

	xmlTemplate := xunitTemplate
	if *xunitnet {
		xmlTemplate = xunitNetTemplate
	} else if *bamboo || (len(suites) > 1) {
		xmlTemplate = multiTemplate
	}

	writeXML(suites, output, xmlTemplate)
	if *fail && hasFailures(suites) {
		os.Exit(1)
	}
}
