package language

import (
	"reflect"
	"testing"
)

func TestNewOperands(t *testing.T) {
	tests := []struct {
		input interface{}
		ops   *Operands
		err   bool
	}{
		{int64(0), &Operands{0.0, 0, 0, 0, 0, 0}, false},
		{int64(1), &Operands{1.0, 1, 0, 0, 0, 0}, false},
		{"0", &Operands{0.0, 0, 0, 0, 0, 0}, false},
		{"1", &Operands{1.0, 1, 0, 0, 0, 0}, false},
		{"1.0", &Operands{1.0, 1, 1, 0, 0, 0}, false},
		{"1.00", &Operands{1.0, 1, 2, 0, 0, 0}, false},
		{"1.3", &Operands{1.3, 1, 1, 1, 3, 3}, false},
		{"1.30", &Operands{1.3, 1, 2, 1, 30, 3}, false},
		{"1.03", &Operands{1.03, 1, 2, 2, 3, 3}, false},
		{"1.230", &Operands{1.23, 1, 3, 2, 230, 23}, false},
		{"20.0230", &Operands{20.023, 20, 4, 3, 230, 23}, false},
		{20.0230, nil, true},
	}
	for _, test := range tests {
		ops, err := newOperands(test.input)
		if err != nil && !test.err {
			t.Errorf("newOperands(%#v) unexpected error: %s", test.input, err)
		} else if err == nil && test.err {
			t.Errorf("newOperands(%#v) returned %#v; expected error", test.input, ops)
		} else if !reflect.DeepEqual(ops, test.ops) {
			t.Errorf("newOperands(%#v) returned %#v; expected %#v", test.input, ops, test.ops)
		}
	}
}

func BenchmarkNewOperand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := newOperands("1234.56780000"); err != nil {
			b.Fatal(err)
		}
	}
}
