package rcv

import (
	"context"
	"fmt"
	"net"
	"syscall"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/u"
	"golang.org/x/sys/unix"
)

func init() {}

type listenerTcp struct {
	ln   net.Listener
	addr string
}

func NewListenerTcp(addr string) (Listener, error) {
	lc := net.ListenConfig{}
	/* if l.md.tproxy {
		lc.Control = log.Control
	}*/
	network := "tcp"
	if u.IsIPv4(addr) {
		network = "tcp4"
	}
	ln, err := lc.Listen(
		context.Background(),
		network,
		fmt.Sprintf("%v:%v", addr, cnf.C.Service.Listen.Port),
	)
	if err != nil {
		return nil, err
	}

	// ln = proxyproto.WrapListener(l.options.ProxyProtocol, ln, 10*time.Second)
	return &listenerTcp{ln, addr}, nil
}

func (l *listenerTcp) Init() (err error) { return }

func (l *listenerTcp) Accept() (conn net.Conn, err error) {
	return l.ln.Accept()
}

func (l *listenerTcp) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *listenerTcp) Close() error {
	return l.ln.Close()
}

func (l *listenerTcp) control(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		if cnf.C.Service.Listen.Transparent {
			if err := unix.SetsockoptInt(int(fd), unix.SOL_IP, unix.IP_TRANSPARENT, 1); err != nil {
				log.ErrorS("listener-tcp", fmt.Sprintf("SetsockoptInt(SOL_IP, IP_TRANSPARENT, 1): %v", err))
			}
		}
	})
}
