package lru

import (
	//	"container/list"
	"fmt"
	"testing"

	"github.com/756445638/somecache/common"
)

func Test_lru(t *testing.T) {
	l := &Lru{maxcachesize: 10}
	var value common.BytesDate
	value = []byte{97, 97} // a a
	for i := 0; i < 10; i++ {
		l.Put("k"+fmt.Sprintf("%d", i), value)
	}

	fmt.Println(l.Get("k9"))

}
