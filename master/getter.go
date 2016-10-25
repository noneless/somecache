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

type Getter interface {
	Get(string) ([]byte, error)
}

var (
	getter Getter
)

func RegisterGetter(g Getter) {
	getter = g
}

func GetFromRemoteServer(k string) ([]byte, error) {
	return service.getRemoteCache(k)
}

func Put(k string, data []byte, remote bool) {
	service.Put(k, data, remote)
}

func Get(k string) ([]byte, error) {
	//	if !validKey(k) {
	//		return nil, fmt.Errorf("key is valid")
	//	}
	return service.Get(k)
}
