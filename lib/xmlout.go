package lib

// XML output
import (
	"fmt"
	"io"
	"strconv"
	"text/template"
	"time"
)

const (
	xmlDeclaration = `<?xml version="1.0" encoding="utf-8"?>`

	// XUnitTemplate is XML template for xunit style reporting
	XUnitTemplate string = `
{{range $suite := .Suites}}  <testsuite name="{{.Name}}" tests="{{.Len}}" errors="0" failures="{{.NumFailed}}" skip="{{.NumSkipped}}">
{{range  $test := $suite.Tests}}    <testcase classname="{{$suite.Name}}" name="{{$test.Name}}" time="{{$test.Time}}">
{{if eq $test.Status $.Skipped }}      <skipped/> {{end}}
{{if eq $test.Status $.Failed }}      <failure type="go.error" message="error">
        <![CDATA[{{$test.Message}}]]>
      </failure>{{end}}    </testcase>
{{end}}  </testsuite>
{{end}}`

	// XMLMultiTemplate is template when we have multiple suites
	XMLMultiTemplate string = `
<testsuites>` + XUnitTemplate + `</testsuites>
`

	// XUnitNetTemplate is XML template for xunit.net
	// see https://xunit.codeplex.com/wikipage?title=XmlFormat
	XUnitNetTemplate string = `
<assembly name="{{.Assembly}}"
          run-date="{{.RunDate}}" run-time="{{.RunTime}}"
          configFile="none"
          time="{{.Time}}"
          total="{{.Len}}"
          passed="{{.NumPassed}}"
          failed="{{.NumFailed}}"
          skipped="{{.NumSkipped}}"
          environment="n/a"
          test-framework="golang">
{{range $suite := .Suites}}
    <class time="{{.Time}}" name="{{.Name}}"
  	     total="{{.Len}}"
  	     passed="{{.NumPassed}}"
  	     failed="{{.NumFailed}}"
  	     skipped="{{.NumSkipped}}">
{{range  $test := $suite.Tests}}
        <test name="{{$test.Name}}"
          type="test"
          method="{{$test.Name}}"
          result={{if eq $test.Status $.Skipped }}"Skip"{{else if eq $test.Status $.Failed }}"Fail"{{else if eq $test.Status $.Passed }}"Pass"{{end}}
          time="{{$test.Time}}">
        {{if eq $test.Status $.Failed }}  <failure exception-type="go.error">
             <message><![CDATA[{{$test.Message}}]]></message>
      	  </failure>
      	{{end}}</test>
{{end}}
    </class>
{{end}}
</assembly>
`
)

// TestResults is passed to XML template
type TestResults struct {
	Suites     []*Suite
	Assembly   string
	RunDate    string
	RunTime    string
	Time       string
	Len        int
	NumPassed  int
	NumFailed  int
	NumSkipped int

	Skipped Status
	Passed  Status
	Failed  Status
}

// calcTotals calculates grand total for all suites
func (r *TestResults) calcTotals() {
	totalTime, _ := strconv.ParseFloat(r.Time, 64)
	for _, suite := range r.Suites {
		r.NumPassed += suite.NumPassed()
		r.NumFailed += suite.NumFailed()
		r.NumSkipped += suite.NumSkipped()

		suiteTime, _ := strconv.ParseFloat(suite.Time, 64)
		totalTime += suiteTime
		r.Time = fmt.Sprintf("%.3f", totalTime)
	}
	r.Len = r.NumPassed + r.NumSkipped + r.NumFailed
}

// WriteXML exits xunit XML of tests to out
func WriteXML(suites []*Suite, out io.Writer, xmlTemplate string, testTime time.Time) {
	testsResult := TestResults{
		Suites:   suites,
		Assembly: suites[len(suites)-1].Name,
		RunDate:  testTime.Format("2006-01-02"),
		RunTime:  testTime.Format("15:04:05"),
		Skipped:  Skipped,
		Passed:   Passed,
		Failed:   Failed,
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
