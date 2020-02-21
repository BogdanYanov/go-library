package library

import (
	"bytes"
	"context"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestNewReader(t *testing.T) {
	var (
		ctx    = context.Background()
		ID     = 1
		lib    = NewLibrary()
		winner = make(chan *Reader)
		once   = &sync.Once{}
	)

	type args struct {
		ctx    context.Context
		ID     int
		lib    *Library
		winner chan *Reader
		once   *sync.Once
	}
	tests := []struct {
		name string
		args args
		want *Reader
	}{
		{
			name: "NewReader() 1",
			args: args{
				ctx:    ctx,
				ID:     ID,
				lib:    lib,
				winner: winner,
				once:   once,
			},
			want: &Reader{
				ID:                      ID,
				lib:                     lib,
				wastepaperRead:          make(map[Wastepaper]*readingInfo),
				wastepaperFinishReading: make([]Wastepaper, 0, len(lib.wastepaper)),
				lastWastepaper:          nil,
				winner:                  winner,
				once:                    once,
				ctx:                     ctx,
			},
		},
		{
			name: "NewReader() 2",
			args: args{
				ctx:    nil,
				ID:     ID,
				lib:    lib,
				winner: winner,
				once:   once,
			},
			want: nil,
		},
		{
			name: "NewReader() 3",
			args: args{
				ctx:    ctx,
				ID:     ID,
				lib:    nil,
				winner: winner,
				once:   once,
			},
			want: nil,
		},
		{
			name: "NewReader() 4",
			args: args{
				ctx:    ctx,
				ID:     ID,
				lib:    lib,
				winner: nil,
				once:   once,
			},
			want: nil,
		},
		{
			name: "NewReader() 5",
			args: args{
				ctx:    ctx,
				ID:     ID,
				lib:    lib,
				winner: winner,
				once:   nil,
			},
			want: nil,
		},
		{
			name: "NewReader() 6",
			args: args{
				ctx:    nil,
				ID:     ID,
				lib:    nil,
				winner: nil,
				once:   nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReader(tt.args.ctx, tt.args.ID, tt.args.lib, tt.args.winner, tt.args.once); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_Read(t *testing.T) {
	var (
		ctx    = context.Background()
		once   = &sync.Once{}
		winner = make(chan *Reader)
		lib    = NewLibrary()
		book   = NewBook("Author", 1976, "Simple Text")
	)

	lib.AddWastepaper(book)

	type fields struct {
		ID     int
		lib    *Library
		winner chan *Reader
		once   *sync.Once
		ctx    context.Context
	}

	testLastWastepaper := []struct {
		name           string
		fields         fields
		do             func()
		lastWastepaper Wastepaper
	}{
		{
			name: "Reader.lastWastepaper 1",
			fields: fields{
				ID:     0,
				lib:    lib,
				winner: winner,
				once:   once,
				ctx:    ctx,
			},
			do:             func() {},
			lastWastepaper: book,
		},
		{
			name: "Reader.lastWastepaper 2",
			fields: fields{
				ID:     0,
				lib:    lib,
				winner: winner,
				once:   once,
				ctx:    ctx,
			},
			do: func() {
				book.Take()
			},
			lastWastepaper: nil,
		},
	}

	// test assignment last wastepaper
	for _, tt := range testLastWastepaper {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()
			r := NewReader(tt.fields.ctx, tt.fields.ID, tt.fields.lib, tt.fields.winner, tt.fields.once)
			if err := r.Read(); err != nil {
				t.Errorf("Read() error: %s", err)
				return
			}
			if r.lastWastepaper != tt.lastWastepaper {
				t.Errorf("Reader.lastWastepaper = %v, want - %v", r.lastWastepaper, tt.lastWastepaper)
				return
			}
		})
	}

	// test reading the last wastepaper
	var oldRInfo, newRInfo *readingInfo
	r := NewReader(ctx, 0, lib, winner, once)
	book.Put()
	err := r.Read()
	if err != nil {
		t.Errorf("Read() error: %s", err)
		return
	}
	oldRInfo, ok := r.wastepaperRead[book]
	if !ok {
		t.Errorf("Read() don`r read first book")
		return
	}
	err = r.Read()
	if err != nil {
		t.Errorf("Read() error: %s", err)
		return
	}
	newRInfo = r.wastepaperRead[book]
	if strings.Compare(oldRInfo.readText, newRInfo.readText) != 0 {
		t.Errorf("Read() read the last wastepaper again")
		return
	}

	// test errors
	testErrors := []struct {
		name     string
		fields   fields
		doFirst  func() *Library
		doSecond func(r *Reader)
		wantErr  bool
	}{
		{
			name: "Read() test errors 1",
			fields: fields{
				ID:     0,
				lib:    lib,
				winner: winner,
				once:   once,
				ctx:    ctx,
			},
			doFirst: func() *Library {
				go func() {
					<-winner
				}()
				lib := NewLibrary()
				lib.AddWastepaper(NewBook("Author", 1976, "abcd"))
				return lib
			},
			doSecond: func(r *Reader) {
				r.lastWastepaper = nil
			},
			wantErr: false,
		},
		{
			name: "Read() test errors 2",
			fields: fields{
				ID:     0,
				lib:    lib,
				winner: winner,
				once:   once,
				ctx:    ctx,
			},
			doFirst: func() *Library {
				lib := NewLibrary()
				lib.AddWastepaper(NewBook("Author", 1976, "abcd abcd"))
				return lib
			},
			doSecond: func(r *Reader) {
				r.lastWastepaper = nil
				rInfo := r.wastepaperRead[r.lib.wastepaper[0]]
				rInfo.lastReadingPos = -15
				r.wastepaperRead[lib.wastepaper[0]] = rInfo
			},
			wantErr: true,
		},
	}

	for _, te := range testErrors {
		t.Run(te.name, func(t *testing.T) {
			r := NewReader(te.fields.ctx, te.fields.ID, te.fields.lib, te.fields.winner, te.fields.once)
			r.lib = te.doFirst()
			_ = r.Read()
			te.doSecond(r)
			if err := r.Read(); (err != nil) != te.wantErr {
				t.Errorf("Read() error: %s, want - %v", err, te.wantErr)
				return
			}
		})
	}
}

func TestReader_WastepaperReadInfo(t *testing.T) {
	var exampleString = `Book by A1 (1976 publishing year)
Read text: abcd

`
	var (
		winner = make(chan *Reader)
		lib    = NewLibrary()
	)

	lib.AddWastepaper(NewBook("A1", 1976, "abcd"))
	reader := NewReader(context.Background(), 0, lib, winner, &sync.Once{})

	go func() {
		<-winner
	}()

	err := reader.Read()
	if err != nil {
		t.Errorf("Read() error: %s", err)
	}

	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	reader.WastepaperReadInfo()

	err = w.Close()
	if err != nil {
		t.Errorf("Error closing pipe: %s\n", err)
	}

	os.Stdout = oldOutput

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Errorf("Error while copy: %s\n", err)
	}

	if equal := strings.Compare(exampleString, buf.String()); equal != 0 {
		t.Errorf("WastepaperReadInfo() = %s, want - %s", buf.String(), exampleString)
	}
}

func Test_countWordBytes(t *testing.T) {
	type args struct {
		reader *strings.Reader
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "countWordBytes() 1",
			args: args{strings.NewReader("abcd")},
			want: 4,
		},
		{
			name: "countWordBytes() 2",
			args: args{strings.NewReader("")},
			want: 0,
		},
		{
			name: "countWordBytes() 3",
			args: args{strings.NewReader("Some text")},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countWordBytes(tt.args.reader); got != tt.want {
				t.Errorf("countWordBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
