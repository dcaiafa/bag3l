package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/dcaiafa/bag3l/internal/errlogger"
	parser "github.com/dcaiafa/bag3l/internal/stub/parser2"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	defFile := flag.Arg(0)
	err := run(defFile)
	if err != nil {
		log.Fatal(err)
	}
}

func run(defFile string) error {
	defFileBytes, err := os.ReadFile(defFile)
	if err != nil {
		return err
	}

	errLogger := errlogger.NewErrLoggerBase(&errlogger.ConsoleErrLogger{})

	unit, err := parser.Parse(defFile, defFileBytes, errLogger)
	if err != nil {
		return err
	}

	outputFilename := filepath.Base(defFile)
	outputFilename = outputFilename + ".go"
	outputFilename = filepath.Join(filepath.Dir(defFile), outputFilename)

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	generator := NewGenerator(unit, outputFile)

	return generator.Run()
}
