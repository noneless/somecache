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
