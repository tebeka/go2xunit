// Genrated by /tmp/go-build621141117/b001/exe/gentmpl

package main

var junit = `<testsuite
    name="go2xunit"
    tests="{{.Count}}"
    errors="0"
    failures="{{index .Stats "fail"}}"
    skip="{{index .Stats "skip"}}">
{{range  $test := .Children}}
    <testcase
	classname="{{$test.Package | escape}}"
	name="{{$test.Name | escape}}"
	time="{{$test.Time}}">
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

var Templates = map[string]string{
	"junit": junit,
}
