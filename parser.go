package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"time"
)

var (
	zeroTime time.Time
	// FAIL	github.com/nuclio/nuclio/pkg/dashboard/test [build failed]
	buildFailRe = regexp.MustCompile("FAIL[\\t ]+([^ ]+)[\\t ]+\\[build failed\\]")
)

// Test data structure
type Test struct {
	Name     string
	Package  string
	Time     time.Time
	Status   string
	Children []*Test
	Message  string
	Elapsed  time.Duration
	Stats    map[string]int

	records []*Record
}

// Record is test in JSON output
type Record struct {
	Test    string
	Time    time.Time
	Action  string
	Package string
	Output  string
	Elapsed float64
}

// Sort tests by time
type byTime []*Record

func (b byTime) Len() int           { return len(b) }
func (b byTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byTime) Less(i, j int) bool { return b[i].Time.Before(b[j].Time) }

type key struct {
	pkg  string
	test string
}

// Parse parsers test output in JSON format
func Parse(input io.Reader) (*Test, error) {
	tests, err := firstScan(input)
	if err != nil {
		return nil, err
	}

	return assembleTests(tests)
}

func firstScan(input io.Reader) (map[key]*Test, error) {
	tests := make(map[key]*Test)
	scan := NewScanner(input)
	//	tests := make(map[string]*Test)
	for scan.Scan() {
		r := &Record{}
		if err := json.Unmarshal(scan.Bytes(), r); err != nil {
			fields := buildFailRe.FindSubmatch(scan.Bytes())
			if len(fields) == 2 {
				r.Action = "fail"
				r.Package = string(fields[1])
				r.Output = "build failed"
				// TODO: r.Time
			} else {
				return nil, fmt.Errorf("%d: error: %s", scan.LineNum(), err)
			}
		}

		k := key{r.Package, r.Test}
		t, ok := tests[k]
		if !ok {
			t = &Test{}
			tests[k] = t
		}
		t.records = append(t.records, r)
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return tests, nil
}

func assembleTests(tests map[key]*Test) (*Test, error) {
	var root *Test
	for _, t := range tests {
		if err := t.assemble(); err != nil {
			return nil, err
		}
	}

	root = &Test{}
	for _, t := range tests {
		root.Children = append(root.Children, t)
	}

	for _, t := range root.Children {
		if isEmptyTime(root.Time) || t.Time.Before(root.Time) {
			root.Time = t.Time
		}
	}

	root.calcStats()
	return root, nil
}

// Count return number of sub tests (including this test)
func (t *Test) Count() int {
	n := 1
	for _, c := range t.Children {
		n += c.Count()
	}

	return n
}

// Stats returns the number of tests and subtests that have status
func (t *Test) calcStats() {
	if t.Stats != nil {
		return
	}

	stats := map[string]int{
		"fail": 0,
		"pass": 0,
		"skip": 0,
	}

	stats[t.Status]++
	for _, c := range t.Children {
		c.calcStats()
		for key := range stats {
			stats[key] += c.Stats[key]
		}
	}

	t.Stats = stats
}

func (t *Test) assemble() error {
	sort.Sort(byTime(t.records))
	var buf bytes.Buffer
	for _, r := range t.records {
		if t.Name == "" {
			t.Name = r.Test
		}
		if t.Package == "" {
			t.Package = r.Package
		}
		if r.Output != "" {
			buf.WriteString(r.Output)
		}

		switch r.Action {
		case "run":
			t.Time = r.Time
		case "output":
		case "pass", "fail", "skip":
			t.Status = r.Action
			t.Elapsed = time.Duration(r.Elapsed) * time.Millisecond
		default:
			return fmt.Errorf("unknown action - %q", r.Action)
		}
	}

	t.Message = buf.String()
	if isEmptyTime(t.Time) {
		for _, r := range t.records {
			if !isEmptyTime(r.Time) {
				t.Time = r.Time
				break
			}
		}
	}
	return nil
}

func isEmptyTime(t time.Time) bool {
	return t.Equal(zeroTime)
}
