package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"

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
	sync.Mutex
}

func HandleTcp(ln net.Listener) error {
	wg.Add(1)
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
		if bytes.Equal(b, common.MagicV1) { //the client should first send communication first
			go (&V1Handle{}).IOLoop(conn) //driver by read
		} else {
			fmt.Println("un support version", string(b))
		}
	}
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
	c, para := v1h.parseCommand(line)
	var e error
	if bytes.Equal(c, common.COMMAND_GET) {
		e = v1h.GET(para)
	} else if bytes.Equal(c, common.COMMAND_PUT) {
		e = v1h.PUT(para)
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
	bs, err := io.ReadFull(v1h.conn, body)
	if err != nil {
		v1h.WriteError([]byte(err.Error()))
		return nil
	}
	return nil
}

func (v1h *V1Handle) parseCommand(line []byte) ([]byte, [][]byte) {
	t := bytes.Split(line, common.WhiteSpace)
	para := [][]byte{}
	for i := 1; i < len(t)-1; i++ {
		para = append(para, t[i])
	}
	return t[0], para
}

//thread safe wirte method
func (v1h *V1Handle) Write(command []byte, parameter [][]byte, content []byte) (int, error) {
	v1h.Lock()
	defer v1h.Unlock()
	total := int(0)
	n, err := v1h.conn.Write(command)
	total += n
	if err != nil {
		return total, err
	}
	n, err = v1h.conn.Write(bytes.Join(parameter, common.WhiteSpace))
	total += n
	if err != nil {
		return total, err
	}
	length := len(content)
	for length > 0 {
		n, err = v1h.conn.Write(content)
		total += n
		if err != nil {
			return total, err
		}
		length -= n
		content = content[n:]
	}
	return total, nil
}

func handlemaster() {

}

func newTcpHandle() TcpHandle {
	return &V1Handle{}
}
