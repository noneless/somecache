package tcphandle

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/756445638/somecache/common"
)

func Server(ln net.Listener) error {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		versionbytes := make([]byte, 4)
		n, err := conn.Read(versionbytes)
		if err != nil || n != 4 {
			fmt.Println("accept error:", err)
			continue
		}
		if bytes.Equal(versionbytes, common.MagicV1) {
			go (&V1Handler{}).IOLoop(conn)
		} else {
			fmt.Println("unsupport protocol version")
		}

	}
}

type V1Handler struct {
}

func (v1h *V1Handler) IOLoop(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		if b[len(b)-1] == '\r' {
			b = b[:len(b)-1]
		}

	}
}

func (v1h *V1Handler) Exec()
