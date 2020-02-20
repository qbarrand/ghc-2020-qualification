package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/qbarrand/ghc-2019-qualification/pkg"
)

func main() {
	fDebug := flag.Bool("debug", false, "enable debug logging")
	fOutDir := flag.String("outdir", "out", "the directory in which the output files should be stored")

	flag.Parse()

	var loggerOut io.Writer

	if !*fDebug {
		loggerOut = ioutil.Discard
	} else {
		loggerOut = os.Stdout
	}

	log.SetOutput(loggerOut)

	log.Printf("Storing the ouputs in %s", *fOutDir)

	var wg sync.WaitGroup
	wg.Add(flag.NArg())

	for _, input := range flag.Args() {
		go func(i string) {
			process(
				i,
				*fOutDir,
				log.New(loggerOut, fmt.Sprintf("%s | ", i), 0),
			)

			wg.Done()
		}(input)
	}

	log.Print("Waiting for the goroutines")

	wg.Wait()
}

func process(input, outDir string, logger *log.Logger) {
	logger.Printf("Starting the goroutine")

	in, err := pkg.ParseInput(input)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %s: %v\n", input, err)
		panic(msg)
	}

	_ = in

	outFilename := filepath.Join(outDir, input)

	logger.Print("Writing " + outFilename)
}
