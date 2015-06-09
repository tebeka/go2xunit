package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

var goCheckFiles []string = []string{
	"gocheck-pass.out",
}

var goTestFiles []string = []string{
	"gotest.out",
	"gotest-0.out",
	"gotest-datarace.out",
	"gotest-empty.out",
	"gotest-fail.out",
	"gotest-log.out",
	"gotest-multi.out",
	"gotest-multierror.out",
	"gotest-nofiles.out",
	"gotest-num.out",
	"gotest-panic.out",
	"gotest-pass.out",
	"gotest-testify-suite.out",
}

func Test_XMLOuptutGoCheckXUnit(t *testing.T) {
	for _, filename := range goCheckFiles {
		suites, err := gc_Parse(getInputData(filename))
		checkError(err)
		generateAndTestXMLXUnit(t, suites, filename)
	}
}

func Test_XMLOuptutGoCheckXUnitNet(t *testing.T) {
	for _, filename := range goCheckFiles {
		suites, err := gc_Parse(getInputData(filename))
		checkError(err)
		generateAndTestXMLXUnitNet(t, suites, filename)
	}
}

func Test_XMLOuptutGoTestXUnit(t *testing.T) {
	for _, filename := range goTestFiles {
		suites, err := gt_Parse(getInputData(filename))
		checkError(err)
		generateAndTestXMLXUnit(t, suites, filename)
	}
}

func Test_XMLOuptutGoTestXunitNet(t *testing.T) {
	for _, filename := range goTestFiles {
		suites, err := gt_Parse(getInputData(filename))
		checkError(err)
		generateAndTestXMLXUnitNet(t, suites, filename)
	}
}

func getInputData(filename string) io.Reader {
	file, err := getInput("data" + string(os.PathSeparator) + filename)
	checkError(err)
	return file
}

func getOutputData(outType string, filename string) io.Reader {
	file, err := getInput("xml" + string(os.PathSeparator) + outType + string(os.PathSeparator) + filename + ".xml")
	checkError(err)
	return file
}

func generateXML(suites []*Suite, filename string, xmlTemplate string) []byte {
	r, w, _ := os.Pipe()
	writeXML(suites, w, xmlTemplate)
	w.Close()
	xml, err := ioutil.ReadAll(r)
	checkError(err)
	return xml
}

func generateAndTestXMLXUnit(t *testing.T, suites []*Suite, filename string) {
	expected, err := ioutil.ReadAll(getOutputData("xunit", filename))
	checkError(err)

	var xmlTemplate = xunitTemplate
	if len(suites) > 1 {
		xmlTemplate = multiTemplate
	}

	actual := generateXML(suites, filename, xmlTemplate)
	if !bytes.Equal(expected, actual) {
		t.Errorf("xUnit XML output %s differs from expected", filename)
	}
}

func generateAndTestXMLXUnitNet(t *testing.T, suites []*Suite, filename string) {
	// run-date="2015-06-05" run-time="18:34:41"
	runTime = time.Date(2015, time.June, 5, 18, 34, 41, 0, time.UTC)

	expected, err := ioutil.ReadAll(getOutputData("xunit.net", filename))
	checkError(err)
	actual := generateXML(suites, filename, xunitNetTemplate)
	if !bytes.Equal(expected, actual) {
		t.Errorf("xUnit.net XML output %s differs from expected", filename)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
