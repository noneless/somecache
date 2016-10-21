package window

import (
	"fmt"
	"testing"
)

func TestWindow(t *testing.T) {
	w := New(3)
	w.Append(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	for k, v := range w.load {
		fmt.Printf("%d:%d\n", k, v)
	}
	w.Append(-1, 33, 12122)
	for k, v := range w.load {
		fmt.Printf("%d:%d\n", k, v)
	}
}
