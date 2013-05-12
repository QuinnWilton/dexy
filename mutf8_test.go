package dexy

import (
	"testing"
)

var decodeMutf8RuneTests = []struct {
	in   []byte
	r    rune
	size int
}{
	{[]byte{0x45}, '\u0045', 1},
	{[]byte{0xC8, 0x85}, '\u0205', 2},
	//{[]byte{0xED, 0xA0, 0x81, 0xED, 0xB0, 0x80}, '\U00010400', 6}, // TODO learn how to handle surrogate pairs
}

var mutf8Tests = []struct {
	in []byte
	s  string
}{
	{[]byte{0x28, 0x5A, 0x54, 0x4B, 0x3B, 0x54, 0x56, 0x3B, 0x54, 0x56, 0x3B, 0x29, 0x56, 0x0},
		"(ZTK;TV;TV;)V"},
	{[]byte{0x0},
		""},
}

func Test_DecodeMutf8Rune(t *testing.T) {
	for i, test := range decodeMutf8RuneTests {
		s, n := DecodeMutf8Rune(test.in)
		if s != test.r || n != test.size {
			t.Errorf("%d. Mutf8(%q) => %q and %q, want %q and %q", i, test.in, s, n, test.r, test.size)
		}
	}
}

func Test_Mutf8(t *testing.T) {
	for i, test := range mutf8Tests {
		s, _ := Mutf8(test.in)
		if s != test.s {
			t.Errorf("%d. Mutf8(%q) => %q, want %#v", i, test.in, s, test.s)
		}
	}
}
