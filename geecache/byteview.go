package geecache

//只读的数据结构
type ByteView struct {
	//储存真实缓存，支持任意数据类型
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

//只读，防止被修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}