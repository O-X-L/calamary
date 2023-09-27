package parse

import (
	"net"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
)

func Parse(l4Proto string, conn net.Conn) ParsedPacket {
	if l4Proto == "udp" {
		return parseUdp(conn)
	} else {
		return parseTcp(conn)
	}
}

func parseTcp(conn net.Conn) ParsedPacket {
	pkt := ParsedPacket{
		L3: &ParsedL3{},
		L4: &ParsedL4{
			Proto: meta.ProtoL4Tcp,
		},
		L5: &ParsedL5{
			Proto:     meta.ProtoNone,
			Encrypted: meta.OptBoolNone,
		},
		L4Tcp:  &ParsedTcp{},
		L5Http: &ParsedHttp{},
	}
	// source address
	tcpSrcAddr, err := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
	if err != nil {
		log.ConnErrorS("proc-parse", conn.RemoteAddr().String(), "?", "Failed to resolve TCP source-address")
	}
	pkt.L3.SrcIP = tcpSrcAddr.IP
	pkt.L4.SrcPort = uint16(tcpSrcAddr.Port)
	pkt.L3.Proto = getL3Proto(pkt.L3.SrcIP)

	log.ConnDebug("proc-parse", PktSrc(pkt), "?", "Parsing TCP connection")

	// destination address
	var dstIpPort net.Addr
	if !cnf.C.Service.Listen.Transparent {
		dstIpPort, err = getTcpOriginalDstAddr(conn)
		if err != nil {
			log.ConnErrorS("proc-parse", PktSrc(pkt), "?", "Failed to get original destination IP")
		}
	} else {
		dstIpPort = conn.LocalAddr()
	}

	pkt.L3.Proto = parseIpProto(dstIpPort)
	tcpDestAddr, err := net.ResolveTCPAddr("tcp", dstIpPort.String())
	if err != nil {
		log.ConnErrorS("proc-parse", PktSrc(pkt), "?", "Failed to resolve TCP destination-address")
	}
	pkt.L3.DestIP = tcpDestAddr.IP
	pkt.L4.DestPort = uint16(tcpDestAddr.Port)

	// additional
	log.ConnDebug("proc-parse", PktSrc(pkt), PktDest(pkt), "Processing TCP")

	return pkt
}

func parseUdp(conn net.Conn) ParsedPacket {
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
			log.ConnErrorS("proc-parse", PktSrcIP(pkt), "?", "Failed to get original destination IP")
		}
		pkt.L3.SrcIP = raddr.String()
		pkt.L3.DestIP = dstAddr.String()
		pkt.L3.Proto = getL3Proto(pkt.L3.SrcIP)

	*/
	/*
		udpAddr, err := net.ResolveUDPAddr("tcp", dstIpPort.String())
		if err != nil {
			log.ConnErrorS("proc-parse", PktSrcIP(pkt), "?", "Failed to resolve UDP address")
		}
		pkt.L3.DestIP = udpAddr.IP
		pkt.L4.DestPort = udpAddr.Port
	*/

	log.ConnDebug("proc-parse", PktSrc(pkt), PktDest(pkt), "Processing UDP")

	return pkt
}

func getL3Proto(ip net.IP) meta.Proto {
	if ip.To4() != nil {
		return meta.ProtoL3IP4
	} else {
		return meta.ProtoL3IP6
	}
}
