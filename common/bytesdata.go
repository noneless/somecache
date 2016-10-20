package common

type BytesDate []byte

func (bd BytesDate) Measure() uint64 {
	return uint64(len(bd))
}
