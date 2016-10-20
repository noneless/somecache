package common

type BytesData []byte

func (bd BytesData) Measure() uint64 {
	return uint64(len(bd))
}
