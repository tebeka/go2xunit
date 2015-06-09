package main

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
