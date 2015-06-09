package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"time"
)

// Time when the test was run
var runTime time.Time

// writeXML exits xunit XML of tests to out
func writeXML(suites []*Suite, out io.Writer, xmlTemplate string) {
	testsResult := TestResults{
		Suites:  suites,
		RunDate: runTime.Format("2006-01-02"),
		RunTime: fmt.Sprintf("%02d:%02d:%02d",
			runTime.Hour(),
			runTime.Minute(),
			runTime.Second()),
	}
	calcTotals(&testsResult)
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

// getInput return input *io.File from file name, if file name is - it will
// return os.Stdin
func getInput(filename string) (*os.File, error) {
	if filename == "-" || filename == "" {
		return os.Stdin, nil
	}
	return os.Open(filename)
}

// getInput return output *io.File from file name, if file name is - it will
// return os.Stdout
func getOutput(filename string) (*os.File, error) {
	if filename == "-" || filename == "" {
		return os.Stdout, nil
	}
	return os.Create(filename)
}

// getIO returns input and output streams from file names
func getIO(inputFile, outputFile string) (io.Reader, io.Writer, error) {
	input, err := getInput(inputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %s for reading: %s", inputFile, err)
	}

	// let's assume test report modification time is also the time when test was run
	statinfo, err := input.Stat()
	runTime = statinfo.ModTime()

	output, err := getOutput(outputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %s for writing: %s", outputFile, err)
	}

	return input, output, nil
}
