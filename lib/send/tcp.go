package send

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/parse"
	"github.com/superstes/calamary/u"
)

func establishTcp(pkt parse.ParsedPacket) net.Conn {
	var connFwd net.Conn
	var err error
	retries := 1
	for {
		connFwd, err = net.DialTimeout("tcp", parse.PktDest(pkt), u.Timeout(cnf.C.Service.Timeout.Connect))

		if err == nil {
			break

		} else {
			if retries >= int(cnf.C.Service.Output.Retries) {
				parse.LogConnError(
					"send", pkt, fmt.Sprintf("Connect retry exceeded (%v/%v): %v",
						retries, cnf.C.Service.Output.Retries, err),
				)
				return nil
			}
			parse.LogConnDebug(
				"send", pkt, fmt.Sprintf("Connection failed (retry %v/%v): %v",
					retries, cnf.C.Service.Output.Retries, err),
			)
			retries++
			time.Sleep(u.Timeout(cnf.DefaultConnectRetryWait))
		}
	}
	parse.LogConnDebug("send", pkt, "Connection established")
	return connFwd
}

func transportTcp(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter, connFwd net.Conn) {
	defer connFwd.Close()
	close := make(chan bool, 1)
	link := Link{src: connIo, close: close}
	parse.LogConnDebug("send", pkt, "Forwarding")
	go link.pipe(pkt, conn, connIo, connFwd)
	go link.pipe(pkt, connFwd, connFwd, connIo)
	<-close
	parse.LogConnDebug("send", pkt, "Closed")
}

func forwardTcp(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter) {
	transportTcp(pkt, conn, connIo, establishTcp(pkt))
}
