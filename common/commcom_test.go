package common

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	s := "get 123"
	cmd, p := ParseCommand([]byte(s))
	t.Logf("%s\n", string(cmd))
	t.Logf("%s\n", string(p[0]))
}
