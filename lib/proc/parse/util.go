package parse

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/u"
)

func PktSrc(pkt ParsedPacket) string {
	return pkt.L3.SrcIP.String()
}

func PktDest(pkt ParsedPacket) string {
	if pkt.L3.DestIP == nil {
		return "?"
	}
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

func SplitHttpHost(host string, encypted meta.OptBool) (dns string, port uint16) {
	if strings.Contains(host, ":") {
		parts := strings.SplitN(host, ":", 2)
		port, err := strconv.Atoi(parts[1])

		if err == nil {
			return parts[0], uint16(port)
		}
	}

	if encypted == meta.OptBoolTrue {
		return host, 443
	}
	return host, 80
}
