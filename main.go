package main

import (
	"flag"
	"io/ioutil"
	"log"
	"multy-go/decoder"
	"multy-go/encoder"
	"multy-go/parser"
	"multy-go/variables"
	"os"
	"path/filepath"
	"time"
)

func main() {
	start := time.Now()
	log.SetFlags(log.Lshortfile)

	var commandLineVars variables.CommandLineVariables

	flag.Var(&commandLineVars, "var", "Variables to be passed to configuration")
	outputFile := flag.String("output", "main.tf", "Name of the output file")
	flag.Parse()
	configFiles := flag.Args()

	if len(configFiles) == 0 {
		log.Println("scanning current directory for .mt files")
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatalf("error while reading current directory: %s", err.Error())
		}
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".mt" {
				log.Println("config file found:", file.Name())
				configFiles = append(configFiles, file.Name())
			}
		}
	}

	if len(configFiles) == 0 {
		log.Fatalf("no .mt config files found")
	}

	p := parser.Parser{CliVars: commandLineVars}
	parsedConfig := p.Parse(configFiles)

	r := decoder.Decode(parsedConfig)

	hclOutput := encoder.Encode(r)

	d1 := []byte(hclOutput)
	err := os.WriteFile(*outputFile, d1, 0644)
	if err != nil {
		log.Fatalf("error creating output file: %s", err.Error())
	}

	// needs terraform cli to run
	//cmd := exec.Command("terraform", "fmt", *outputFile)
	//err = cmd.Run()
	//if err != nil {
	//	log.Fatalf("error formatting output file: %s", err.Error())
	//}

	log.Printf("multy finished running: %s\n\n", time.Since(start).Round(time.Second))
	log.Println("output file:", *outputFile)
}
