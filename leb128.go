package dexy

func Sleb128(p []byte) (n int64) {
	shift := uint(0)
	for _, b := range p {
		n = (n | int64(b&0x7F)) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}

}

func Uleb128(p []byte) (n int64) {
	shift := uint(0)
	for _, b := range p {
		n = (n | int64(b&0x7F)) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return
}

func Uleb128p1(p []byte) (n int64) {
	return
}
