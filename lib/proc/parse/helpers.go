package parse

import (
	"fmt"
	"net"

	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/u"
)

func PkgSrc(pkg ParsedPackage) string {
	return pkg.L3.SrcIP.String()
}

func PkgDest(pkg ParsedPackage) string {
	return fmt.Sprintf("%s:%v", pkg.L3.DestIP.String(), pkg.L4.DestPort)
}

func PkgDestIP(pkg ParsedPackage) string {
	return pkg.L3.DestIP.String()
}

func parseIpProto(addr net.Addr) meta.Proto {
	if u.IsIPv4(addr.String()) {
		return meta.ProtoL3IP4
	}
	return meta.ProtoL3IP6
}
