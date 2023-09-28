package send

import (
	"fmt"
	"io"

	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/parse"
)

type Link struct {
	sentBytes uint64
	rcvBytes  uint64
	src       io.ReadWriter
	closed    chan bool
	erred     bool
}

func (l *Link) pipe(pkt parse.ParsedPacket, src io.ReadWriter, dst io.ReadWriter) {
	islocal := src == l.src

	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			l.close(pkt, "Read failed", err)
			return
		}
		b := buff[:n]

		n, err = dst.Write(b)
		if err != nil {
			l.close(pkt, "Write failed", err)
			return
		}
		if islocal {
			l.sentBytes += uint64(n)
			log.ConnDebug("send", parse.PktSrc(pkt), parse.PktDest(pkt), fmt.Sprintf("%d bytes sent", l.sentBytes))
		} else {
			l.rcvBytes += uint64(n)
			log.ConnDebug("send", parse.PktSrc(pkt), parse.PktDest(pkt), fmt.Sprintf("%d bytes received", l.rcvBytes))
		}
	}
}

func (l *Link) close(pkt parse.ParsedPacket, errMsg string, err error) {
	if l.erred {
		return
	}
	if err != io.EOF {
		log.ConnErrorS(
			"send", parse.PktSrc(pkt), parse.PktDest(pkt),
			fmt.Sprintf("Read failed: %s - %v", errMsg, err),
		)
	}
	l.closed <- true
	l.erred = true
}
