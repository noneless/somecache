/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package master

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	lock      sync.Mutex
	conn      net.Conn
	ctx       context
	slave     *Slave
	jobschan  chan *job
	pingpool  sync.Pool
	jobid     uint64
	jobs      map[uint64]*job
	jobs_lock sync.Mutex
}

type job struct {
	jobid     uint64
	c         *common.Command
	diff      interface{} //diff is a field can receive any kind of data
	errorchan chan error  //errorchan is also an endchanv that will notify caller to exit
}

func (v1s *V1Slave) Login(c chan bool) error {
	_, line, err := common.ReadLine(v1s.reader)
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
	_, err = common.NewCommand(common.OK, nil, nil, 0).Write(v1s.conn)
	if err != nil {
		return err
	}
	return nil
}

//main loop
func (v1s *V1Slave) MainLoop(conn net.Conn, c chan bool) {
	v1s.conn = conn
	defer v1s.Close()
	v1s.reader = bufio.NewReader(v1s.conn)
	err := v1s.Login(c)
	c <- (err == nil)
	if err != nil {
		log.Println("error:login failed,err:", err)
		return
	}
	errchan := make(chan error)
	go func() {
		errchan <- v1s.writeLoop()
	}()
	go func() {
		errchan <- v1s.readLoop()
	}()
	go func() {
		errchan <- v1s.pingLoop()
	}()
	e := <-errchan
	if err != nil {
		log.Printf("error:err poped:%v\n", e)
	}
}

func sendError(c chan error, err error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Printf("warn: panic[%v] recovered\n", x)
		}
	}()
	c <- err
}

//it is a runtime send timeouterror to errorchan
func (v1s *V1Slave) ticking() {
	timeout := func() {
		for _, v := range v1s.jobs {
			sendError(v.errorchan, common.TimeOutErr)
		}
	}
	timeoutticker := time.NewTicker(time.Second)
	for !v1s.stoped {
		select {
		case <-timeoutticker.C:
			timeout()
		}
	}
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
		jobid, line, err := common.ReadLine(v1s.reader)
		if err != nil {
			return err
		}
		log.Printf("debug: read line[%s]\n", string(line))
		err = v1s.processRead(jobid, line)
		if err != nil {
			log.Printf("warn:process line error:%v\n", err)
		}
	}
	return nil
}
func (v1s *V1Slave) processRead(jobid uint64, line []byte) error {
	defer func() {
		x := recover()
		if x != nil {
			log.Printf("warn:panic[%v] recovered\n", x)
		}
	}()
	var err error
	c, _ := common.ParseCommand(line)
	if err != nil {
		return err
	}
	job, ok := v1s.jobs[jobid]
	if !ok {
		return fmt.Errorf("can`t find job binded on jobid")
	}
	defer func(e *error) {
		sendError(job.errorchan, *e)
		v1s.deljob(jobid)
	}(&err)
	if !bytes.Equal(c, common.OK) {
		return fmt.Errorf("slave return error:", string(c))
	}
	if bytes.Equal(job.c.Command, common.COMMAND_GET) {
		dest := job.diff.(*[]byte)
		*dest, _, err = common.ReadBody4(v1s.reader, nil)
	} else if bytes.Equal(job.c.Command, common.COMMAND_PUT) { //ok is just enough
		err = nil
	} else if bytes.Equal(job.c.Command, common.COMMAND_PING) {
		dest := job.diff.(*[]byte)
		*dest, _, err = common.ReadBody4(v1s.reader, nil)
	} else {
		err = fmt.Errorf("unkown command")
	}
	return err
}

func (v1s *V1Slave) pingLoop() error {
	for !v1s.stoped {
		time.Sleep(time.Second * 3)
		log.Printf("debug: master start to ping\n")

		js, err := v1s.ping()
		if err != nil {
			return err
		}
		log.Printf("info: ping finished,raw message are %v\n", string(js))

	}
	return nil
}
func (v1s *V1Slave) Close() {
	v1s.lock.Lock()
	defer v1s.lock.Unlock()
	v1s.stoped = true
	v1s.conn.Close()
	close(v1s.jobschan)
}

func (v1s *V1Slave) addjob(jobid uint64, j *job) {
	v1s.jobs_lock.Lock()
	v1s.jobs[jobid] = j
	v1s.jobs_lock.Unlock()
}
func (v1s *V1Slave) deljob(jobid uint64) {
	v1s.jobs_lock.Lock()
	delete(v1s.jobs, jobid)
	v1s.jobs_lock.Unlock()
}
func (v1s *V1Slave) selectTimeout(ch chan error) error {
	var err error
	select {
	case <-time.After(time.Second * 2):
		return fmt.Errorf("timeout")
	case err = <-ch:
		return err
	}
	return nil
}

func (v1s *V1Slave) newJobId() uint64 {
	return atomic.AddUint64(&v1s.jobid, 1)
}

func (v1s *V1Slave) Put(key string, data []byte) error {
	v1s.lock.Lock()
	defer v1s.lock.Unlock()
	if v1s.stoped {
		return stoppedError
	}
	errorchan := make(chan error)
	defer close(errorchan)
	jobid := v1s.newJobId()
	j := &job{
		c:         common.NewCommand(common.COMMAND_PUT, [][]byte{[]byte(key)}, data, jobid),
		diff:      nil,
		errorchan: errorchan,
		jobid:     jobid,
	}
	v1s.addjob(jobid, j)
	v1s.jobschan <- j
	return v1s.selectTimeout(errorchan)
}

func (v1s *V1Slave) ping() ([]byte, error) {
	v1s.lock.Lock() // in case ,chan is closed
	defer v1s.lock.Unlock()
	if v1s.stoped {
		return nil, stoppedError
	}
	errorchan := make(chan error)
	defer close(errorchan)
	jobid := v1s.newJobId()
	var dest []byte
	j := &job{
		c:         common.NewCommand(common.COMMAND_PING, nil, nil, jobid),
		diff:      &dest,
		errorchan: errorchan,
		jobid:     jobid,
	}
	v1s.jobschan <- j
	v1s.addjob(jobid, j)
	return dest, v1s.selectTimeout(errorchan)
}

func (v1s *V1Slave) Get(key string) ([]byte, error) {
	v1s.lock.Lock()
	defer v1s.lock.Unlock()
	if v1s.stoped {
		return nil, stoppedError
	}
	var b []byte
	errorchan := make(chan error)
	defer close(errorchan)
	jobid := v1s.newJobId()
	j := &job{
		c:         common.NewCommand(common.COMMAND_GET, [][]byte{[]byte(key)}, nil, jobid),
		diff:      &b,
		errorchan: errorchan,
		jobid:     jobid,
	}
	v1s.addjob(jobid, j)
	v1s.jobschan <- j
	return b, v1s.selectTimeout(errorchan)
}

var stoppedError = errors.New("slave stoped")
