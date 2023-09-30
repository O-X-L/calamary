package rcv

import (
	"context"
	"fmt"
	"net"
	"syscall"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
	"golang.org/x/sys/unix"
)

type listenerTransparentTcp struct {
	ln    net.Listener
	addr  string
	Lncnf cnf.ServiceListener
}

func newServerTransparentTcp(addr string, lncnf cnf.ServiceListener) (Server, error) {
	lc := net.ListenConfig{}
	ln, err := lc.Listen(
		context.Background(),
		"tcp",
		fmt.Sprintf("%v:%v", addr, lncnf.Port),
	)
	if err != nil {
		return Server{}, err
	}

	return Server{
		Listener: &listenerTransparentTcp{ln, addr, lncnf},
		Cnf:      lncnf,
		L4Proto:  meta.ProtoL4Tcp,
	}, nil
}

func (l *listenerTransparentTcp) Accept() (conn net.Conn, err error) {
	return l.ln.Accept()
}

func (l *listenerTransparentTcp) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *listenerTransparentTcp) Close() error {
	return l.ln.Close()
}

func (l *listenerTransparentTcp) control(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		if l.Lncnf.TProxy {
			if err := unix.SetsockoptInt(int(fd), unix.SOL_IP, unix.IP_TRANSPARENT, 1); err != nil {
				log.ErrorS("listener-tcp", fmt.Sprintf("SetsockoptInt(SOL_IP, IP_TRANSPARENT, 1): %v", err))
			}
		}
	})
}
