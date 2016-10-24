package master

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/756445638/somecache/common"
	"github.com/756445638/somecache/message"
)

type V1Slave struct {
	reader    *bufio.Reader
	stoped    bool
	t         int64
	lock      sync.Mutex
	conn      net.Conn
	ctx       context
	slave     *Slave
	closechan chan struct{}
	jobschan  chan *job
	pingpool  sync.Pool
	jobid     uint64
	notify    map[uint64]*job
	wg        sync.WaitGroup
}

type job struct {
	jobid     uint64
	c         *common.Command
	diff      interface{} //diff is a field can receive any kind of data
	errorchan chan error  //errorchan is also an endchanv that will notify caller to exit
}

func (v1s *V1Slave) close() {
	v1s.lock.Lock()
	defer v1s.lock.Unlock()
	v1s.stoped = true
	v1s.conn.Close()
	close(v1s.closechan)
	close(v1s.jobschan)
}

func (v1s *V1Slave) Login(c chan bool) error {
	line, err := common.ReadLine(v1s.reader)
	if err != nil {
		return err
	}
	if !bytes.Equal(line, common.COMMAND_LOGIN) {
		return fmt.Errorf("first package must be login")
	}
	body, _, err := common.ReadBody4(v1s.reader, nil)
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
	return nil
}

//main loop
func (v1s *V1Slave) MainLoop(conn net.Conn, c chan bool) {
	v1s.conn = conn
	defer v1s.close()
	v1s.reader = bufio.NewReader(v1s.conn)
	err := v1s.Login(c)
	c <- (err == nil)
	if err != nil {
		fmt.Println("login failed,err:", err)
		return
	}
	v1s.wg.Add(3)
	go func() {
		defer v1s.wg.Done()
		e := v1s.writeLoop()
		if e != nil {
			fmt.Println("loop error:", e)
		}
	}()
	go func() {
		defer v1s.wg.Done()
		e := v1s.readLoop()
		if e != nil {
			fmt.Println("loop error:", e)
		}
	}()
	go func() {
		defer v1s.wg.Done()
		e := v1s.pingLoop()
		if e != nil {
			fmt.Println("loop error:", e)
		}
	}()
	v1s.wg.Wait()
}

func (v1s *V1Slave) writeLoop() error {
	for job := range v1s.jobschan {
		_, err := job.c.Write(v1s.conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v1s *V1Slave) readLoop() error {
	v1s.reader = bufio.NewReader(v1s.conn)
	for {
		line, err := common.ReadLine(v1s.reader)
		if err != nil {
			fmt.Println("read error:", err)
			return err
		}
		err = v1s.processRead(line)
		if err != nil {
			fmt.Println("process line error:", err)
		}
	}
	return nil
}
func (v1s *V1Slave) processRead(line []byte) error {
	var err error
	c, jobid, _, err := common.ParseCommandJobid(line)
	if err != nil {
		return err
	}
	job, ok := v1s.notify[jobid]
	if !ok {
		return fmt.Errorf("can`t find job binded on jobid")
	}

	defer func(e *error) {
		job.errorchan <- *e
		delete(v1s.notify, jobid)
	}(&err)
	if !bytes.Equal(c, common.OK) {
		err = fmt.Errorf("slave return error:", string(c))
		return err
	}
	if bytes.Equal(job.c.Command, common.COMMAND_GET) {
		dest := job.diff.(*[]byte)
		*dest, _, err = common.ReadBody4(v1s.reader, nil)
	} else if bytes.Equal(job.c.Command, common.COMMAND_PUT) {
		err = nil
	} else {
		err = fmt.Errorf("unkown command")
	}

	return err
}

func (v1s *V1Slave) pingLoop() error {
	for {
		time.Sleep(time.Second)
		v1s.ping()
	}
	return nil
}
func (v1s *V1Slave) Close() {
	v1s.closechan <- struct{}{}
}

func (v1s *V1Slave) Get(key string) ([]byte, error) {
	var b []byte
	errorchan := make(chan error)
	defer close(errorchan)
	jobid, b := v1s.newJobId()
	job := &job{
		c:         common.NewCommand(common.COMMAND_GET, [][]byte{b, []byte(key)}, nil),
		diff:      &b,
		errorchan: errorchan,
		jobid:     jobid,
	}
	v1s.notify[jobid] = job
	v1s.jobschan <- job
	var e error
	select {
	case e = <-errorchan:
	}
	return b, e
}

func (v1s *V1Slave) newJobId() (uint64, []byte) {
	jobid := atomic.AddUint64(&v1s.jobid, 1)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, jobid)
	return jobid, b
}

func (v1s *V1Slave) Put(key string, data []byte) error {
	errorchan := make(chan error)
	defer close(errorchan)
	jobid, b := v1s.newJobId()
	job := &job{
		c:         common.NewCommand(common.COMMAND_PUT, [][]byte{b, []byte(key)}, data),
		diff:      nil,
		errorchan: errorchan,
		jobid:     jobid,
	}
	v1s.notify[jobid] = job
	v1s.jobschan <- job
	var e error
	select {
	case e = <-errorchan:
	}
	return e
}
func (v1s *V1Slave) ping() error {
	errorchan := make(chan error)
	defer close(errorchan)
	jobid, b := v1s.newJobId()
	job := &job{
		c:         common.NewCommand(common.COMMAND_PING, [][]byte{b}, nil),
		diff:      nil,
		errorchan: errorchan,
		jobid:     jobid,
	}
	v1s.notify[jobid] = job
	v1s.jobschan <- job
	var e error
	select {
	case e = <-errorchan:
	}
	return e
}

var stopped = errors.New("slave stoped")
