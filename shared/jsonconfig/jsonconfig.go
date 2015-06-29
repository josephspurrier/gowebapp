package jsonconfig

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Parser must implement ParseJSON
type Parser interface {
	ParseJSON([]byte) error
}

// Load the JSON config file
func Load(configFile string, p Parser) {
	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Parse the config
	if err := p.ParseJSON(jsonBytes); err != nil {
		log.Printf("Could not parse %q: %v", configFile, err)
		os.Exit(2)
	}
}
