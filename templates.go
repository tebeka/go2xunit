// Genrated by /tmp/go-build722822768/b001/exe/gentmpl

package main

var junit = `<testsuite
    name="{{.Name | escape}}"
    tests="{{.Count}}"
    errors="0"
    failures="{{.Stats["fail"]}}"
    skip="{{.Stats["fail"]}}">
{{range  $test := .Children}}
    <testcase
	classname="{{$suite.Name | escape}}"
	name="{{$test.Name | escape}}"
	time="{{$test.Time}}">
    {{if eq $test.Status "skip" }}
	<skipped/>
    {{end}}
    {{if eq $test.Status "fail" }}
	<failure type="go.error" message="error">
        <![CDATA[{{$test.Message}}]]>
        </failure>
    {{end}
    </testcase>
{{end}}
</testsuite>
`

var Templates = map[string]string{
	"junit": junit,
}
