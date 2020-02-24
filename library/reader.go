package library

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// readingInfo stores information necessary for the next reading
type readingInfo struct {
	readText       string
	lastReadingPos int
}

// Reader stores information about books read
type Reader struct {
	ID                      int
	lib                     *Library
	wastepaperRead          map[Wastepaper]*readingInfo
	wastepaperFinishReading []Wastepaper
	lastWastepaper          Wastepaper
	winner                  chan *Reader
	once                    *sync.Once
	ctx                     context.Context
}

// NewReader creates new Reader
func NewReader(ctx context.Context, ID int, lib *Library, winner chan *Reader, once *sync.Once) *Reader {
	if ctx == nil {
		return nil
	}
	if lib == nil {
		return nil
	}
	if winner == nil {
		return nil
	}
	if once == nil {
		return nil
	}
	return &Reader{
		ID:                      ID,
		lib:                     lib,
		wastepaperRead:          make(map[Wastepaper]*readingInfo),
		wastepaperFinishReading: make([]Wastepaper, 0, len(lib.wastepaper)),
		winner:                  winner,
		once:                    once,
		ctx:                     ctx,
	}
}

// Read randomly selects the position of the book in the library and reads the word from there
func (r *Reader) Read() {
	var (
		wastepaperPos        int
		targetWastepaper     Wastepaper
		targetWastepaperInfo *readingInfo
		targetWastepaperText string
		word                 string
		wordEndByte          int
	)

	rand.Seed(time.Now().UnixNano())

	wastepaperPos = rand.Intn(len(r.lib.wastepaper)) + 1

	targetWastepaper = r.lib.GetWastepaper(wastepaperPos)
	if targetWastepaper == nil {
		return
	}

	if targetWastepaper == r.lastWastepaper {
		r.lib.PutWastepaper(targetWastepaper)
		return
	}

	r.lastWastepaper = targetWastepaper

	if _, ok := r.wastepaperRead[targetWastepaper]; !ok {
		r.wastepaperRead[targetWastepaper] = &readingInfo{
			readText:       "",
			lastReadingPos: 0,
		}
	}

	targetWastepaperInfo = r.wastepaperRead[targetWastepaper]

	targetWastepaperText = targetWastepaper.GetText()

	r.lib.PutWastepaper(targetWastepaper)

	if len(targetWastepaperText) == targetWastepaperInfo.lastReadingPos {
		return
	}

	word, wordEndByte = countWordBytes(targetWastepaperText[targetWastepaperInfo.lastReadingPos:])

	targetWastepaperInfo.lastReadingPos += wordEndByte
	targetWastepaperInfo.readText += word

	r.wastepaperRead[targetWastepaper] = targetWastepaperInfo

	if len(targetWastepaperText) == targetWastepaperInfo.lastReadingPos {
		r.wastepaperFinishReading = append(r.wastepaperFinishReading, targetWastepaper)

		if len(r.wastepaperFinishReading) == len(r.lib.wastepaper) {
			r.once.Do(func() {
				r.winner <- r
			})
		}
	}
}

func countWordBytes(text string) (string, int) {
	var (
		wordEndByte int
		word        string
	)

	wordEndByte = strings.Index(text, " ")
	if wordEndByte == -1 {
		word = text[:]
		wordEndByte = len(text)
		return word, wordEndByte
	}
	word = text[:wordEndByte+1]
	wordEndByte++
	return word, wordEndByte
}

// WastepaperReadInfo displays information about wastepaper and the text read by the Reader
func (r *Reader) WastepaperReadInfo() {
	var rInfo *readingInfo
	for wp, info := range r.wastepaperRead {
		rInfo = info
		wp.Info()
		fmt.Printf("Read text: %s\n\n", rInfo.readText)
	}
}
