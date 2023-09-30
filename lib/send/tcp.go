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

func forwardTcp(pkt parse.ParsedPacket, conn net.Conn, connIo io.ReadWriter) {
	var connFwd net.Conn
	var err error
	retries := 1
	for {
		connFwd, err = net.DialTimeout("tcp", parse.PktDest(pkt), u.Timeout(cnf.C.Service.Timeout.Connect))

		if err == nil {
			defer connFwd.Close()
			break

		} else {
			if retries >= int(cnf.C.Service.Output.Retries) {
				parse.LogConnError(
					"send", pkt, fmt.Sprintf("Connect retry exceeded (%v/%v): %v",
						cnf.C.Service.Output.Retries, cnf.C.Service.Output.Retries, err),
				)
				return
			}
			parse.LogConnDebug(
				"send", pkt, fmt.Sprintf("Connection failed (retry %v/%v): %v",
					cnf.C.Service.Output.Retries, cnf.C.Service.Output.Retries, err),
			)
			retries++
			time.Sleep(u.Timeout(cnf.DefaultConnectRetryWait))
		}
	}

	close := make(chan bool, 1)
	link := Link{src: connIo, close: close}
	parse.LogConnDebug("send", pkt, "Forwarding")
	go link.pipe(pkt, conn, connIo, connFwd)
	go link.pipe(pkt, connFwd, connFwd, connIo)
	<-close
	parse.LogConnDebug("send", pkt, "Closed")
}
