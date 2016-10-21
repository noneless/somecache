package tcphandle

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/756445638/somecache/common"
)

type V1Slave struct {
	conn   net.Conn
	reader *bufio.Reader
	lock   sync.Mutex
	ctx    context
}

//main routine
func (v1s *V1Slave) CommandLoop() {
	defer v1s.conn.Close()
	v1s.reader = bufio.NewReader(v1s.conn)
	pingticker := time.NewTicker(time.Second)
	for {
		select {
		case <-pingticker.C:
			err := v1s.Ping()
			if err != nil {
				fmt.Println("ping failed,err:", err)
				return
			}
		}
	}
}
func (v1s *V1Slave) readLine() ([]byte, error) {
	line, err := v1s.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	line = line[0 : len(line)-1]
	fmt.Println("read line:", string(line))
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[0 : len(line)-1]
	}
	return line, nil
}

func (v1s *V1Slave) Ping() error {
	v1s.lock.Lock()
	defer v1s.lock.Unlock()
	v1s.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	_, err := common.NewCommand(common.COMMAND_PING, nil, nil).Write(v1s.conn)
	if err != nil {
		return err
	}
	v1s.conn.SetReadDeadline(time.Now().Add(readTimeout))
	res, err := v1s.readLine()
	if err != nil {
		return err
	}
	if !bytes.Equal(res, common.OK) {
		return fmt.Errorf("slave send response but not OK")
	}
	return nil
}
