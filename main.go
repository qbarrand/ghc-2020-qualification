package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
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

	in, err := pkg.ParseInput(input, logger)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %s: %v\n", input, err)
		panic(msg)
	}

	ints := make([]int, len(in.Libraries))

	for i, l := range in.Libraries {
		ints[i] = l.SignupTime
	}

	std := stdDev(ints)

	logger.Printf("stddev: %f", std)

	day := 0

	i := 0

	type signedUpLibrary struct {
		libID int
		books []*pkg.Book
	}

	signedUpLibraries := make([]*signedUpLibrary, 0)

	for day < in.DaysForScanning {
		remainingDays := in.DaysForScanning - day

		selected := pkg.GetBestLibrary(in.Libraries, remainingDays, std)

		// Maybe all libraries are selected already?
		if selected == nil {
			break
		}

		selected.MarkAsSelected()

		if in.DaysForScanning-day >= selected.SignupTime {
			day += selected.SignupTime
			i++

			sul := &signedUpLibrary{
				libID: selected.ID,
				books: selected.GetBooksToBeSent(remainingDays),
			}

			pkg.MarkBooksAsScanned(sul.books)

			if len(sul.books) > 0 {
				// Only mention libraries from which we'll scan books in the output
				signedUpLibraries = append(signedUpLibraries, sul)
			}
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

	for _, sul := range signedUpLibraries {
		fmt.Fprintf(fd, "%d %d\n", sul.libID, len(sul.books))

		for _, book := range sul.books {
			fmt.Fprintf(fd, "%d ", book.ID)
		}

		fmt.Fprint(fd, "\n")
	}
}

func stdDev(ints []int) float64 {
	sum := 0
	n := len(ints)

	for _, i := range ints {
		sum += i
	}

	average := float64(sum) / float64(n)

	var top float64 = 0

	for _, i := range ints {
		top += math.Pow(float64(i)-average, 2)
	}

	top /= float64(n)

	return math.Sqrt(top)
}
