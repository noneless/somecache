package common

import (
	"bytes"
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
	n, err = w.Write(bytes.Join(c.Parameter, WhiteSpace))
	total += n
	if err != nil {
		return total, err
	}
	n, err = w.Write(ENDL)
	total += n
	if err != nil {
		return total, err
	}

	if c.Content != nil && len(c.Content) > 0 {
		length := len(c.Content)
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
