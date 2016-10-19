package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"net"

	"github.com/756445638/somecache/common"
)

type V1Handle struct {
	conn net.Conn
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

		if bytes.Equal(b, common.MagicV1) {
			go (&V1Handle).IOLoop()
		}
		go th.handleConn(conn)
	}
}

func (th *V1Handle) handleConn(conn net.Conn) {

}

func handlemaster() {

}

func newTcpHandle() TcpHandle {
	return &TcpV1Handle{}
}
