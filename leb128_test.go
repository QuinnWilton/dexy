package dexy

import (
	"testing"
)

var sleb128Tests = []struct {
	in  []byte
	out int64
}{
	{[]byte{0x00}, 0},
	{[]byte{0x01}, 1},
	{[]byte{0x7F}, -1},
	{[]byte{0x80, 0x7F}, -128},
}

var uleb128Tests = []struct {
	in  []byte
	out int64
}{
	{[]byte{0x00}, 0},
	{[]byte{0x01}, 1},
	{[]byte{0x7F}, 127},
	{[]byte{0x80, 0x7F}, 16256},
}

var uleb128p1Tests = []struct {
	in  []byte
	out int64
}{
	{[]byte{0x00}, -1},
	{[]byte{0x01}, 0},
	{[]byte{0x7F}, 126},
	{[]byte{0x80, 0x7F}, 16255},
}

func Test_Sleb128(t *testing.T) {
	for i, test := range sleb128Tests {
		n := Sleb128(test.in)
		if n != test.out {
			t.Errorf("%d. Sleb128(%v) => %v, want %v", i, test.in, n, test.out)
		}
	}
}

func Test_Uleb128(t *testing.T) {
	for i, test := range uleb128Tests {
		n := Uleb128(test.in)
		if n != test.out {
			t.Errorf("%d. Uleb128(%v) => %v, want %v", i, test.in, n, test.out)
		}
	}
}

func Test_Uleb128p1(t *testing.T) {
	for i, test := range uleb128p1Tests {
		n := Uleb128p1(test.in)
		if n != test.out {
			t.Errorf("%d. Uleb128p1(%v) => %v, want %v", i, test.in, n, test.out)
		}
	}
}
