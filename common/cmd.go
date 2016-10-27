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
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Command struct {
	Command   []byte
	Parameter [][]byte
	Content   []byte
	Jobid     uint64
}

//not thread safe
func (c *Command) Write(w io.Writer) (int, error) {
	total := int(0)

	jb := Uint642byte(c.Jobid)
	n, err := w.Write(jb)
	total += n
	if err != nil {
		return total, err
	}

	n, err = w.Write(c.Command)
	total += n
	if err != nil {
		return total, err
	}
	//parameter
	if c.Parameter != nil && len(c.Parameter) > 0 {
		n, err = w.Write(WhiteSpace)
		total += n
		if err != nil {
			return total, err
		}
		bs := bytes.Join(c.Parameter, WhiteSpace)
		n, err = w.Write(bytes.Join(c.Parameter, WhiteSpace))
		total += n
		if err != nil {
			return total, err
		}
		if n != len(bs) {
			return total, fmt.Errorf("not enough length")
		}
	}
	// \n must have to do this
	n, err = w.Write(ENDL)
	total += n
	if err != nil {
		return total, err
	}

	if c.Content != nil && len(c.Content) > 0 {
		length := len(c.Content)
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(length))
		n, err = w.Write(b)
		total += n
		if err != nil {
			return total, err
		}
		if n != 4 {
			return total, fmt.Errorf("length must be 4")
		}

		for length > 0 {
			n, err = w.Write(c.Content)
			total += n
			if err != nil {
				return total, err
			}
			length -= n
			c.Content = c.Content[n:]
		}
	}
	return total, nil
}

func NewCommand(command []byte, paras [][]byte, content []byte, jobid uint64) *Command {
	return &Command{command, paras, content, jobid}
}
