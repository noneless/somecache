package lru

import (
	"container/list"
	"fmt"
	"testing"
	"time"
)

type el struct {
	key  string
	data interface{}
}

func (e *el) Key() string {
	return e.key
}
func (e *el) Measure() int64 {
	return int64(len(e.key)) + int64(4)
}

func Test_lru(t *testing.T) {
	l := &Lru{maxcachesize: 1 << 30}
	l.maincache.eles = make(map[string]*list.Element)
	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("k%d", i)
		l.Put(k, &el{key: k, data: i})
	}

	l.Get("k4")
	l.Get("k0")
	l.maincache.travel()
	fmt.Println("\n\n\n\n")

}

func Test_insert(t *testing.T) {
	l := &Lru{maxcachesize: 500 * 1024 * 1024}
	l.maincache.lru = l
	l.maincache.eles = make(map[string]*list.Element)
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println("l.cachedsize:", l.CachedSize())
			fmt.Println("l.maincache.ele.Len():", len(l.maincache.eles))
			fmt.Println("l.maincache.lis.Len():", l.maincache.lis.Len())
		}
	}()
	now := time.Now()
	for i := 0; i < 10*1024*1024; i++ {
		k := fmt.Sprintf("k%d", i)
		l.Put(k, &el{key: k, data: i})

	}
	fmt.Println(time.Now().Sub(now).Seconds())
}

//func BenchmarkPut(b *testing.B) {
//	l := &Lru{maxcachesize: 20 << 20}
//	l.maincache.eles = make(map[string]*list.Element)
//	for i := 0; i < b.N; i++ {
//		l.Put(fmt.Sprintf("k%d", i), m(i))
//	}
//}
