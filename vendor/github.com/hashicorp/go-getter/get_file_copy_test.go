package getter

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"
)

// OneDoneContext is a context that is
// cancelled after a first done is called.
type OneDoneContext bool

func (*OneDoneContext) Deadline() (deadline time.Time, ok bool) { return }
func (*OneDoneContext) Value(key interface{}) interface{}       { return nil }

func (o *OneDoneContext) Err() error {
	if *o == false {
		return nil
	}
	return context.Canceled
}

func (o *OneDoneContext) Done() <-chan struct{} {
	if *o == false {
		*o = true
		return nil
	}
	c := make(chan struct{})
	close(c)
	return c
}

func (o *OneDoneContext) String() string {
	if *o {
		return "done OneDoneContext"
	}
	return "OneDoneContext"
}

func TestCopy(t *testing.T) {
	const text3lines = `line1
	line2
	line3
	`

	cancelledContext, cancel := context.WithCancel(context.Background())
	_ = cancelledContext
	cancel()
	type args struct {
		ctx context.Context
		src io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantDst string
		wantErr error
	}{
		{"read all", args{context.Background(), bytes.NewBufferString(text3lines)}, int64(len(text3lines)), text3lines, nil},
		{"read none", args{cancelledContext, bytes.NewBufferString(text3lines)}, 0, "", context.Canceled},
		{"cancel after read", args{new(OneDoneContext), bytes.NewBufferString(text3lines)}, int64(len(text3lines)), text3lines, context.Canceled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &bytes.Buffer{}
			got, err := Copy(tt.args.ctx, dst, tt.args.src)
			if err != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
			if gotDst := dst.String(); gotDst != tt.wantDst {
				t.Errorf("Copy() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}
