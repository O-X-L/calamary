package rcv

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/u"
	"golang.org/x/sys/unix"
)

type listenerTransparentUdp struct {
	ln    *net.UDPConn
	addr  string
	Lncnf cnf.ServiceListener
}

func newServerTransparentUdp(addr string, lncnf cnf.ServiceListener) (Server, error) {
	return Server{}, fmt.Errorf("UDP listener is not yet implemented!")
	/*
		lnu := &listenerTransparentUdp{
			addr:  addr,
			lncnf: lncnf,
		}
		ln, err := lnu.listenUdp(addr)
		if err != nil {
			return nil, err
		}
		lnu.ln = ln
		return lnu, nil
	*/
}

func (l *listenerTransparentUdp) Accept() (conn net.Conn, err error) {
	conn, err = l.accept()
	if err != nil {
		return
	}
	return
}

func (l *listenerTransparentUdp) Addr() net.Addr {
	return l.ln.LocalAddr()
}

func (l *listenerTransparentUdp) Close() error {
	return l.ln.Close()
}

type connUdp struct {
	net.Conn
	buf  []byte
	ttl  time.Duration
	once sync.Once
}

func (c *connUdp) Read(b []byte) (n int, err error) {
	if c.ttl > 0 {
		c.SetReadDeadline(time.Now().Add(c.ttl))
		defer c.SetReadDeadline(time.Time{})
	}

	c.once.Do(func() {
		n = copy(b, c.buf)
		u.PutBufferPool(&c.buf)
	})

	if n == 0 {
		n, err = c.Conn.Read(b)
	}

	return
}

func (c *connUdp) Write(b []byte) (n int, err error) {
	if c.ttl > 0 {
		c.SetWriteDeadline(time.Now().Add(c.ttl))
		defer c.SetWriteDeadline(time.Time{})
	}
	return c.Conn.Write(b)
}

func (l *listenerTransparentUdp) listenUdp(addr string) (*net.UDPConn, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				if l.Lncnf.TProxy {
					if err := unix.SetsockoptInt(int(fd), unix.SOL_IP, unix.IP_TRANSPARENT, 1); err != nil {
						log.ErrorS("listener-udp", fmt.Sprintf("SetsockoptInt(SOL_IP, IP_TRANSPARENT, 1): %v", err))
					}
				}
				if err := unix.SetsockoptInt(int(fd), unix.SOL_IP, unix.IP_RECVORIGDSTADDR, 1); err != nil {
					log.ErrorS("listener-udp", fmt.Sprintf("SetsockoptInt(SOL_IP, IP_RECVORIGDSTADDR, 1): %v", err))
				}
			})
		},
	}

	network := "udp"
	if u.IsIPv4(addr) {
		network = "udp4"
	}
	pc, err := lc.ListenPacket(
		context.Background(),
		network,
		fmt.Sprintf("%v:%v", addr, l.Lncnf.Port),
	)
	if err != nil {
		return nil, err
	}

	return pc.(*net.UDPConn), nil
}

func (l *listenerTransparentUdp) accept() (conn net.Conn, err error) {
	panic(fmt.Errorf("processing of UDP traffic is not yet implemented"))
	/*
		b := u.GetBufferPool(cnf.UDP_BUFFER_SIZE)

		_, raddr, dstAddr, err := readAdressesFromUDP(l.ln, *b)
		logSrc := raddr.String()
		logDst := dstAddr.String()
		if err != nil {
			if !strings.Contains(fmt.Sprintf("%v", err), "use of closed network connection") {
				log.ConnError("listener-udp", logSrc, logDst, err)
			}
			return
		}

		log.ConnInfo("listener-udp", logSrc, logDst, "establishing connection")

		network := "udp"
		if u.IsIPv4(l.addr) {
			network = "udp4"
		}
		freePort, err := net.ResolveUDPAddr(network, "localhost:0")
		if err != nil {
			log.ConnErrorS("listener-udp", logSrc, logDst, "Unable to get free port")
			return
		}
		conn, err = net.ListenUDP(network, freePort)
		if err != nil {
			log.ConnErrorS("listener-udp", logSrc, logDst, "Unable to get free port")
			return
		}
	*/
	/*
		conn = &connUdp{
			Conn: c,
			buf:  (*b)[:n],
		}

		c, err := nw.DialUDP(network, dstAddr, raddr)
		if err != nil {
			if strings.Contains(u.ToStr(err), "address already in use") && strings.Contains(u.ToStr(err), u.ToStr(cnf.ListenPort)) {
				err = fmt.Errorf("Denied connection targeting proxy directly")
			}
			log.ConnError("listener-udp", logSrc, logDst, err)
			return
		}

		conn = &connUdp{
			Conn: c,
			buf:  (*b)[:n],
		}
	*/
	return l.ln, nil
}

type RemoteAddr interface {
	RemoteAddr() net.Addr
}

type SetBuffer interface {
	SetReadBuffer(bytes int) error
	SetWriteBuffer(bytes int) error
}

type SyscallConn interface {
	SyscallConn() (syscall.RawConn, error)
}

type Conn interface {
	net.PacketConn
	io.Reader
	io.Writer
	ReadUDP
	WriteUDP
	SetBuffer
	SyscallConn
	RemoteAddr
}

type ReadUDP interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *net.UDPAddr, err error)
}

type WriteUDP interface {
	WriteToUDP(b []byte, addr *net.UDPAddr) (int, error)
	WriteMsgUDP(b, oob []byte, addr *net.UDPAddr) (n, oobn int, err error)
}
