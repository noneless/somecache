package abstract

import (
	"net"
)

type RunMaster interface {
	RunMaster(net.Listener) error
}
