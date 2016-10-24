package common

type BytesData struct {
	Data []byte
	K    string
}

func (b *BytesData) Measure() int64 {
	length := len(b.K) + len(b.Data)
	return int64(length)
}

func (b *BytesData) Key() string {
	return b.K
}
