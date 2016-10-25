/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
