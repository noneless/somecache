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

var cache lru.CaChe
