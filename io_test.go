package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var runTime = time.Date(2015, time.June, 5, 18, 34, 41, 0, time.UTC)

var goCheckFiles []string = []string{
	"gocheck-pass.out",
	"gocheck-fail.out",
	"gocheck-panic.out",
	"gocheck-nofiles.out",
	"gocheck-empty.out",
	"gocheck-setup-miss.out",
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
	"gotest-1.5.out",
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
	file, err := getInput(filepath.Join("data", filename))
	checkError(err)
	return file
}

func getOutputData(outType string, filename string) io.Reader {
	path := filepath.Join("xml", outType, fmt.Sprintf("%s.xml", filename))
	file, err := getInput(path)
	checkError(err)
	return file
}

func generateXML(suites []*Suite, filename string, xmlTemplate string) []byte {
	if len(suites) == 0 {
		return []byte("error: no tests found")
	}
	r, w, _ := os.Pipe()
	writeXML(suites, w, xmlTemplate, runTime)
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
