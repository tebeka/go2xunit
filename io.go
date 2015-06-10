package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
	"time"
)

// Time when the test was run
var runTime time.Time

// getInput return input *io.File from file name, if file name is - it will
// return os.Stdin
func getInput(filename string) (*os.File, error) {
	if isEmpty(filename) {
		return os.Stdin, nil
	}
	return os.Open(filename)
}

// getInput return output *io.File from file name, if file name is - it will
// return os.Stdout
func getOutput(filename string) (*os.File, error) {
	if isEmpty(filename) {
		return os.Stdout, nil
	}
	return os.Create(filename)
}

// returns true if no filename or "-" was specified, false otherwise
func isEmpty(filename string) bool {
	return filename == "-" || filename == ""
}

// getIO returns input and output streams from file names
func getIO(inputFile, outputFile string) (io.Reader, io.Writer) {
	input, err := getInput(inputFile)
	checkFatal(err)

	output, err := getOutput(outputFile)
	checkFatal(err)

	setRunTimeFrom(input)

	return input, output
}

// set test execution time from file date
func setRunTimeFrom(file *os.File) {
	statinfo, err := file.Stat()
	checkFatal(err)
	runTime = statinfo.ModTime()
}

// writeXML exits xunit XML of tests to out
func writeXML(suites []*Suite, out io.Writer, xmlTemplate string) {
	results := TestResults{
		Suites:   suites,
		Assembly: suites[len(suites)-1].Name,
		RunDate:  runTime.Format("2006-01-02"),
		RunTime:  fmt.Sprintf("%02d:%02d:%02d", runTime.Hour(), runTime.Minute(), runTime.Second()),
	}
	results.calcTotals()
	t := template.New("test template")

	t, err := t.Parse(xmlDeclaration + xmlTemplate)
	checkFatal(err)

	err = t.Execute(out, results)
	checkFatal(err)
}

func checkFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
