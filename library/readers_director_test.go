package library

import (
	"context"
	"reflect"
	"testing"
)

func TestNewReadersDirector(t *testing.T) {
	type args struct {
		ctx        context.Context
		lib        *Library
		readersNum int
	}
	tests := []struct {
		name string
		args args
		want *ReadersDirector
	}{
		{
			name: "NewReadersDirector() 1",
			args: args{
				ctx:        nil,
				lib:        NewLibrary(),
				readersNum: 2,
			},
			want: nil,
		},
		{
			name: "NewReadersDirector() 2",
			args: args{
				ctx:        context.Background(),
				lib:        nil,
				readersNum: 2,
			},
			want: nil,
		},
		{
			name: "NewReadersDirector() 3",
			args: args{
				ctx:        context.Background(),
				lib:        NewLibrary(),
				readersNum: -1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReadersDirector(tt.args.ctx, tt.args.lib, tt.args.readersNum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReadersDirector() = %v, want %v", got, tt.want)
			}
		})
	}

	var (
		ctx        = context.Background()
		lib        = NewLibrary()
		readersNum = 2
	)

	rd1 := NewReadersDirector(ctx, lib, readersNum)

	rd2 := NewReadersDirector(ctx, lib, readersNum)

	if reflect.DeepEqual(rd1, rd2) {
		t.Errorf("NewReadersDirector() = %v should not equal %v", rd1, rd2)
	}
}

func TestReadersDirector_StartWork(t *testing.T) {
	var (
		ctx, _     = context.WithCancel(context.Background())
		lib        = NewLibrary()
		readersNum = 2
	)

	lib.AddWastepaper(NewBook("A1", 1876, "abcd"))
	lib.AddWastepaper(NewBook("A2", 1876, "efgh"))
	lib.AddWastepaper(NewBook("A3", 1876, "ijkl"))

	rd := NewReadersDirector(ctx, lib, readersNum)

	winnerCh := rd.StartWork()

	if winnerCh == nil {
		t.Errorf("StartWork() channel != nil")
	}

}
