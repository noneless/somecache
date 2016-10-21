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

func New(size int) *Window {
	if size < 0 {
		panic("size must be positive number")
	}
	return &Window{size: size, load: []interface{}{}}
}
