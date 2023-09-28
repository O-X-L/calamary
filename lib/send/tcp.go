package send

import (
	"io"
	"net"

	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/parse"
)

func forwardTcp(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter) {
	var connFwd net.Conn
	var err error
	connFwd, err = net.Dial("tcp", parse.PktDest(pkt))

	if err != nil {
		log.ConnError("send", parse.PktSrc(pkt), parse.PktDest(pkt), err)
		return
	}
	defer connFwd.Close()

	link := Link{src: connIo}
	log.ConnDebug(
		"send", parse.PktSrc(pkt), parse.PktDest(pkt),
		"Forwarding",
	)
	go link.pipe(pkt, connIo, connFwd)
	go link.pipe(pkt, connFwd, connIo)
	<-link.closed
}
