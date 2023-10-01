package send

import (
	"io"
	"net"
	"net/http"

	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

func Forward(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter) {
	if pkt.L4.Proto == meta.ProtoL4Udp {
		// connFwd, err = DialUdp(pkt)
		parse.LogConnError("send", pkt, "UDP not yet implemented!")

	} else {
		forwardTcp(pkt, conn, connIo)
	}
}

func ForwardHttp(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter, req *http.Request) {
	// http-proxy - rewrite request; only for plaintext http
	connFwd := establishTcp(pkt)

	if err := req.Write(connFwd); err != nil {
		return
	}

	transportTcp(pkt, conn, connIo, connFwd)
}
