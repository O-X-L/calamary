package parse

import (
	"fmt"
	"net"

	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/u"
)

func PktSrc(pkt ParsedPacket) string {
	return pkt.L3.SrcIP.String()
}

func PktDest(pkt ParsedPacket) string {
	return fmt.Sprintf("%s:%v", pkt.L3.DestIP.String(), pkt.L4.DestPort)
}

func PktDestIP(pkt ParsedPacket) string {
	return pkt.L3.DestIP.String()
}

func parseIpProto(addr net.Addr) meta.Proto {
	if u.IsIPv4(addr.String()) {
		return meta.ProtoL3IP4
	}
	return meta.ProtoL3IP6
}
