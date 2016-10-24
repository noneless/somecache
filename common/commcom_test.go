package common

import (
	"fmt"
	"testing"
)

func TestParseCommand(t *testing.T) {
	b := []byte{80, 73, 78, 71, 32, 0, 0, 0, 0, 0, 0, 0, 8}
	fmt.Println(ParseCommandJobid(b))
}
