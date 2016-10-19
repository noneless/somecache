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

type V1Handle struct {
	conn net.Conn
	sync.Mutex
}

func HandleTcp(ln net.Listener) error {
	wg.Add(1)
	defer wg.Done()
	b := make([]byte, 4)
	for {
		conn, err := th.ln.Accept()
		if err != nil {
			return err
		}
		n, err := io.ReadFull(conn, b)
		if err != nil && n != 4 {
			continue
		}
		if bytes.Equal(b, common.MagicV1) { //the client should first send v1
			go (&V1Handle{conn: conn}).ReadLoop() //driver by read
		} else {
			fmt.Println("un support version", string(b))
		}
	}
}

func (v1h *V1Handle) ReadLoop() {
	reader := bufio.NewReader(v1h.conn)
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
func (v1h *V1Handle) Exec(line []byte) {
	c, para := v1h.parseCommand(line)
	var e error
	if bytes.Equal(c, common.COMMAND_GET) {
		e = v1h.GET(para)
	} else if bytes.Equal(c, common.COMMAND_PUT) {
		e = v1h.PUT(para)
	} else { //
		v1h.WriteError([]byte(e.Error()))
	}

}

func (v1h *V1Handle) WriteError(reason []byte) (int64, error) {
	v1h.Write(common.E_ERROR)
}
func (v1h *V1Handle) GET() error {

}
func (v1h *V1Handle) PUT() error {

}

func (v1h *V1Handle) parseCommand(line []byte) ([]byte, [][]byte) {
	t := bytes.Split(line, common.WhiteSpace)
	para := [][]byte{}
	for i = 1; i < len(t)-1; i++ {
		para = append(para, t[i])
	}
	return t[0], para
}

func (v1h *V1Handle) Write(command []byte, parameter [][]byte, content []byte) error {
	v1h.Lock()
	defer v1h.Unlock()
}

func handlemaster() {

}

func newTcpHandle() TcpHandle {
	return &TcpV1Handle{}
}
