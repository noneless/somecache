package master

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

func newVersionHandler(v []byte, slave *Slave) (ProtocolHandler, error) {
	if bytes.Equal(v, common.MagicV1) {
		v1 := &V1Slave{
			slave:     slave,
			closechan: make(chan struct{}),
		}
		return v1, nil

	} else {
		return nil, errors.New("unkown version")
	}
}

type Slave struct {
	addr         net.Addr
	service      *Service
	handle       ProtocolHandler
	loginmessage *message.Login
}

type V1Slave struct {
	conn      net.Conn
	reader    *bufio.Reader
	lock      sync.Mutex
	ctx       context
	slave     *Slave
	closechan chan struct{}
}

func (v1s *V1Slave) Login(c chan struct{}) error {
	line, err := common.ReadLine(v1s.reader)
	if err != nil {
		return err
	}
	if !bytes.Equal(line, common.COMMAND_LOGIN) {
		return fmt.Errorf("first package must be login")
	}
	body, _, err := common.Read4BytesBody(v1s.reader)
	if err != nil {
		return err
	}
	login := &message.Login{}
	err = json.Unmarshal(body, login)
	if err != nil {
		return err
	}
	v1s.slave.loginmessage = login
	_, err = common.NewCommand(common.OK, nil, nil).Write(v1s.conn)
	if err != nil {
		return err
	}
	c <- struct{}{}

	return nil
}

//main loop
func (v1s *V1Slave) MainLoop(conn net.Conn, c chan struct{}) {
	v1s.conn = conn
	defer v1s.conn.Close()
	v1s.reader = bufio.NewReader(v1s.conn)
	if err := v1s.Login(c); err != nil {
		fmt.Println("login failed,err:", err)
		return
	}
	pingticker := time.NewTicker(time.Second)
	for {
		select {
		case <-pingticker.C:
			err := v1s.Ping()
			if err != nil {
				fmt.Println("ping failed,err:", err)
				return
			}
		case <-v1s.closechan:
			break
		}
	}
}

//ping is  hearbeat
func (v1s *V1Slave) Ping() error {
	v1s.lock.Lock()
	defer v1s.lock.Unlock()
	v1s.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	_, err := common.NewCommand(common.COMMAND_PING, nil, nil).Write(v1s.conn)
	if err != nil {
		return err
	}
	v1s.conn.SetReadDeadline(time.Now().Add(readTimeout))
	res, err := common.ReadLine(v1s.reader)
	if err != nil {
		return err
	}
	if !bytes.Equal(res, common.OK) {
		return fmt.Errorf("slave send response but not OK")
	}
	fmt.Printf("slave[%s] process ping ok\n", v1s.slave.addr.String())
	return nil
}

func (v1s *V1Slave) Close() {
	v1s.closechan <- struct{}{}
}
