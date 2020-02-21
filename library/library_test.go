package library

import (
	"reflect"
	"testing"
)

func TestLibrary_AddWastepaper(t *testing.T) {
	type args struct {
		wp Wastepaper
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Library.AddWastepaper() 1",
			args: args{wp: nil},
			want: 0,
		},
		{
			name: "Library.AddWastepaper() 2",
			args: args{wp: NewBook("Author", 869, "Simple text")},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lib := NewLibrary()
			lib.AddWastepaper(tt.args.wp)
			if len(lib.wastepaper) != tt.want {
				t.Errorf("len(Library.wastepaper) = %d, want - %d", len(lib.wastepaper), tt.want)
				return
			}
		})
	}
	testSliceMake := []struct {
		name string
		args func() *Library
		want []Wastepaper
	}{
		{
			name: "Test making slice 1",
			args: func() *Library {
				return NewLibrary()
			},
			want: nil,
		},
		{
			name: "Test making slice 2",
			args: func() *Library {
				lib := NewLibrary()
				lib.AddWastepaper(&Book{
					author:         "Author",
					publishingYear: 1968,
					text:           "Simple Text",
				})
				return lib
			},
			want: make([]Wastepaper, 1, defaultWastepaperSliceSize),
		},
	}
	for _, tt := range testSliceMake {
		t.Run(tt.name, func(t *testing.T) {
			lib := tt.args()
			if len(lib.wastepaper) != len(tt.want) || cap(lib.wastepaper) != cap(tt.want) {
				t.Errorf("len(Library.wastepaper) = %d, want - %d\n", len(lib.wastepaper), len(tt.want))
				t.Errorf("cap(Library.wastepaper) = %d, want - %d\n", cap(lib.wastepaper), cap(tt.want))
				return
			}
		})
	}
}

func TestLibrary_GetWastepaper(t *testing.T) {
	var testBook1 = NewBook("Author", 1976, "Simple text")
	var testBook2 = NewBook("Author", 1976, "Simple text")
	testBook2.Take()

	type fields struct {
		wastepaper []Wastepaper
	}
	type args struct {
		pos int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Wastepaper
	}{
		{
			name:   "Library.GetWastepaper() 1",
			fields: fields{wastepaper: nil},
			args:   args{pos: 1},
			want:   nil,
		},
		{
			name:   "Library.GetWastepaper() 2",
			fields: fields{wastepaper: nil},
			args:   args{pos: -1},
			want:   nil,
		},
		{
			name:   "Library.GetWastepaper() 3",
			fields: fields{wastepaper: []Wastepaper{testBook1}},
			args:   args{pos: 2},
			want:   nil,
		},
		{
			name:   "Library.GetWastepaper() 4",
			fields: fields{wastepaper: []Wastepaper{testBook2}},
			args:   args{pos: 1},
			want:   nil,
		},
		{
			name:   "Library.GetWastepaper() 5",
			fields: fields{wastepaper: []Wastepaper{testBook1}},
			args:   args{pos: 1},
			want:   testBook1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lib := &Library{
				wastepaper: tt.fields.wastepaper,
			}
			if got := lib.GetWastepaper(tt.args.pos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWastepaper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLibrary_PutWastepaper(t *testing.T) {
	var testBook1 = NewBook("Author1", 1976, "Simple Text")
	testBook1.Take()
	var testBook2 = NewBook("Author2", 1976, "Simple Text")
	testBook2.Take()
	var testBook3 = NewBook("Author3", 1976, "Simple Text")

	var wp = []Wastepaper{testBook1, testBook2, testBook3}

	var lib = &Library{wastepaper: wp}

	lib.PutWastepaper(testBook1)
	if testBook1.IsTaken() {
		t.Errorf("Book taken = %v, want - %v", testBook1.IsTaken(), false)
	}

	lib.PutWastepaper(testBook2)
	if testBook2.IsTaken() {
		t.Errorf("Book taken = %v, want - %v", testBook2.IsTaken(), false)
	}

	lib.PutWastepaper(testBook3)
	if testBook3.IsTaken() {
		t.Errorf("Book taken = %v, want - %v", testBook3.IsTaken(), false)
	}
}

func TestNewLibrary(t *testing.T) {
	lib := NewLibrary()
	if lib == nil {
		t.Errorf("Library is nil!")
	}
}
