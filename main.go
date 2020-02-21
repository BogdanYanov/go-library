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

func generateSimpleString() (string, error) {
	simpleText := "Lorem ipsum dolor sit amet, nostrud scribentur eam ad. Ea utamur aliquip corpora usu. " +
		"Has in laboramus definitionem, sit ei fugit adipisci disputationi. Ea eos volumus theophrastus. Timeam aeterno interpretaris sed id, " +
		"his no eros erroribus sadipscing. Sed cu apeirian persecuti, cum novum voluptatibus ne, ad nec quaestio adolescens. Cibo altera eos id. " +
		"Luptatum moderatius vel no, qui at cetero intellegam, nobis audire ius te."
	rd := strings.NewReader(simpleText)
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, rand.Intn(rd.Len()-16)+16)
	_, err := rd.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func main() {

	var (
		defaultBooksNum       = 10
		defaultReadersNum     = 5
		defaultPublishingYear = 1976
		lib                   *library.Library
		textSlice             []string
		_                     chan *library.Reader
		_                     *library.Reader
	)

	ctx, cancel := context.WithCancel(context.Background())

	lib = library.NewLibrary()

	textSlice = make([]string, 0, defaultBooksNum)

	for i := 0; i < defaultBooksNum; i++ {
		simpleText, err := generateSimpleString()
		if err != nil {
			log.Println(err)
			i--
			continue
		}
		textSlice = append(textSlice, simpleText)
	}

	for i := 0; i < defaultBooksNum; i++ {
		lib.AddWastepaper(library.NewBook(fmt.Sprintf("Author %d", i), defaultPublishingYear, textSlice[i]))
	}

	rd := library.NewReadersDirector(ctx, lib, defaultReadersNum)

	_ = rd.StartWork()

	cancel()

	//winner = <-winnerCh

	//cancel()

	//log.Printf("Winner - goroutine#%d\n", winner.ID)

	//winner.WastepaperReadInfo()
}
