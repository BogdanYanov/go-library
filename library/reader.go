package library

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
func (r *Reader) Read() error {
	var (
		wastepaperPos        int
		targetWastepaper     Wastepaper
		lastReadingPos       int
		targetWastepaperInfo *readingInfo
	)

	rand.Seed(time.Now().UnixNano())

	//randomly choose a wastepaper from the library
	wastepaperPos = rand.Intn(len(r.lib.wastepaper)) + 1
	targetWastepaper = r.lib.GetWastepaper(wastepaperPos)
	//if wastepaper is nil do nothing
	if targetWastepaper == nil {
		return nil
	}

	//if wastepaper is read at the previous iteration put it back
	if targetWastepaper == r.lastWastepaper {
		r.lib.PutWastepaper(targetWastepaper)
		return nil
	}

	//else this wastepaper is last
	r.lastWastepaper = targetWastepaper

	//check whether this wastepaper have already read
	if _, ok := r.wastepaperRead[targetWastepaper]; !ok {
		r.wastepaperRead[targetWastepaper] = &readingInfo{
			readText:       "",
			lastReadingPos: 0,
		}
	}

	targetWastepaperInfo = r.wastepaperRead[targetWastepaper]
	//get last reading position
	lastReadingPos = targetWastepaperInfo.lastReadingPos
	//create Reader to read only one word
	rd := strings.NewReader(targetWastepaper.GetText())
	r.lib.PutWastepaper(targetWastepaper)

	//if last reading position equal string length then return
	if rd.Len() == lastReadingPos {
		return nil
	}
	//move carriage to last reading position
	_, err := rd.Seek(int64(lastReadingPos), io.SeekStart)
	if err != nil {
		return err
	}

	//create buffer of countWordBytes() size
	buf := make([]byte, countWordBytes(rd))

	//move carriage to last reading position again
	_, err = rd.Seek(int64(lastReadingPos), io.SeekStart)
	if err != nil {
		return err
	}

	//read word to buffer
	_, err = rd.Read(buf)
	if err != nil {
		return err
	}

	//increase last reading position
	lastReadingPos += len(buf)
	//add a word to the read text
	targetWastepaperInfo.readText += string(buf)
	targetWastepaperInfo.lastReadingPos = lastReadingPos
	r.wastepaperRead[targetWastepaper] = targetWastepaperInfo
	//return carriage to compare length of reading text
	rd.Seek(0, io.SeekStart)
	if rd.Len() == lastReadingPos {
		r.wastepaperFinishReading = append(r.wastepaperFinishReading, targetWastepaper)
		//if finish reading all wastepapers send myself to winner channel
		if len(r.wastepaperFinishReading) == len(r.lib.wastepaper) {
			r.once.Do(func() {
				r.winner <- r
			})
		}
	}

	return nil
}

func countWordBytes(reader *strings.Reader) int {
	r := bufio.NewReader(reader)
	token, _ := r.ReadBytes(' ')
	return len(token)
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
