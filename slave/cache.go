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
package slave

import (
	"github.com/756445638/somecache/lru"
)

//type Group struct {
//	groups map[string]*lru.Group
//}

//var (
//	groups = Group{groups: make(map[string]*lru.Group)}
//)

//func (g *Group) getGroup(name string) *lru.Group {
//	ele, ok := g.groups[name]
//	if ok {
//		return ele
//	}
//	ne := lru.NewGroup(name)
//	g.groups[name] = ne
//	return ne

//}

var cache lru.Lru
