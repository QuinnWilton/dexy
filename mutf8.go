package dexy

import (
	"bytes"
	"errors"
)

const (
	RuneError = '\uFFFD' // the "error" Rune or "Unicode replacement character"
	RuneSelf  = '\u0080' // characters below Runeself are represented as themselves in a single byte.

	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
)

func DecodeMutf8Rune(p []byte) (r rune, n int) {
	a := p[n]
	n++
	switch {
	case a < RuneSelf:
		r = rune(a)
		return
	case (a & t3) == t2:
		b := p[n]
		n++
		if (b & t2) != tx {
			return RuneError, 2
		}
		r = rune(a&0x1F)<<6 | rune(b&0x3F)
		return r, 2
	case (a & t4) == t3:
		b := p[n]
		n++
		if (b & t2) != tx {
			return RuneError, 2
		}
		c := p[n]
		n++
		if (c & t2) != tx {
			return RuneError, 3
		}
		r = rune(a&0x0F)<<12 | rune(b&0x3F)<<6 | rune(c&0x3F)
		return r, 3
	}
	return RuneError, 3
}

func Mutf8(p []byte) (s string, err error) {
	var buffer bytes.Buffer
	for i := range p {
		if p[i] == 0x0 {
			s = buffer.String()
			return
		}
		r, _ := DecodeMutf8Rune(p[i:])
		if r == RuneError {
			err = errors.New("Mutf8: String contains invalid code point")
			return
		}
		buffer.WriteRune(r)
	}
	err = errors.New("Mutf8: Byte array must not be empty")
	return
}
