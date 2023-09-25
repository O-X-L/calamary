package parse

import (
	"net"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
)

func Parse(l4Proto string, conn net.Conn) ParsedPackage {
	if l4Proto == "udp" {
		return parseUdp(conn)
	} else {
		return parseTcp(conn)
	}
}

func parseTcp(conn net.Conn) ParsedPackage {
	pkg := ParsedPackage{
		L3: &ParsedL3Package{
			L4Proto: meta.ProtoL4Tcp,
		},
		L4: &ParsedL4Package{
			L5Proto: meta.ProtoNone,
		},
		L4Tcp:  &ParsedTcpPackage{},
		L5Http: &ParsedHttpPackage{},
	}
	// source address
	tcpSrcAddr, err := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
	if err != nil {
		log.ConnErrorS("proc-parse", conn.RemoteAddr().String(), "?", "Failed to resolve TCP source-address")
	}
	pkg.L3.SrcIP = tcpSrcAddr.IP
	pkg.L4.SrcPort = uint16(tcpSrcAddr.Port)
	pkg.L3.Proto = getL3Proto(pkg.L3.SrcIP)

	log.ConnDebug("proc-parse", PkgSrc(pkg), "?", "Parsing TCP connection")

	// destination address
	var dstIpPort net.Addr
	if !cnf.C.Service.Listen.Transparent {
		dstIpPort, err = getTcpOriginalDstAddr(conn)
		if err != nil {
			log.ConnErrorS("proc-parse", PkgSrc(pkg), "?", "Failed to get original destination IP")
		}
	} else {
		dstIpPort = conn.LocalAddr()
	}

	pkg.L3.Proto = parseIpProto(dstIpPort)
	tcpDestAddr, err := net.ResolveTCPAddr("tcp", dstIpPort.String())
	if err != nil {
		log.ConnErrorS("proc-parse", PkgSrc(pkg), "?", "Failed to resolve TCP destination-address")
	}
	pkg.L3.DestIP = tcpDestAddr.IP
	pkg.L4.DestPort = uint16(tcpDestAddr.Port)

	// additional
	log.ConnDebug("proc-parse", PkgSrc(pkg), PkgDest(pkg), "Processing TCP")

	return pkg
}

func parseUdp(conn net.Conn) ParsedPackage {
	pkg := ParsedPackage{
		L3: &ParsedL3Package{
			L4Proto: meta.ProtoL4Udp,
		},
		L4: &ParsedL4Package{
			L5Proto: meta.ProtoNone,
		},
		L4Udp: &ParsedUdpPackage{},
	}
	/*
		b := u.GetBufferPool(cnf.UDP_BUFFER_SIZE)

		_, raddr, dstAddr, err := getUdpOriginalDstAddr(*conn, *b)
		if err != nil {
			log.ConnErrorS("proc-parse", PkgSrcIP(pkg), "?", "Failed to get original destination IP")
		}
		pkg.L3.SrcIP = raddr.String()
		pkg.L3.DestIP = dstAddr.String()
		pkg.L3.Proto = getL3Proto(pkg.L3.SrcIP)

	*/
	/*
		udpAddr, err := net.ResolveUDPAddr("tcp", dstIpPort.String())
		if err != nil {
			log.ConnErrorS("proc-parse", PkgSrcIP(pkg), "?", "Failed to resolve UDP address")
		}
		pkg.L3.DestIP = udpAddr.IP
		pkg.L4.DestPort = udpAddr.Port
	*/

	log.ConnDebug("proc-parse", PkgSrc(pkg), PkgDest(pkg), "Processing UDP")

	return pkg
}

func getL3Proto(ip net.IP) meta.Proto {
	if ip.To4() != nil {
		return meta.ProtoL3IP4
	} else {
		return meta.ProtoL3IP6
	}
}
