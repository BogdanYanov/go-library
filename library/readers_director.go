package library

import (
	"context"
	"log"
	"sync"
)

// ReadersDirector keeps and controls all readers
type ReadersDirector struct {
	readers []*Reader
	winner  chan *Reader
	ctx     context.Context
	once    *sync.Once
}

// NewReadersDirector creates new ReadersDirector
func NewReadersDirector(ctx context.Context, lib *Library, readersNum int) *ReadersDirector {
	if ctx == nil {
		return nil
	}
	if lib == nil {
		return nil
	}
	if readersNum < 1 {
		return nil
	}
	rd := &ReadersDirector{
		readers: make([]*Reader, 0, readersNum),
		winner:  make(chan *Reader),
		ctx:     ctx,
		once:    &sync.Once{},
	}
	for i := 0; i < readersNum; i++ {
		rd.readers = append(rd.readers, NewReader(rd.ctx, i+1, lib, rd.winner, rd.once))
	}
	return rd
}

// StartWork begins reading for all readers and returns the channel with the winning reader
func (rd *ReadersDirector) StartWork() chan *Reader {
	for i := 0; i < len(rd.readers); i++ {
		go func(reader *Reader) {
			for {
				select {
				case <-reader.ctx.Done():
					rd.once.Do(func() {
						close(rd.winner)
					})
					return
				default:
					err := reader.Read()
					if err != nil {
						log.Println(err)
						log.Printf("Stop the reader#%d\n", reader.ID)
						return
					}
				}
			}
		}(rd.readers[i])
	}
	return rd.winner
}
