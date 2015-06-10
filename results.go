package main

import (
	"fmt"
	"strconv"
)

type Test struct {
	Name, Time, Message string
	Failed              bool
	Skipped             bool
	Passed              bool
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
