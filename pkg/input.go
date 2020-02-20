package pkg

import (
	"fmt"
	"log"
	"os"
	"sort"
)

type input struct {
	Books           []*Book
	Libraries       []*Library
	DaysForScanning int
}

type Book struct {
	ID    int
	Score int
}

func (b *Book) String() string {
	return fmt.Sprintf("Book %d has a score of %d", b.ID, b.Score)
}

type Library struct {
	Books       []*Book
	BooksPerDay int
	ID          int
	SignupTime  int

	totalScore int
}

func (l *Library) GetBooksToBeSent(days int) []int {
	// Copy the library's books and sort the copied slice
	sortedBooks := make([]*Book, len(l.Books))

	for i, b := range l.Books {
		sortedBooks[i] = b
	}

	sort.Slice(sortedBooks, func(i, j int) bool {
		return sortedBooks[i].Score > sortedBooks[j].Score
	})

	booksToScan := days * l.BooksPerDay

	//fmt.Printf("DBG days: %d\n", days)
	//fmt.Printf("DBG BooksPerDay: %d\n", l.BooksPerDay)
	//fmt.Printf("DBG booksToScan: %d\n", booksToScan)

	if booksToScan > len(sortedBooks) {
		booksToScan = len(sortedBooks)
	}

	ids := make([]int, booksToScan)

	for i := 0; i < booksToScan; i++ {
		ids[i] = sortedBooks[i].ID
	}

	return ids
}

func (l *Library) String() string {
	return fmt.Sprintf(
		"Library %d holds %d books, takes %d days to signup, can send %d books per day, total score %d",
		l.ID,
		len(l.Books),
		l.SignupTime,
		l.BooksPerDay,
		l.totalScore,
	)
}

func (l *Library) TotalScore() int {
	return l.totalScore
}

func ParseInput(filename string, logger *log.Logger) (*input, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	in := input{}

	var (
		nBooks     int
		nLibraries int
	)

	if _, err := fmt.Fscanf(fd, "%d %d %d", &nBooks, &nLibraries, &in.DaysForScanning); err != nil {
		return nil, fmt.Errorf("could not read the header: %v", err)
	}

	in.Books = make([]*Book, nBooks)

	for i := 0; i < nBooks; i++ {
		b := Book{ID: i}

		_, err := fmt.Fscanf(fd, "%d", &b.Score)
		if err != nil {
			panic(err)
		}

		in.Books[i] = &b
	}

	logger.Printf("Read %d books", len(in.Books))

	in.Libraries = make([]*Library, nLibraries)

	for i := 0; i < nLibraries; i++ {
		l := Library{ID: i}

		var booksInThisLibrary int

		_, err := fmt.Fscanf(fd, "%d %d %d", &booksInThisLibrary, &l.SignupTime, &l.BooksPerDay)
		if err != nil {
			panic(err)
		}

		l.Books = make([]*Book, booksInThisLibrary)

		for j := 0; j < booksInThisLibrary; j++ {
			var bookID int

			_, err := fmt.Fscanf(fd, "%d", &bookID)
			if err != nil {
				panic(err)
			}

			b := in.Books[bookID]

			l.Books[j] = b
			l.totalScore += b.Score
		}

		in.Libraries[i] = &l

		logger.Print(l.String())
	}

	logger.Printf("Read %d libraries", len(in.Libraries))

	return &in, nil
}
