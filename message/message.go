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
package message

//	"github.com/756445638/somecache/lru"

type Login struct {
	TcpPort      int
	BroadAddress string
	Sign         string
}
type HeartBeat struct {
	Lru_hit          int64
	Lru_cachedsize   int64
	Lru_maxcachesize int64
	Lru_puts         int64
	Lru_gets         int64
	Meminfo          string // in linux it is "cat /proc/meminfo"
}
