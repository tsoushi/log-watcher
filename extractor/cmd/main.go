package main

import (
	"flag"
	"logwatcher/extractor"
)

var (
	text         = flag.String("text", "", "Text")
	regexString  = flag.String("regexString", "", "Regex string")
	outputFormat = flag.String("outputFormat", "", "Output format")
)

func main() {
	flag.Parse()
	outputs, err := extractor.ExtractAndReplaceText(*text, *regexString, *outputFormat)
	if err != nil {
		panic(err)
	}
	for _, output := range outputs {
		println(output)
	}
}
