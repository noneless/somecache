package master

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/756445638/somecache/common"
	"github.com/756445638/somecache/message"
)

func newVersionHandler(v []byte, slave *Slave) (ProtocolHandler, error) {
	if bytes.Equal(v, common.MagicV1) {
		v1 := &V1Slave{
			slave:       slave,
			closechan:   make(chan struct{}),
			commandchan: make(chan *commandFn),
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
	stoped      bool
	t           int64
	conn        net.Conn
	reader      *bufio.Reader
	ctx         context
	slave       *Slave
	closechan   chan struct{}
	commandchan chan *commandFn
}

type commandFn struct {
	c  *common.Command
	fn func(reader io.Reader, size int) error
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
			v1s.t = time.Now().UnixNano()
			err := v1s.Ping()
			v1s.t = -1
			if err != nil {
				fmt.Println("ping failed,err:", err)
				return
			}
		case d := <-v1s.commandchan:
			v1s.t = time.Now().UnixNano()
			v1s.exec(d.c, d.fn)
			v1s.t = -1
		case <-v1s.closechan:
			break
		}
	}
}

func (v1s *V1Slave) exec(c *common.Command, fn func(reader io.Reader, size int) error) {

}

//ping is  hearbeat
func (v1s *V1Slave) Ping() error {
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

//-1 means not busy,positive numbre menas how long I have been busy
func (v1s *V1Slave) IfBusy() int64 {
	if v1s.t == -1 {
		return -1
	}
	return time.Now().UnixNano() - v1s.t
}

func (v1s *V1Slave) Exec(c *common.Command, fn func(reader io.Reader, size int) error) {
	v1s.commandchan <- &commandFn{
		c:  c,
		fn: fn,
	}
}
