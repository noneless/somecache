package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/756445638/somecache/common"
	"github.com/756445638/somecache/message"
)

func Connection2Master(tcp_addr string) {
	defer wg.Done()
	for {
		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", tcp_addr)
		if err != nil {
			fmt.Println("dail master server failed,err:", err)
			continue
		}
		_, err = conn.Write(common.MagicV1)
		if err != nil {
			fmt.Println("write magic version v1 failed,err:", err)
			continue
		}
		v1s := &V1Slave{conn: conn}
		e := v1s.MainLoop()
		if e != nil {
			fmt.Println("IOLoop failed,err:", e)
		}
	}
}

type V1Slave struct {
	conn   net.Conn
	reader *bufio.Reader
}

func (v1s *V1Slave) MainLoop() error {
	defer v1s.conn.Close()
	v1s.reader = bufio.NewReader(v1s.conn)
	if err := v1s.Login(); err != nil {
		return fmt.Errorf("login failed,err:%v", err)
	}
	for {
		line, err := common.ReadLine(v1s.reader)
		if err != nil {
			return err
		}
		fmt.Printf("read line,data[%s]\n", string(line))
		err = v1s.Exec(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v1s *V1Slave) Login() error {
	l := message.Login{}
	body, err := json.Marshal(l)
	if err != nil {
		return err
	}
	_, err = common.NewCommand(common.COMMAND_LOGIN, nil, body).Write(v1s.conn)
	if err != nil {
		return err
	}
	line, err := common.ReadLine(v1s.reader)
	if err != nil {
		return err
	}
	if !bytes.Equal(common.OK, line) {
		return fmt.Errorf("master response something[%s] but not ok,", string(line))
	}
	return nil
}

func (v1s *V1Slave) Exec(line []byte) error { // error just for log
	cmd, _ := common.ParseCommand(line)
	if bytes.Equal(cmd, common.COMMAND_PING) {
		return v1s.Ping()
	} else {
		v1s.WtiteError(common.E_NOT_FOUND)
		return errors.New(string(common.E_NOT_FOUND))
	}
	return nil
}

func (v1s *V1Slave) Ping() error {
	_, err := common.NewCommand(common.OK, nil, nil).Write(v1s.conn)
	return err
}

func (v1s *V1Slave) WtiteError(reason []byte) error {
	_, err := common.NewCommand(reason, nil, nil).Write(v1s.conn)
	return err
}
