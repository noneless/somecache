package lru

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
)

type Lru struct {
	maincache    CaChe
	cachedsize   int64 `json:"cachedsize"`
	maxcachesize int64
	hit          int64
	gets         int64
	puts         int64
	lock         sync.Mutex
}

func (l *Lru) SetMaxCacheSize(size int64) {
	l.maxcachesize = size
}
func (l *Lru) GetMaxCacheSize() int64 {
	return l.maxcachesize
}

//func (l *Lru) GroupName() string {
//	return l.groupname
//}

type Measureable interface {
	Measure() int64
	Key() string
}

func (l *Lru) CachedSize() int64 {
	return l.cachedsize
}
func (l *Lru) MaxCacheSize() int64 {
	return l.maxcachesize
}

func (l *Lru) Hit() int64 {
	return atomic.LoadInt64(&l.hit)
}
func (l *Lru) Gets() int64 {
	return l.gets
}
func (l *Lru) Puts() int64 {
	return l.puts
}

func (l *Lru) Get(k string) interface{} {
	l.gets++
	e := l.maincache.Get(k)
	if e == nil {
		return nil
	}
	atomic.AddInt64(&l.hit, 1)
	return e.Value
}

// will overwrite
func (l *Lru) Put(k string, v Measureable) (*list.Element, error) {
	l.puts++
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
	lru  *Lru
	eles map[string]*list.Element
}

//debug method
func (c *CaChe) travel() {
	first := c.lis.Front()
	fmt.Println(first.Value)
	next := first.Next()
	for first != next && next != nil {
		fmt.Println(next.Value)
		next = next.Next()
	}
}

func (c *CaChe) mkroom(size int64) int64 {
	if c.eles == nil {
		c.eles = make(map[string]*list.Element)
	}
	var length int64
	released := int64(0)
	c.lock.Lock()
	for size > 0 {
		t := c.lis.Back() //clear
		m := t.Value.(Measureable)
		length = m.Measure()
		c.lis.Remove(t)
		delete(c.eles, m.Key())
		size -= length
		released += length
	}
	c.lock.Unlock()
	return released
}

func (c *CaChe) Get(k string) *list.Element {
	if c.eles == nil {
		return nil
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
		m := e.Value.(Measureable)
		c.lis.Remove(e)
		delete(c.eles, m.Key())
		c.lru.cachedsize -= m.Measure()
	}
	e := c.lis.PushFront(v)
	c.eles[k] = e
	c.lock.Unlock()
	return e
}
