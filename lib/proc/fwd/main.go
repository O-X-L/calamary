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
			"forward", parse.PkgSrc(pkt), parse.PkgDest(pkt),
			"Accept",
		)

	} else {
		log.ConnInfo(
			"forward", parse.PkgSrc(pkt), parse.PkgDest(pkt),
			"Denied",
		)
		conn.Close()
	}
}
