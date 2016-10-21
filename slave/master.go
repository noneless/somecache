package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/756445638/somecache/common"
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
		e := v1s.IOLoop()
		if e != nil {
			fmt.Println("IOLoop failed,err:", e)
		}
	}
}

type V1Slave struct {
	conn net.Conn
}

func (v1s *V1Slave) IOLoop() error {
	defer v1s.conn.Close()
	reader := bufio.NewReader(v1s.conn)
	for {
		line, err := common.ReadLine(reader)
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

func (v1s *V1Slave) Exec(line []byte) error { // error just for log
	cmd, _ := common.ParseCommand(line)
	if bytes.Equal(cmd, common.COMMAND_PING) {
		v1s.Ping()
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
	_, err := common.NewCommand(common.E_ERROR, [][]byte{reason}, nil).Write(v1s.conn)
	return err
}
