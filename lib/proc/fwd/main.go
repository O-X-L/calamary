package fwd

import (
	"net"

	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/filter"
	"github.com/superstes/calamary/proc/parse"
)

func Forward(l4Proto string, conn net.Conn) {
	pkt := parse.Parse(l4Proto, conn)
	if filter.Filter(pkt) {
		// dialer := nw.NewDialerDirect()
		log.ConnInfo(
			"forward", parse.PktSrc(pkt), parse.PktDest(pkt),
			"Accept",
		)

	} else {
		log.ConnInfo(
			"forward", parse.PktSrc(pkt), parse.PktDest(pkt),
			"Denied",
		)
		conn.Close()
	}
}
