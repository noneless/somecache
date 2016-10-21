package common

import (
	"bufio"
	"bytes"
	"errors"
	"net"
)

type Getter interface {
	Get(string) ([]byte, error)
}

type TcpServer interface {
	TcpServer(net.Listener)
}

var (
	MagicV1      = []byte("  V1")
	COMMAND_PUT  = []byte("PUT")
	COMMAND_GET  = []byte("GET")
	COMMAND_PING = []byte("PING")
	E_ERROR      = []byte("E_ERROR")
	E_NOT_FOUND  = []byte("NOT_FOUND")
	ENDL         = []byte("\n")
	OK           = []byte("OK")
	WhiteSpace   = []byte(" ")
)

func ParseCommand(line []byte) ([]byte, [][]byte) {
	t := bytes.Split(line, WhiteSpace)
	para := [][]byte{}
	for i := 1; i < len(t)-1; i++ {
		para = append(para, t[i])
	}
	return t[0], para
}

var EmptyLineError = errors.New("empty line")

func ReadLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, EmptyLineError
	}
	line = line[0 : len(line)-1]
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[0 : len(line)-1]
	}
	return line, nil
}
