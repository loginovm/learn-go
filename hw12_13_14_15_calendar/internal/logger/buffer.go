package logger

type Buffer []byte

func (b *Buffer) WriteString(s string) {
	*b = append(*b, s...)
}

func (b *Buffer) WriteByte(c byte) {
	*b = append(*b, c)
}
