package lru

import (
	"container/list"
	"sync"
)

type Lru struct {
	groupname    string
	maincache    CaChe
	cachedsize   uint64
	maxcachesize uint64
	lock         sync.Mutex
}

func (l *Lru) GroupName() string {
	return l.groupname
}

type Measureable interface {
	Measure() uint64
}

func (l *Lru) CachedSize() uint64 {
	return l.cachedsize
}

func (l *Lru) Get(k string) interface{} {
	e := l.maincache.Get(k)
	if e == nil {
		return nil
	}
	return e.Value
}

// will overwrite
func (l *Lru) Put(k string, v Measureable) (*list.Element, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	length := v.Measure()
	if (length + l.cachedsize) > l.maxcachesize { //force to make room for new cache
		l.cachedsize -= l.maincache.mkroom(length + l.cachedsize - l.maxcachesize)
	}
	l.cachedsize += length
	e := l.maincache.Put(k, v)
	return e, nil
}

type CaChe struct {
	lock sync.RWMutex
	lis  list.List
	eles map[string]*list.Element
}

func (c *CaChe) mkroom(size uint64) uint64 {
	var length uint64
	released := uint64(0)
	c.lock.Lock()
	for size > 0 {
		t := c.lis.Back() //clear
		length = uint64(t.Value.(Measureable).Measure())
		c.lis.Remove(t)
		size -= length
		released += length
	}
	c.lock.Unlock()
	return released
}

func (c *CaChe) Get(k string) *list.Element {
	if c.eles == nil {
		c.eles = make(map[string]*list.Element)
	}
	c.lock.RLock()
	e, ok := c.eles[k]
	c.lock.RUnlock()
	if ok {
		c.MoveToFront(e)
		return e
	}
	c.lock.Lock()
	e, ok = c.eles[k]
	c.lock.Unlock()
	if ok {
		c.MoveToFront(e)
		return e
	}
	return nil
}

func (c *CaChe) MoveToFront(e *list.Element) {
	c.lock.Lock()
	c.lis.MoveToFront(e)
	c.lock.Unlock()
}

func (c *CaChe) Put(k string, v interface{}) *list.Element {
	if c.eles == nil {
		c.eles = make(map[string]*list.Element)
	}
	c.lock.Lock()
	if e, ok := c.eles[k]; ok {
		c.lis.Remove(e)
	}
	e := c.lis.PushBack(v)
	c.eles[k] = e
	c.lock.Unlock()
	return e
}
