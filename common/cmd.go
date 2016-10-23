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
}

//not thread safe
func (c *Command) Write(w io.Writer) (int, error) {
	total := int(0)
	n, err := w.Write(c.Command)
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
		n, err = w.Write(bytes.Join(c.Parameter, WhiteSpace))
		total += n
		if err != nil {
			return total, err
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

func NewCommand(command []byte, paras [][]byte, content []byte) *Command {
	return &Command{command, paras, content}
}
