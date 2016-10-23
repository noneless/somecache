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
			slave:     slave,
			closechan: make(chan struct{}),
			jobschan:  make(chan *job),
		}
		return v1, nil

	} else {
		return nil, errors.New("unkown version")
	}
}

type V1Slave struct {
	stoped    bool
	t         int64
	conn      net.Conn
	reader    *bufio.Reader
	ctx       context
	slave     *Slave
	closechan chan struct{}
	jobschan  chan *job
}

type job struct {
	c         *common.Command
	diff      interface{} //diff is a field can receive any kind of data
	errorchan chan error
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
		case d := <-v1s.jobschan:
			v1s.t = time.Now().UnixNano()
			v1s.exec(d)
			v1s.t = -1
		case <-v1s.closechan:
			goto exit
		}
	}
exit:
	close(v1s.closechan)
	close(v1s.jobschan)
}

func (v1s *V1Slave) exec(d *job) {
	var e error
	if bytes.Equal(d.c.Command, common.COMMAND_GET) {

	} else if bytes.Equal(d.c.Command, common.COMMAND_GET_STREAM) {

	} else {
		e = fmt.Errorf("no such command")
	}
	d.errorchan <- e
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

func (v1s *V1Slave) Get(key string, dest *[]byte) error {
	errorchan := make(chan error)
	v1s.jobschan <- &job{
		c:         common.NewCommand(common.COMMAND_GET, [][]byte{[]byte(key)}, nil),
		diff:      dest,
		errorchan: errorchan,
	}
	var e error
	select {
	case e = <-errorchan:
	}
	return e
}

func (v1s *V1Slave) Transfer2Writer(key string, w io.Writer) error {
	return nil
}
