package send

import (
	"io"
	"net"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/parse"
	"github.com/superstes/calamary/u"
)

func forwardTcp(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter) {
	var connFwd net.Conn
	var err error
	connFwd, err = net.DialTimeout("tcp", parse.PktDest(pkt), u.Timeout(cnf.C.Service.Timeout.Connect))

	if err != nil {
		log.ConnError("send", parse.PktSrc(pkt), parse.PktDest(pkt), err)
		return
	}
	defer connFwd.Close()

	close := make(chan bool, 1)
	link := Link{src: connIo, close: close}
	log.ConnDebug("send", parse.PktSrc(pkt), parse.PktDest(pkt), "Forwarding")
	go link.pipe(pkt, conn, connIo, connFwd)
	go link.pipe(pkt, connFwd, connFwd, connIo)
	<-close
	log.ConnDebug("send", parse.PktSrc(pkt), parse.PktDest(pkt), "Closed")
}
