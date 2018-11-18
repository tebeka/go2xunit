package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Status is test status
type Status string

// Test statuses
const (
	UnknownStatus Status = "Unknown"
	FailedStatus  Status = "Failed"
	SkippedStatus Status = "Skipped"
	PassedStatus  Status = "Passed"
)

// Test data structure
type Test struct {
	Name     string
	Package  string
	Time     time.Time
	Status   Status
	Children []*Test
	Message  string

	records []*Record
}

// Record is test in JSON output
type Record struct {
	Name    string
	Time    time.Time
	Action  string
	Package string
	Output  string
}

type key struct {
	pkg  string
	name string
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
			return nil, fmt.Errorf("%d: error: %s", scan.LineNum(), err)
		}
		k := key{r.Package, r.Name}
		t, ok := tests[k]
		if !ok {
			t = &Test{
				Name:    r.Name,
				Package: r.Package,
			}
		}
		t.records = append(t.records, r)
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return tests, nil
}

func assembleTests(tests map[key]*Test) (*Test, error) {
	return nil, nil
}

// Len return number of sub tests
func (t *Test) Len() int {
	return len(t.Children)
}

// NumStatus returns the number of tests and subtests that have status
func (t *Test) NumStatus(status Status) int {
	n := 0
	if t.Status == status {
		n = 1
	}

	for _, st := range t.Children {
		n += st.NumStatus(status)
	}

	return n
}
