package slave

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/756445638/somecache/common"
	"github.com/756445638/somecache/message"
)

func Connection2Master(tcp_addr string, cachesize int64) {
	cache.SetMaxCacheSize(cachesize)
	for {
		conn, err := net.Dial("tcp", tcp_addr)
		if err != nil {
			fmt.Println("dail master server failed,err:", err)
			continue
		}
		_, err = conn.Write(common.MagicV1)
		if err != nil {
			fmt.Println("write magic version v1 failed,err:", err)
			continue
		}
		v1s := &V1Slave{
			conn: conn,
		}
		e := v1s.MainLoop()
		if e != nil {
			fmt.Println("IOLoop failed,err:", e)
		}
	}
}

type V1Slave struct {
	conn     net.Conn
	reader   *bufio.Reader
	pingpool sync.Pool
	lock     sync.Mutex
}

func (v1s *V1Slave) MainLoop() error {
	defer v1s.conn.Close()
	v1s.reader = bufio.NewReader(v1s.conn)
	if err := v1s.Login(); err != nil {
		return fmt.Errorf("login failed,err:%v", err)
	}
	fmt.Println("login ok")
	for {
		v1s.conn.SetDeadline(time.Now().Add(30 * time.Second))
		jodid, line, err := common.ReadLine(v1s.reader)
		if err != nil {
			return err
		}
		fmt.Printf("read line,data[%s] raw[%v]\n", string(line), line)
		cmd := &common.Command{}
		err = v1s.Exec(cmd, line)
		cmd.Jobid = jodid
		if err != nil {
			cmd.Command = []byte(err.Error())
		} else {
			cmd.Command = common.OK
		}
		_, err = cmd.Write(v1s.conn)
		if err != nil { // write faild,must return,fetal error
			return err
		}
	}
	return nil
}

//ret is returen value
func (v1s *V1Slave) Exec(ret *common.Command, line []byte) error { // error just for log
	var err error
	cmd, para := common.ParseCommand(line)
	if err != nil {
		return err
	}
	if bytes.Equal(cmd, common.COMMAND_PING) {
		return v1s.Ping(ret)
	} else if bytes.Equal(cmd, common.COMMAND_GET) {
		return v1s.Get(ret, para)
	} else if bytes.Equal(cmd, common.COMMAND_PUT) {
		return v1s.Put(para)
	} else {
		return errors.New(string(common.E_COMMAND_NOT_FOUND))
	}
	return nil
}

func (v1s *V1Slave) Login() error {
	l := message.Login{}
	body, err := json.Marshal(l)
	if err != nil {
		return err
	}
	_, err = common.NewCommand(common.COMMAND_LOGIN, nil, body, 0).Write(v1s.conn)
	if err != nil {
		return err
	}
	_, line, err := common.ReadLine(v1s.reader)
	if err != nil {
		return err
	}
	if !bytes.Equal(common.OK, line) {
		return fmt.Errorf("master response something[%s] but not ok,", string(line))
	}
	return nil
}

func (v1s *V1Slave) Get(ret *common.Command, para [][]byte) error {
	if len(para) != 1 {
		return fmt.Errorf("must have 1 parameter")
	}
	v := cache.Get(string(para[0]))
	if v == nil {
		return fmt.Errorf(string(common.E_NOT_FOUND))
	}
	data := v.(*common.BytesData).Data
	ret.Content = data
	return nil
}

func (v1s *V1Slave) Put(para [][]byte) error {
	if len(para) != 1 {
		return fmt.Errorf("must have 1 parameter")
	}
	key := string(para[0])
	buf, _, err := common.ReadBody4(v1s.reader, nil)
	if err != nil {
		return err
	}
	d := &common.BytesData{Data: buf, K: key}
	cache.Put(key, d)
	return nil
}
func (v1s *V1Slave) Ping(cmd *common.Command) error {
	v := v1s.pingpool.Get()
	if v == nil {
		v = &message.HeartBeat{}
	}
	defer v1s.pingpool.Put(v)
	m := v.(*message.HeartBeat)
	m.Lru_hit = cache.Hit()
	m.Lru_cachedsize = cache.CachedSize()
	m.Lru_maxcachesize = cache.MaxCacheSize()
	m.Lru_gets = cache.Gets()
	m.Lru_puts = cache.Puts()
	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal error:", err)
	}
	cmd.Content = body
	return nil
}
