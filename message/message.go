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
