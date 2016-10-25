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
package master

//	"fmt"

type Getter interface {
	Get(string) ([]byte, error)
}

var (
	getter Getter
)

func RegisterGetter(g Getter) {
	getter = g
}

func Get(k string) ([]byte, error) {
	//	if !validKey(k) {
	//		return nil, fmt.Errorf("key is valid")
	//	}
	return service.Get(k)
}

/*
//type ReaderAndSeekerAndCloser interface {
//	Read(p []byte) (n int, err error)
//	// Seek sets the offset for the next Read or Write on file to offset, interpreted
//	// according to whence: 0 means relative to the origin of the file, 1 means
//	// relative to the current offset, and 2 means relative to the end.
//	// It returns the new offset and an error, if any.
//	// The behavior of Seek on a file opened with O_APPEND is not specified.
//	Seek(offset int64, whence int) (ret int64, err error)
//	Close()
//}

//	streamGetter ReaderAndSeekerAndCloser
//func RegisterSteamGetter(g ReaderAndSeekerAndCloser) {
//	streamGetter = g
//}
*/
