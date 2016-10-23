//window is a fixed size of array
package window

import (
	"fmt"
)

type Window struct {
	size int
	load []interface{}
}

func (w *Window) Append(d ...interface{}) {
	w.load = append(w.load, d...)
	fmt.Println(w.load)
	if t := len(w.load) - w.size; t > 0 {
		w.load = w.load[t:]
	}
}

var defaultSize int = 10

func New(size int) *Window {
	if size < 0 {
		size = defaultSize
	}
	return &Window{size: size, load: []interface{}{}}
}
