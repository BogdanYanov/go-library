package library

import (
	"reflect"
	"sync"
	"testing"
)

func TestBook_GetText(t *testing.T) {
	type fields struct {
		author         string
		publishingYear int
		text           string
		mu             *sync.Mutex
		isTaken        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Book.GetText() 1",
			fields: fields{text: "One two three"},
			want:   "One two three",
		},
		{
			name:   "Book.GetText() 2",
			fields: fields{text: ""},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Book{
				author:         tt.fields.author,
				publishingYear: tt.fields.publishingYear,
				text:           tt.fields.text,
				mu:             tt.fields.mu,
				isTaken:        tt.fields.isTaken,
			}
			if got := b.GetText(); got != tt.want {
				t.Errorf("GetText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBook_IsTaken(t *testing.T) {
	type fields struct {
		author         string
		publishingYear int
		text           string
		mu             *sync.Mutex
		isTaken        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Book.IsTaken() 1",
			fields: fields{},
			want:   false,
		},
		{
			name:   "Book.IsTaken() 2",
			fields: fields{isTaken: true},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Book{
				author:         tt.fields.author,
				publishingYear: tt.fields.publishingYear,
				text:           tt.fields.text,
				mu:             tt.fields.mu,
				isTaken:        tt.fields.isTaken,
			}
			if got := b.IsTaken(); got != tt.want {
				t.Errorf("IsTaken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBook_Put(t *testing.T) {
	type fields struct {
		author         string
		publishingYear int
		text           string
		mu             *sync.Mutex
		isTaken        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Book.Put() 1",
			fields: fields{},
			want:   false,
		},
		{
			name:   "Book.Put() 2",
			fields: fields{isTaken: true},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Book{
				author:         tt.fields.author,
				publishingYear: tt.fields.publishingYear,
				text:           tt.fields.text,
				mu:             tt.fields.mu,
				isTaken:        tt.fields.isTaken,
			}
			b.Put()
			if got := b.IsTaken(); got != tt.want {
				t.Errorf("IsTaken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBook_Take(t *testing.T) {
	type fields struct {
		author         string
		publishingYear int
		text           string
		mu             *sync.Mutex
		isTaken        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Book.Take() 1",
			fields: fields{},
			want:   true,
		},
		{
			name:   "Book.Take() 2",
			fields: fields{isTaken: true},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Book{
				author:         tt.fields.author,
				publishingYear: tt.fields.publishingYear,
				text:           tt.fields.text,
				mu:             tt.fields.mu,
				isTaken:        tt.fields.isTaken,
			}
			b.Take()
			if got := b.IsTaken(); got != tt.want {
				t.Errorf("IsTaken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBook(t *testing.T) {
	type args struct {
		author         string
		publishingYear int
		text           string
	}
	tests := []struct {
		name string
		args args
		want Wastepaper
	}{
		{
			name: "NewBook() 1",
			args: args{
				author:         "",
				publishingYear: 1976,
				text:           "Some text",
			},
			want: nil,
		},
		{
			name: "NewBook() 2",
			args: args{
				author:         "    ",
				publishingYear: 1976,
				text:           "Some text",
			},
			want: nil,
		},
		{
			name: "NewBook() 4",
			args: args{
				author:         "Author",
				publishingYear: 0,
				text:           "Some text",
			},
			want: nil,
		},
		{
			name: "NewBook() 5",
			args: args{
				author:         "Author",
				publishingYear: 867,
				text:           "Some text",
			},
			want: nil,
		},
		{
			name: "NewBook() 6",
			args: args{
				author:         "Author",
				publishingYear: 2021,
				text:           "Some text",
			},
			want: nil,
		},
		{
			name: "NewBook() 8",
			args: args{
				author:         "Author",
				publishingYear: 1976,
				text:           "",
			},
			want: nil,
		},
		{
			name: "NewBook() 9",
			args: args{
				author:         "Author",
				publishingYear: 1976,
				text:           "      ",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBook(tt.args.author, tt.args.publishingYear, tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBook() = %v, want %v", got, tt.want)
				return
			}
		})
	}
	testBook := NewBook("   author", 869, "   simple text")
	if testBook == nil {
		t.Errorf("NewBook() is nil!")
		return
	}
}
