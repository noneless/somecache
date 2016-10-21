package lru

////group is a logic volume for lru
//type Group struct {
//	Name string
//	lru  *Lru
//	Hit  uint64
//}

//const (
//	default_cache_size = 1 << 30
//)

//func NewGroup(name string) *Group {
//	return &Group{
//		Name: name,
//		lru:  &Lru{cachedsize: default_cache_size},
//	}
//}

//func (g *Group) Get(k string) interface{} {
//	return g.lru.Get(k)
//}

//func (g *Group) Put(k string, m Measureable) {
//	g.lru.Put(k, m)
//}
