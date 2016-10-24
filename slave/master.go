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
		time.Sleep(time.Second)
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
}

func (v1s *V1Slave) MainLoop() error {
	defer v1s.conn.Close()
	v1s.reader = bufio.NewReader(v1s.conn)
	if err := v1s.Login(); err != nil {
		return fmt.Errorf("login failed,err:%v", err)
	}
	for {
		line, err := common.ReadLine(v1s.reader)
		if err != nil {
			return err
		}
		fmt.Printf("read line,data[%s]\n", string(line))
		err = v1s.Exec(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v1s *V1Slave) Login() error {
	l := message.Login{}
	body, err := json.Marshal(l)
	if err != nil {
		return err
	}
	_, err = common.NewCommand(common.COMMAND_LOGIN, nil, body).Write(v1s.conn)
	if err != nil {
		return err
	}
	line, err := common.ReadLine(v1s.reader)
	if err != nil {
		return err
	}
	if !bytes.Equal(common.OK, line) {
		return fmt.Errorf("master response something[%s] but not ok,", string(line))
	}
	return nil
}

func (v1s *V1Slave) Exec(line []byte) error { // error just for log
	cmd, para := common.ParseCommand(line)
	fmt.Printf("debug cmd[%s] para[%v]\n", string(cmd), para)
	if bytes.Equal(cmd, common.COMMAND_PING) {
		return v1s.Ping()
	} else if bytes.Equal(cmd, common.COMMAND_GET) {
		v1s.Get(para)
		return nil
	} else if bytes.Equal(cmd, common.COMMAND_PUT) {
		v1s.Put(para)
		return nil
	} else {
		v1s.WtiteError(common.E_NOT_FOUND)
		return errors.New(string(common.E_NOT_FOUND))
	}
	return nil
}

func (v1s *V1Slave) Get(para [][]byte) error {
	if len(para) != 1 {
		v1s.WtiteError(common.E_PARAMETER_ERROR)
		return fmt.Errorf("must have 1 parameter")
	}
	v := cache.Get(string(para[0]))
	if v == nil {
		v1s.WtiteError(common.E_NOT_FOUND)
		return nil
	}
	data := v.(*common.BytesData)
	_, err := common.NewCommand(common.OK, nil, data.Data).Write(v1s.conn)
	return err
}

func (v1s *V1Slave) Put(para [][]byte) error {
	if len(para) != 1 {
		v1s.WtiteError(common.E_PARAMETER_ERROR)
		return fmt.Errorf("must have 1 parameter")
	}
	key := string(para[0])
	buf, _, err := common.ReadBody4(v1s.reader, nil)
	if err != nil {
		v1s.WtiteError(common.E_READ_ERROR)
	}
	d := &common.BytesData{Data: buf, K: key}

	cache.Put(key, d)
	return nil
}
func (v1s *V1Slave) Ping() error {
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
	_, err = common.NewCommand(common.OK, nil, body).Write(v1s.conn)
	return err
}

func (v1s *V1Slave) WtiteError(reason []byte) error {
	_, err := common.NewCommand(reason, nil, nil).Write(v1s.conn)
	return err
}
