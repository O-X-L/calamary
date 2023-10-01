package fwd

import (
	"fmt"
	"net"
	"time"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/parse"
	"github.com/superstes/calamary/u"
)

func FirstReachableTarget(pkt parse.ParsedPacket, dest4 []net.IP, dest6 []net.IP, l4Proto string, port uint16) net.IP {
	retries := 1
	for {
		for _, ip6 := range dest6 {
			if TargetIsReachable(l4Proto, ip6, port) {
				return ip6
			}
		}
		for _, ip4 := range dest4 {
			if TargetIsReachable(l4Proto, ip4, port) {
				return ip4
			}
		}
		if retries >= int(cnf.C.Service.Output.Retries) {
			break
		}
		parse.LogConnDebug(
			"send", pkt, fmt.Sprintf("Connection probe to all available IPs failed (retry %v/%v)",
				retries, cnf.C.Service.Output.Retries),
		)
		retries++
		time.Sleep(u.Timeout(cnf.DefaultConnectRetryWait))
	}
	return nil
}

func TargetIsReachable(l4Proto string, target net.IP, port uint16) bool {
	targetAddr := fmt.Sprintf("%s:%v", target.String(), port)
	dummyConn, err := net.DialTimeout(l4Proto, targetAddr, u.Timeout(cnf.C.Service.Timeout.Probe))
	if dummyConn != nil {
		dummyConn.Close()
	}
	return err == nil
}
