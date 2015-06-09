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

// calculates the totals to populate xunit.net assembly attributes
func calcTotals(r *TestResults) {
	for _, suite := range r.Suites {
		r.Passed += suite.NumPassed()
		r.Failed += suite.NumFailed()
		r.Skipped += suite.NumSkipped()

		totalTime, _ := strconv.ParseFloat(r.Time, 64)
		suiteTime, _ := strconv.ParseFloat(suite.Time, 64)
		totalTime += suiteTime
		r.Time = fmt.Sprintf("%.3f", totalTime)

		// last suite is the go package name
		r.Assembly = suite.Name
	}
	r.Total = r.Passed + r.Skipped + r.Failed
}
