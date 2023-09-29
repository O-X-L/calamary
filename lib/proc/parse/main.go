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

func Parse(l4Proto string, conn net.Conn, connIo io.Reader) (pkt ParsedPacket) {
	// get packet L5-header
	conn.SetReadDeadline(time.Now().Add(u.Timeout(cnf.C.Service.Timeout.Process)))
	var hdr [cnf.BYTES_HDR_L5]byte
	n, err := io.ReadFull(connIo, hdr[:])
	if err != nil {
		log.Warn("parse", fmt.Sprintf("Error parsing L5Header: %v", err))
	}
	connIo = io.MultiReader(bytes.NewReader(hdr[:n]), connIo) // write header back to stream

	// main L4 split
	if l4Proto == "udp" {
		pkt = parseUdp(conn, connIo, hdr)
	} else {
		pkt = parseTcp(conn, connIo, hdr)
	}

	conn.SetReadDeadline(time.Time{})
	l5ProtoStr := meta.RevProto(pkt.L5.Proto)
	tlsVersionStr := meta.RevTlsVersion(pkt.L5.TlsVersion)
	if cnf.Metrics() {
		metrics.RuleReqL3Proto.WithLabelValues(meta.RevProto(pkt.L3.Proto)).Inc()
		metrics.RuleReqL5Proto.WithLabelValues(l5ProtoStr).Inc()
		metrics.RuleReqTlsVersion.WithLabelValues(tlsVersionStr).Inc()
	}
	log.ConnDebug(
		"parse", PktSrc(pkt), PktDest(pkt),
		fmt.Sprintf("Packet L5Proto: %v | TLS: v%v", l5ProtoStr, tlsVersionStr),
	)
	return
}

func parseTcp(conn net.Conn, connIo io.Reader, hdr [cnf.BYTES_HDR_L5]byte) ParsedPacket {
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
		log.ConnErrorS("parse", conn.RemoteAddr().String(), "?", "Failed to resolve TCP source-address")
	}
	pkt.L3.SrcIP = tcpSrcAddr.IP
	pkt.L4.SrcPort = uint16(tcpSrcAddr.Port)
	pkt.L3.Proto = getL3Proto(pkt.L3.SrcIP)

	log.ConnDebug("parse", PktSrc(pkt), "?", "Parsing TCP connection")

	// destination address
	var dstIpPort net.Addr
	if !cnf.C.Service.Listen.TProxy {
		dstIpPort, err = getTcpOriginalDstAddr(conn)
		if err != nil {
			log.ConnErrorS("parse", PktSrc(pkt), "?", "Failed to get original destination IP")
		}
	} else {
		dstIpPort = conn.LocalAddr()
	}

	pkt.L3.Proto = parseIpProto(dstIpPort)
	tcpDestAddr, err := net.ResolveTCPAddr("tcp", dstIpPort.String())
	if err != nil {
		log.ConnErrorS("parse", PktSrc(pkt), "?", "Failed to resolve TCP destination-address")
	}
	pkt.L3.DestIP = tcpDestAddr.IP
	pkt.L4.DestPort = uint16(tcpDestAddr.Port)

	// additional
	log.ConnDebug("parse", PktSrc(pkt), PktDest(pkt), "Processing TCP")
	pkt.L5.Encrypted, pkt.L5.TlsVersion, pkt.L5.TlsSni = parseTls(pkt, conn, connIo, hdr)

	if pkt.L5.Encrypted == meta.OptBoolTrue {
		pkt.L5.Proto = meta.ProtoL5Tls
	} else {
		parseHttp(pkt, hdr)
	}

	return pkt
}

func parseUdp(conn net.Conn, connIo io.Reader, hdr [cnf.BYTES_HDR_L5]byte) ParsedPacket {
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
			log.ConnErrorS("parse", PktSrcIP(pkt), "?", "Failed to get original destination IP")
		}
		pkt.L3.SrcIP = raddr.String()
		pkt.L3.DestIP = dstAddr.String()
		pkt.L3.Proto = getL3Proto(pkt.L3.SrcIP)

	*/
	/*
		udpAddr, err := net.ResolveUDPAddr("tcp", dstIpPort.String())
		if err != nil {
			log.ConnErrorS("parse", PktSrcIP(pkt), "?", "Failed to resolve UDP address")
		}
		pkt.L3.DestIP = udpAddr.IP
		pkt.L4.DestPort = udpAddr.Port
	*/

	log.ConnDebug("parse", PktSrc(pkt), PktDest(pkt), "Processing UDP")

	return pkt
}

func getL3Proto(ip net.IP) meta.Proto {
	if ip.To4() != nil {
		return meta.ProtoL3IP4
	} else {
		return meta.ProtoL3IP6
	}
}
