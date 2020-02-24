package library

import (
	"context"
	"reflect"
	"sync"
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
	var lib = NewLibrary()

	lib.AddWastepaper(NewBook("A1", 1876, "abcd"))
	lib.AddWastepaper(NewBook("A2", 1876, "efgh"))
	lib.AddWastepaper(NewBook("A3", 1876, "ijkl"))

	type args struct {
		lib        *Library
		readersNum int
	}

	tests := []struct {
		name    string
		args    args
		canSend bool
		do      func(wg *sync.WaitGroup, cancel context.CancelFunc)
	}{
		{
			name: "StartWork() case 1",
			args: args{
				lib:        lib,
				readersNum: 2,
			},
			canSend: false,
			do: func(wg *sync.WaitGroup, cancel context.CancelFunc) {
				cancel()
				wg.Wait()
			},
		},
		{
			name: "StartWork() case 2",
			args: args{
				lib:        lib,
				readersNum: 2,
			},
			canSend: true,
			do: func(wg *sync.WaitGroup, cancel context.CancelFunc) {
				wg.Wait()
				cancel()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var canSend bool

			ctx, cancel := context.WithCancel(context.Background())
			wg := &sync.WaitGroup{}

			rd := NewReadersDirector(ctx, tt.args.lib, tt.args.readersNum)

			wg.Add(1)
			go func() {
				_, canSend = <-rd.winner
				wg.Done()
			}()

			rd.StartWork()

			tt.do(wg, cancel)

			if tt.canSend != canSend {
				t.Errorf("StartWork(): canSend = %v, want - %v", canSend, tt.canSend)
			}
		})
	}

}
