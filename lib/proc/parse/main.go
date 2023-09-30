package parse

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/metrics"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/u"
)

func Parse(srvCnf cnf.ServiceListener, l4Proto meta.Proto, conn net.Conn, connIo io.Reader) (pkt ParsedPacket) {
	// get packet L5-header
	conn.SetReadDeadline(time.Now().Add(u.Timeout(cnf.C.Service.Timeout.Process)))
	var hdr [cnf.BYTES_HDR_L5]byte
	n, err := io.ReadFull(connIo, hdr[:])
	if err != nil {
		log.Warn("parse", fmt.Sprintf("Error parsing L5Header: %v", err))
	}
	connIo = io.MultiReader(bytes.NewReader(hdr[:n]), connIo) // write header back to stream

	// main L4 split
	if l4Proto == meta.ProtoL4Udp {
		pkt = parseUdp(srvCnf, conn, connIo, hdr)
	} else {
		pkt = parseTcp(srvCnf, conn, connIo, hdr)
	}

	conn.SetReadDeadline(time.Time{})
	l5ProtoStr := meta.RevProto(pkt.L5.Proto)
	tlsVersionStr := meta.RevTlsVersion(pkt.L5.TlsVersion)
	if cnf.Metrics() {
		metrics.ReqL3Proto.WithLabelValues(meta.RevProto(pkt.L3.Proto)).Inc()
		metrics.ReqL5Proto.WithLabelValues(l5ProtoStr).Inc()
		metrics.ReqTlsVersion.WithLabelValues(tlsVersionStr).Inc()
	}
	LogConnDebug("parse", pkt, fmt.Sprintf("Packet L5Proto: %v | TLS: %v", l5ProtoStr, tlsVersionStr))
	return
}

func parseTcp(srvCnf cnf.ServiceListener, conn net.Conn, connIo io.Reader, hdr [cnf.BYTES_HDR_L5]byte) ParsedPacket {
	pkt := ParsedPacket{
		L3: &ParsedL3{},
		L4: &ParsedL4{
			Proto: meta.ProtoL4Tcp,
		},
		L5: &ParsedL5{
			Proto:     meta.ProtoNone,
			Encrypted: meta.OptBoolNone,
		},
		L4Tcp: &ParsedTcp{},
	}
	// source address
	tcpSrcAddr, err := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
	if err != nil {
		log.ConnError("parse", conn.RemoteAddr().String(), "?", "Failed to resolve TCP source-address")
	}
	pkt.L3.SrcIP = tcpSrcAddr.IP
	pkt.L4.SrcPort = uint16(tcpSrcAddr.Port)
	pkt.L3.Proto = getL3Proto(pkt.L3.SrcIP)

	LogConnDebug("parse", pkt, "Parsing TCP connection")

	// destination address
	var dstIpPort net.Addr
	if srvCnf.TProxy {
		dstIpPort = conn.LocalAddr()

	} else {
		dstIpPort, err = getTcpOriginalDstAddr(conn)
		if err != nil {
			LogConnError("parse", pkt, "Failed to get original destination IP")
		}
	}

	pkt.L3.Proto = parseIpProto(dstIpPort)
	tcpDestAddr, err := net.ResolveTCPAddr("tcp", dstIpPort.String())
	if err != nil {
		LogConnError("parse", pkt, "Failed to resolve TCP destination-address")
	}
	pkt.L3.DestIP = tcpDestAddr.IP
	pkt.L4.DestPort = uint16(tcpDestAddr.Port)

	// additional
	LogConnDebug("parse", pkt, "Processing TCP")
	pkt.L5.Encrypted, pkt.L5.TlsVersion, pkt.L5.TlsSni = parseTls(pkt, conn, connIo, hdr)

	if pkt.L5.Encrypted == meta.OptBoolTrue {
		pkt.L5.Proto = meta.ProtoL5Tls
	} else {
		parseHttp(pkt, hdr)
	}

	return pkt
}

func parseUdp(srvCnf cnf.ServiceListener, conn net.Conn, connIo io.Reader, hdr [cnf.BYTES_HDR_L5]byte) ParsedPacket {
	pkt := ParsedPacket{
		L3: &ParsedL3{},
		L4: &ParsedL4{
			Proto: meta.ProtoL4Udp,
		},
		L5: &ParsedL5{
			Proto:     meta.ProtoNone,
			Encrypted: meta.OptBoolNone,
		},
		L4Udp: &ParsedUdp{},
	}
	/*
		b := u.GetBufferPool(cnf.UDP_BUFFER_SIZE)

		_, raddr, dstAddr, err := getUdpOriginalDstAddr(*conn, *b)
		if err != nil {
			LogConnError("parse", pkt, "Failed to get original destination IP")
		}
		pkt.L3.SrcIP = raddr.String()
		pkt.L3.DestIP = dstAddr.String()
		pkt.L3.Proto = getL3Proto(pkt.L3.SrcIP)

	*/
	/*
		udpAddr, err := net.ResolveUDPAddr("tcp", dstIpPort.String())
		if err != nil {
			LogConnError("parse", pkt, "Failed to resolve UDP address")
		}
		pkt.L3.DestIP = udpAddr.IP
		pkt.L4.DestPort = udpAddr.Port
	*/

	LogConnDebug("parse", pkt, "Processing UDP")

	return pkt
}

func getL3Proto(ip net.IP) meta.Proto {
	if ip.To4() != nil {
		return meta.ProtoL3IP4
	} else {
		return meta.ProtoL3IP6
	}
}
