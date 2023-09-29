package send

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/metrics"
	"github.com/superstes/calamary/proc/parse"
	"github.com/superstes/calamary/u"
)

type Link struct {
	sentBytes uint64
	rcvBytes  uint64
	src       io.ReadWriter
	close     chan bool
	closed    bool
}

func (l *Link) pipe(pkt parse.ParsedPacket, srcConn net.Conn, src io.ReadWriter, dst io.ReadWriter) {
	islocal := src == l.src

	buff := make([]byte, 0xffff)
	for {
		srcConn.SetDeadline(time.Now().Add(u.Timeout(cnf.C.Service.Timeout.Idle)))
		if l.closed {
			return
		}
		n, err := src.Read(buff)
		if err != nil {
			l.closer(pkt, "Read failed", err)
			return
		}
		b := buff[:n]

		n, err = dst.Write(b)
		if err != nil {
			l.closer(pkt, "Write failed", err)
			return
		}
		if islocal {
			if cnf.Metrics() {
				metrics.BytesSent.Add(float64(n))
			}
			l.sentBytes += uint64(n)
			log.ConnDebug("send", parse.PktSrc(pkt), parse.PktDest(pkt), fmt.Sprintf("%d bytes sent", l.sentBytes))
		} else {
			if cnf.Metrics() {
				metrics.BytesRcv.Add(float64(n))
			}
			l.rcvBytes += uint64(n)
			log.ConnDebug("send", parse.PktSrc(pkt), parse.PktDest(pkt), fmt.Sprintf("%d bytes received", l.rcvBytes))
		}
	}
}

func (l *Link) closer(pkt parse.ParsedPacket, errMsg string, err error) {
	if l.closed {
		return
	}
	if err != io.EOF {
		log.ConnErrorS(
			"send", parse.PktSrc(pkt), parse.PktDest(pkt),
			fmt.Sprintf("Read failed: %s - %v", errMsg, err),
		)
	}
	l.close <- true
	l.closed = true
}
