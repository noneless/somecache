package main

import (
	"fmt"
	"net"
	"time"

	"github.com/756445638/somecache/common"
)

func Connection2Master(tcp_addr string) {
	wg.Add(1)
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
		v1c := &V1Client{}
		go v1c.IOLoop(conn)
	}
}

type V1Client struct {
	conn net.Conn
}

func (v1c *V1Client) IOLoop(conn net.Conn) {
	v1c.conn = conn
	defer conn.Close()

}
