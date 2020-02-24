package main

import (
	"context"
	"fmt"
	"github.com/BogdanYanov/go-library/library"
	"log"
	"math/rand"
	"strings"
	"time"
)

func generateSimpleString() string {
	simpleText := "To wait for multiple goroutines to finish, we can use a wait group. This is the function we’ll run in every goroutine. Note that a WaitGroup must be passed to functions by pointer."
	rand.Seed(time.Now().UnixNano())
	result := simpleText[:rand.Intn(len(simpleText)-16)+16+1]
	wordEndByte := strings.LastIndex(result, " ")
	if wordEndByte == -1 {
		return result
	}
	return result[:wordEndByte]
}

func main() {

	var (
		defaultBooksNum       = 10
		defaultReadersNum     = 5
		defaultPublishingYear = 1976
		lib                   *library.Library
		textSlice             []string
		winnerCh              chan *library.Reader
		winner                *library.Reader
	)

	ctx, cancel := context.WithCancel(context.Background())

	lib = library.NewLibrary()

	textSlice = make([]string, 0, defaultBooksNum)

	textSlice = append(textSlice, "abcd")
	for i := 1; i < defaultBooksNum; i++ {
		simpleText := generateSimpleString()
		textSlice = append(textSlice, simpleText)
	}

	for i := 0; i < defaultBooksNum; i++ {
		lib.AddWastepaper(library.NewBook(fmt.Sprintf("Author %d", i), defaultPublishingYear, textSlice[i]))
	}

	rd := library.NewReadersDirector(ctx, lib, defaultReadersNum)

	winnerCh = rd.StartWork()

	winner = <-winnerCh

	cancel()

	log.Printf("Winner - goroutine#%d\n", winner.ID)

	winner.WastepaperReadInfo()
}
