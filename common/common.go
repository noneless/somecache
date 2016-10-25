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
package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	MagicV1                 = []byte("  V1")
	COMMAND_PUT             = []byte("PUT")
	COMMAND_PUT_FROM_READER = []byte("PUT_FROM_READER")
	COMMAND_GET             = []byte("GET")
	COMMAND_GET_STREAM      = []byte("GET_STREAM") //stream type get
	COMMAND_PING            = []byte("PING")
	COMMAND_LOGIN           = []byte("LOGIN")
	E_ERROR                 = []byte("E_ERROR")
	E_PARAMETER_ERROR       = []byte("E_PARAMETER_ERROR")
	E_READ_ERROR            = []byte("E_READ_ERROR")
	E_NOT_FOUND             = []byte("NOT_FOUND")
	E_COMMAND_NOT_FOUND     = []byte("E_COMMAND_NOT_FOUND")
	ENDL                    = []byte("\n")
	OK                      = []byte("OK")
	WhiteSpace              = []byte(" ")
)

func ParseCommand(line []byte) ([]byte, [][]byte) {
	t := bytes.Split(line, WhiteSpace)
	para := [][]byte{}
	para = append(para, t[1:]...)
	return t[0], para
}

var packerror error = errors.New("pack error")

var EmptyLineError = errors.New("empty line")

func ReadLine(reader *bufio.Reader) (uint64, []byte, error) {
	b := make([]byte, 8)
	n, err := io.ReadFull(reader, b)
	if err != nil {
		return 0, nil, err
	}
	if n != 8 {
		return 0, nil, fmt.Errorf("wrong size")
	}
	jobid := binary.BigEndian.Uint64(b)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		return 0, nil, err
	}
	if len(line) == 0 {
		return 0, nil, EmptyLineError
	}
	line = line[0 : len(line)-1]
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[0 : len(line)-1]
	}
	return jobid, line, nil
}

/*
	read 4 bytes body,body length in 4 bytes bigendian
	if w is not nil,will copy write to w directy,this make it easy to use
*/
func ReadBody4(reader *bufio.Reader, w io.Writer) ([]byte, int, error) {
	l := make([]byte, 4)
	n, err := io.ReadFull(reader, l)
	if err != nil {
		return nil, n, err
	}
	length := binary.BigEndian.Uint32(l)
	if w == nil {
		buf := make([]byte, length)
		n, err = io.ReadFull(reader, buf)
		return buf, n, err
	} else {
		x, err := io.CopyN(w, reader, int64(length))
		return nil, int(x), err
	}
}

func Uint642byte(d uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, d)
	return b
}

type TimeoutError struct {
	Err string
}

func (t *TimeoutError) Error() string {
	return t.Err
}

var TimeOutErr *TimeoutError = &TimeoutError{"timeout"}
