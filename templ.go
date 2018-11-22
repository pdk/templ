package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

var (
	inFile   string
	outFile  string
	dataFile string
)

func init() {
	flag.StringVar(&inFile, "template", "", "name of file containing template")
	flag.StringVar(&inFile, "t", "", "name of file containing template")
	flag.StringVar(&outFile, "write", "", "output file name")
	flag.StringVar(&outFile, "w", "", "output file name")
	flag.StringVar(&dataFile, "data", "", "name of file containing data (json)")
	flag.StringVar(&dataFile, "d", "", "name of file containing data (json)")
}

func main() {
	flag.Parse()

	checkInvocation()

	data := readData(dataFile)

	templ := readTemplate(inFile)

	result := executeTemplate(templ, data)

	writeResult(outFile, result)
}

func writeResult(outFile string, result []byte) {

	if outFile == "" {
		_, err := os.Stdout.Write(result)
		if err != nil {
			log.Fatalf("cannot write to stdout: %s", err)
		}

		return
	}

	err := ioutil.WriteFile(outFile, result, os.ModePerm)
	if err != nil {
		log.Fatalf("cannot write result %s: %s", outFile, err)
	}
}

func checkInvocation() {
	if dataFile == "" {
		flag.Usage()
		fmt.Fprintf(flag.CommandLine.Output(), "-d/-data required")
		os.Exit(1)
	}
}

func readData(dataFile string) interface{} {

	dataBytes, err := ioutil.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("cannot read datafile: %s", err)
	}

	var data interface{}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		log.Fatalf("cannot parse data from %s: %s", dataFile, err)
	}

	return data
}

func readTemplate(inFile string) *template.Template {

	var inputBytes []byte
	var err error

	if inFile == "" {
		inputBytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		inputBytes, err = ioutil.ReadFile(inFile)
	}
	if err != nil {
		log.Fatalf("cannot read input %s: %s", inFile, err)
	}

	templ, err := template.New("templ").Parse(string(inputBytes))
	if err != nil {
		log.Fatalf("cannot parse template %s: %s", inFile, err)
	}

	return templ
}

func executeTemplate(templ *template.Template, data interface{}) []byte {

	// execute the template into a buffer so we don't touch the output file
	// until we're sure we have a result.
	resultBuf := &bytes.Buffer{}
	err := templ.Execute(resultBuf, data)
	if err != nil {
		log.Fatalf("cannot execute template %s: %s", inFile, err)
	}

	return resultBuf.Bytes()
}
