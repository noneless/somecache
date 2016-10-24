package common

import (
	"io"
)

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

func WriteFull(write io.Writer, buf []byte) error {
	length := len(buf)
	for length > 0 {
		n, err := write.Write(buf)
		if err != nil {
			return err
		}
		length -= n
		buf = buf[n:]
	}
	return nil
}
