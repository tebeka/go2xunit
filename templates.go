// Genrated by gentmpl.go

package main

var junit = `<!DOCTYPE xml>
<testsuite
    name="go2xunit"
    tests="{{.Count}}"
    errors="0"
    failures="{{index .Stats "fail"}}"
    skip="{{index .Stats "skip"}}">
{{range  $test := .Children}}
    <testcase
	classname="{{$test.Package | escape}}"
	name="{{$test.Name | escape}}"
	time="{{$test.Elapsed.Seconds}}">
    {{if eq $test.Status "skip" }}
	<skipped/>
    {{end}}
    {{if eq $test.Status "fail" }}
	<failure type="go.error" message="error">
        <![CDATA[{{$test.Message}}]]>
        </failure>
    {{end}}
    </testcase>
{{end}}
</testsuite>
`

var xunit = `<!DOCTYPE xml>
<assembly name="{{.Name | escape}}"
          run-date="{{.Time}}" run-time="{{.Time}}"
          configFile="none"
          time="{{.Time}}"
          total="{{.Count}}"
          passed="{{index .Stats "passed"}}"
          failed="{{index .Stats "failed"}}"
          skipped="{{index .Stats "skipped"}}"
          environment="n/a"
          test-framework="golang">
    <class time="{{.Time}}" name="{{.Name | escape}}"
  	     total="{{.Count}}"
          passed="{{index .Stats "passed"}}"
          failed="{{index .Stats "failed"}}"
          skipped="{{index .Stats "skipped"}}"
    >
    {{range  $test := $.Children}}
        <test name="{{$test.Name | escape}}"
          type="test"
          method="{{$test.Name | escape}}"
          result={{$test.Status}}
          time="{{$test.Time}}">
        {{if eq $test.Status "fail" }}  <failure exception-type="go.error">
             <message><![CDATA[{{$test.Message}}]]></message>
      	  </failure>
      	{{end}}</test>
    {{end}}
    </class>
</assembly>
`

var internalTemplates = map[string]string{
	"junit": junit,
	"xunit": xunit,
}
