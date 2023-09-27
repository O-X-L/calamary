package parse

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
	tls_dissector "github.com/superstes/calamary/proc/parse/tls"
	"github.com/superstes/calamary/u"
)

func Parse(l4Proto string, conn net.Conn, connIo io.ReadWriter) ParsedPacket {
	// get packet L5-header
	conn.SetReadDeadline(time.Now().Add(
		time.Duration(cnf.C.Service.Timeout.Intercept * int(time.Millisecond)),
	))
	var hdr [cnf.L5HDRLEN]byte
	n, err := io.ReadFull(connIo, hdr[:])
	conn.SetReadDeadline(time.Time{})
	if err == nil {
		connIo = u.NewReadWriter(io.MultiReader(bytes.NewReader(hdr[:n]), connIo), connIo)
	} else {
		log.Warn("parse", fmt.Sprintf("Error parsing L5Header: %v", err))
	}

	// main L4 split
	var pkt ParsedPacket
	if l4Proto == "udp" {
		pkt = parseUdp(conn, connIo, hdr)
	} else {
		pkt = parseTcp(conn, connIo, hdr)
	}

	log.ConnDebug(
		"parse", PktSrc(pkt), PktDest(pkt),
		fmt.Sprintf("Packet L5Proto: %v | TLS: %v", meta.RevProto(pkt.L5.Proto), meta.RevTlsVersion(pkt.L5.TlsVersion)),
	)
	return pkt
}

func parseTcp(conn net.Conn, connIo io.ReadWriter, hdr [cnf.L5HDRLEN]byte) ParsedPacket {
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
	if !cnf.C.Service.Listen.Transparent {
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
	} else if hdrL5Http(hdr) {
		pkt.L5.Proto = meta.ProtoL5Http
		pkt.L5Http = &ParsedHttp{}
	}

	return pkt
}

func parseUdp(conn net.Conn, connIo io.ReadWriter, hdr [cnf.L5HDRLEN]byte) ParsedPacket {
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

func parseTls(pkt ParsedPacket, conn net.Conn, connIo io.ReadWriter, hdr [cnf.L5HDRLEN]byte) (isTls meta.OptBool, tlsVersion uint16, sni string) {
	isTlsRaw := hdr[0] == tls_dissector.Handshake
	tlsVersion = binary.BigEndian.Uint16(hdr[1:3])
	if isTlsRaw {
		isTls = meta.OptBoolTrue
	} else {
		isTls = meta.OptBoolFalse
	}

	if isTlsRaw {
		buf := new(bytes.Buffer)
		r := io.TeeReader(connIo, buf)
		record, err := tls_dissector.ReadRecord(r)
		if err != nil {
			return
		}

		clientHello := tls_dissector.ClientHelloMsg{}
		if err = clientHello.Decode(record.Opaque); err != nil {
			return
		}

		for _, ext := range clientHello.Extensions {
			if ext.Type() == tls_dissector.ExtServerName {
				snExtension := ext.(*tls_dissector.ServerNameExtension)
				sni = snExtension.Name
				break
			}
		}
	}
	log.ConnDebug("parse", PktSrc(pkt), PktDest(pkt), fmt.Sprintf(
		"TLS information: IsTls=%v, TlsVersion=%v, TlsSni=%s", isTlsRaw, tlsVersion, sni,
	))
	return
}

func hdrL5Http(hdr [cnf.L5HDRLEN]byte) bool {
	s := string(hdr[:])
	return strings.HasPrefix(http.MethodGet, s[:3]) ||
		strings.HasPrefix(http.MethodPost, s[:4]) ||
		strings.HasPrefix(http.MethodPut, s[:3]) ||
		strings.HasPrefix(http.MethodDelete, s) ||
		strings.HasPrefix(http.MethodOptions, s) ||
		strings.HasPrefix(http.MethodPatch, s) ||
		strings.HasPrefix(http.MethodHead, s[:4]) ||
		strings.HasPrefix(http.MethodConnect, s) ||
		strings.HasPrefix(http.MethodTrace, s)
}
