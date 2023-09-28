package fwd

import (
	"bytes"
	"io"
	"net"

	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/filter"
	"github.com/superstes/calamary/proc/parse"
	"github.com/superstes/calamary/send"
	"github.com/superstes/calamary/u"
)

func Forward(l4Proto string, conn net.Conn) {
	defer conn.Close()
	var connIo io.ReadWriter = conn
	connIoBuf := new(bytes.Buffer)
	connIoTee := io.TeeReader(connIo, connIoBuf)
	pkt := parse.Parse(l4Proto, conn, connIoTee)
	if !filter.Filter(pkt) {
		log.ConnInfo(
			"forward", parse.PktSrc(pkt), parse.PktDest(pkt),
			"Denied",
		)
		conn.Close()
		return
	}
	// write header back to stream
	connIo = u.NewReadWriter(io.MultiReader(bytes.NewReader(connIoBuf.Bytes()), connIo), connIo)

	log.ConnInfo(
		"forward", parse.PktSrc(pkt), parse.PktDest(pkt),
		"Accept",
	)
	send.Forward(pkt, conn, connIo)
}
