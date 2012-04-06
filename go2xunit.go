package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	startPrefix = "=== RUN "
	passPrefix = "--- PASS: "
	failPrefix = "--- FAIL: "
)

var endRegexp *regexp.Regexp = regexp.MustCompile(`([^ ]+) \((\d+\.\d+)`)

type Test struct {
	Name, Time, Message string
	Failed bool
}

func parseEnd(prefix, line string) (string, string, error) {
	matches := endRegexp.FindStringSubmatch(line[len(prefix):])

	if len(matches) == 0 {
		return "", "", fmt.Errorf("can't parse %s", line)
	}

	return matches[1], matches[2], nil
}

func parseOutput(rd io.Reader) ([]*Test, error) {
	tests := []*Test{}

	reader := bufio.NewReader(rd)
	var test *Test = nil
	for {
		/* FIXME: Handle isPrefix */
		buf, _, err := reader.ReadLine()

		switch err {
		case io.EOF:
			if test != nil {
				tests = append(tests, test)
			}
			return tests, nil
		case nil:
			;
		default:
			return nil, err
		}

		line := string(buf)
		switch {
		case strings.HasPrefix(line, startPrefix):
			if test != nil {
				tests = append(tests, test)
			}
			test = &Test{Name: line[len(startPrefix):]}
		case strings.HasPrefix(line, failPrefix):
			if test == nil {
				return nil, fmt.Errorf("fail not inside test")
			}
			test.Failed = true
			name, time, err := parseEnd(failPrefix, line)
			if err != nil {
				return nil, err
			}
			if name != test.Name {
				return nil, fmt.Errorf("wrong test end (%s!=%s)", name, test.Name)
			}
			test.Time = time
		case strings.HasPrefix(line, passPrefix):
			if test == nil {
				return nil, fmt.Errorf("pass not inside test")
			}
			test.Failed = false
			name, time, err := parseEnd(passPrefix, line)
			if err != nil {
				return nil, err
			}
			if name != test.Name {
				return nil, fmt.Errorf("wrong test end (%s!=%s)", name, test.Name)
			}
			test.Time = time
		default:
			if test != nil {
				test.Message += line + "\n"
			}
		}
	}

	return tests, nil
}

func getInput(filename string) (io.Reader, error) {
	if filename == "-" || filename == "" {
		return os.Stdin, nil
	}

	return os.Open(filename)
}

func numFailures(tests []*Test) int {
	count := 0
	for _, test := range tests {
		if test.Failed {
			count ++
		}
	}

	return count
}

func writeXML(tests []*Test, out io.Writer) {
	newline := func() { fmt.Fprintln(out) }

	fmt.Fprintf(out, `<?xml version="1.0" encoding="utf-8"?>`)
	newline()
	fmt.Fprintf(out, `<testsuite name="go2xunit" tests="%d" errors="0" failures="%d" skip="0">`,
					 len(tests), numFailures(tests))
	newline()
	for _, test := range(tests) {
		fmt.Fprintf(out, `  <testcase classname="go2xunit" name="%s" time="%s"`,
		                 test.Name, test.Time)
		if !test.Failed {
			fmt.Fprintf(out, " />\n")
			continue
		}
		fmt.Fprintln(out, ">")
		fmt.Fprintf(out,  `    <failure type="go.error" message="error">`)
		newline()
		fmt.Fprintf(out, "<![CDATA[%s]]>\n", test.Message)
		fmt.Fprintln(out, "    </failure>")
		fmt.Fprintln(out, "  </testcase>")
	}
	fmt.Fprintln(out, "</testsuite>")
}

func getOutput(filename string) (io.Writer, error) {
	if filename == "-" || filename == "" {
		return os.Stdout, nil
	}

	return os.Create(filename)
}

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

func fatal(formt string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, formt, args...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func main() {
	inputFile := flag.String("input", "", "input file (default to stdin)")
	outputFile := flag.String("output", "", "output file (default to stdout)")
	//run := flag.Bool("run", false, "run go test yourself")
	flag.Parse()

	if flag.NArg() > 0 {
		fatal("error: %s does not take parameters (did you mean -input?)", os.Args[0])
	}

	/*
	if len(*inputFile) > 0 && *run {
		log.Fatalf("error: can't specify -run and -input at the same time\n")
	}
	*/

	input, output, err := getIO(*inputFile, *outputFile)
	if err != nil {
		fatal("error: %s", err)
	}


	tests, err := parseOutput(input)
	if err != nil {
		fatal("error: %s", err)
	}

	writeXML(tests, output)
}
