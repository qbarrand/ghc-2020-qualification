package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
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

	in, err := pkg.ParseInput(input, logger)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %s: %v\n", input, err)
		panic(msg)
	}

	sortedLibraries := make([]*pkg.Library, len(in.Libraries))

	for i, l := range in.Libraries {
		sortedLibraries[i] = l
	}

	sort.Slice(sortedLibraries, func(i, j int) bool {
		return sortedLibraries[i].TotalScore() > sortedLibraries[j].TotalScore()
	})

	day := 0

	i := 0

	daysToScan := make(map[int]int)

	signedUpLibraries := make([]*pkg.Library, 0)

	//curr := sortedLibraries[i]

	for day < in.DaysForScanning && i < len(sortedLibraries) {
		curr := sortedLibraries[i]

		if in.DaysForScanning-day >= curr.SignupTime {
			day += curr.SignupTime
			i++

			signedUpLibraries = append(signedUpLibraries, curr)
			daysToScan[curr.ID] = in.DaysForScanning - day
		} else {
			break
		}
	}

	outFilename := filepath.Join(outDir, filepath.Base(input))

	logger.Print("Writing " + outFilename)

	fd, err := os.Create(outFilename)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fmt.Fprintf(fd, "%d\n", len(signedUpLibraries))

	for _, l := range signedUpLibraries {
		nBooks := l.GetBooksToBeSent(daysToScan[l.ID])

		fmt.Fprintf(fd, "%d %d\n", l.ID, len(nBooks))

		for _, bookID := range nBooks {
			fmt.Fprintf(fd, "%d ", bookID)
		}

		fmt.Fprint(fd, "\n")
	}
}
