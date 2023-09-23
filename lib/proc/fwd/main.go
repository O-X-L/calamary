package fwd

import (
	"net"

	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/filter"
	parse "github.com/superstes/calamary/proc/parse"
)

func Forward(l4Proto string, conn net.Conn) {
	pkg := parse.Parse(l4Proto, conn)
	if filter.Filter(pkg) {
		// dialer := nw.NewDialerDirect()
		log.ConnInfo(
			"forward", parse.PkgSrc(pkg), parse.PkgDest(pkg),
			"Accept",
		)

	} else {
		log.ConnInfo(
			"forward", parse.PkgSrc(pkg), parse.PkgDest(pkg),
			"Denied",
		)
		conn.Close()
	}
}
