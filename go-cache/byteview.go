package go_cache

// ByteView 只读
type ByteView struct {
	b []byte
}

func (b *ByteView) Len() int {
	return len(b.b)
}

func (b *ByteView) String() string {
	return string(b.b)
}

// ByteSlice 返回拷贝，防止缓存值被外部程序修改
func (b *ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
