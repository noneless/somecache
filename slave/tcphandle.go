package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	//	"sync"
	//	"time"

	"github.com/756445638/somecache/common"
)

var (
	get common.Getter = &filegetter{}
)

type TcpHandle interface {
	IOLoop(net.Conn)
}

type V1Handle struct {
	conn net.Conn
	//sync.Mutex
}

func HandleTcp(ln net.Listener) error {
	defer wg.Done()
	b := make([]byte, 4)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		n, err := io.ReadFull(conn, b)
		if err != nil && n != 4 {
			continue
		}
		if bytes.Equal(b, common.MagicV1) { //the client should first send protocol version first
			v1h := &V1Handle{}
			go v1h.IOLoop(conn) //I am server
		} else {
			fmt.Println("un support version", string(b))
		}
	}
}

//
func (v1h *V1Handle) Ping() {
	v1h.Write(common.OK, nil, nil)
}

func (v1h *V1Handle) IOLoop(conn net.Conn) {
	reader := bufio.NewReader(v1h.conn)
	v1h.conn = conn
	defer conn.Close() //close
	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		if len(b) > 0 && b[len(b)-1] == '\r' {
			b = b[:len(b)-1]
		}
		if err := v1h.Exec(b); err != nil {
		} else {

		}
	}
}

func (v1h *V1Handle) Exec(line []byte) error {
	c, para := common.ParseCommand(line)
	var e error
	if bytes.Equal(c, common.COMMAND_GET) {
		e = v1h.GET(para)
	} else if bytes.Equal(c, common.COMMAND_PUT) {
		e = v1h.PUT(para)
	} else if bytes.Equal(c, common.COMMAND_PING) {
		v1h.Ping()
	} else { //
		v1h.WriteError([]byte(e.Error()))
	}
	return e

}

func (v1h *V1Handle) WriteError(reason []byte) (int64, error) {
	v1h.Write(common.E_ERROR, [][]byte{reason}, nil)
	return 0, nil
}
func (v1h *V1Handle) GET(para [][]byte) error {
	if len(para) != 2 {
		v1h.WriteError([]byte("must be 2 parameters"))
		return nil
	}
	groupname := string(para[1])
	key := string(para[1])
	group := groups.getGroup(groupname)
	v := group.Get(key)
	if v == nil {
		v1h.WriteError(common.E_NOT_FOUND)
		return nil
	}
	return nil
}

/*
	para:groupname key
*/
func (v1h *V1Handle) PUT(para [][]byte) error {
	if len(para) != 2 {
		v1h.WriteError([]byte("must be 2 parameters"))
		return nil
	}
	groupname := string(para[1])
	key := string(para[1])
	group := groups.getGroup(groupname)
	length_bytes := make([]byte, 4)
	_, err := io.ReadFull(v1h.conn, length_bytes)
	if err != nil {
		v1h.WriteError([]byte(err.Error()))
		return nil
	}
	length := binary.BigEndian.Uint32(length_bytes)
	body := make([]byte, length)
	_, err = io.ReadFull(v1h.conn, body)
	if err != nil {
		v1h.WriteError([]byte(err.Error()))
		return nil
	}
	group.Put(key, common.BytesData(body))
	v1h.Write(common.OK, nil, nil)
	return nil
}

//thread safe wirte method,
func (v1h *V1Handle) Write(command []byte, parameter [][]byte, content []byte) (int, error) {
	//	v1h.Lock()
	//	defer v1h.Unlock()
	return common.NewCommand(command, parameter, content).Write(v1h.conn)
}

func newTcpHandle() TcpHandle {
	return &V1Handle{}
}
