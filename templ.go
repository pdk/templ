package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/pdk/qkjson/parser"
)

var (
	inFile     string
	outFile    string
	dataFile   string
	removeJSON bool
)

func init() {
	flag.StringVar(&inFile, "input", "", "name of file containing template")
	flag.StringVar(&inFile, "i", "", "name of file containing template")
	flag.StringVar(&outFile, "output", "", "name of file to write")
	flag.StringVar(&outFile, "o", "", "name of file to write")
	flag.StringVar(&dataFile, "json", "", "name of file containing json data")
	flag.StringVar(&dataFile, "j", "", "name of file containing json data")
	flag.BoolVar(&removeJSON, "rm", false, "remove the json file afterwards")
	flag.BoolVar(&removeJSON, "r", false, "remove the json file afterwards")
}

func main() {
	flag.Parse()

	data := readData(dataFile, flag.Args())

	templ := readTemplate(inFile)

	result := executeTemplate(templ, data)

	writeResult(outFile, result)

	if !removeJSON {
		return
	}

	err := os.Remove(dataFile)
	if err != nil {
		log.Fatalf("cannot remove data file: %s", err)
	}
}

func readData(dataFile string, args []string) interface{} {

	if dataFile != "" && len(args) > 0 {
		log.Fatalf("cannot use data from file and from command line simultaneously")
	}

	if dataFile == "" {
		return parser.ParseArgs(args)
	}

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
