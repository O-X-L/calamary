package rcv

import (
	"net"
)

type Listener interface {
	Init() error
	Accept() (net.Conn, error)
	Addr() net.Addr
	Close() error
}
