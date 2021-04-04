package utils

import (
	"plugin"
	"log"
)

func LoadPrepPlugin(filename string) func(string, string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot open plugin %v", filename)
	}
	xprep, err := p.Lookup("Preprocess")
	if err != nil {
		log.Fatalf("cannot find Map in %v", filename)
	}
	prepf := xprep.(func(string, string))

	return prepf
}



