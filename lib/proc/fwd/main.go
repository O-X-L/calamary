package fwd

import (
	"bytes"
	"io"
	"net"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/metrics"
	"github.com/superstes/calamary/proc/filter"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
	"github.com/superstes/calamary/send"
	"github.com/superstes/calamary/u"
)

func Forward(srvCnf cnf.ServiceListener, l4Proto meta.Proto, conn net.Conn) {
	defer conn.Close()

	if cnf.Metrics() {
		metrics.ReqTcp.Inc()
		metrics.CurrentConn.Inc()
		defer metrics.CurrentConn.Dec()
	}

	var connIo io.ReadWriter = conn
	connIoBuf := new(bytes.Buffer)
	connIoTee := io.TeeReader(connIo, connIoBuf)
	pkt := parse.Parse(srvCnf, l4Proto, conn, connIoTee)
	if !filter.Filter(pkt) {
		parse.LogConnInfo("forward", pkt, "Denied")
		conn.Close()
		return
	}
	// write read bytes back to stream so we can forward them
	connIo = u.NewReadWriter(io.MultiReader(bytes.NewReader(connIoBuf.Bytes()), connIo), connIo)

	parse.LogConnInfo("forward", pkt, "Accept")
	send.Forward(pkt, conn, connIo)
}
