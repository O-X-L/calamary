package send

import (
	"io"
	"net"

	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

func Forward(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter) {
	if pkt.L4.Proto == meta.ProtoL4Udp {
		// connFwd, err = DialUdp(pkt)
		log.ConnErrorS("send", parse.PktSrc(pkt), parse.PktDest(pkt), "UDP not yet implemented!")

	} else {
		forwardTcp(pkt, conn, connIo)
	}
}
