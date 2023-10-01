package u

import (
	"net"
	"testing"

	"github.com/superstes/calamary/cnf"
)

func TestDnsLookup(t *testing.T) {
	// NOTE: tests will fail if the IPs change.. should not be common
	cnf.C = &cnf.Config{
		Service: cnf.ServiceConfig{
			DnsNameservers: []string{"1.1.1.1"},
		},
	}
	resp := DnsLookup("one.one.one.one")

	if !isInIpList("1.1.1.1", resp) || !isInIpList("1.0.0.1", resp) ||
		!isInIpList("2606:4700:4700::1001", resp) || !isInIpList("2606:4700:4700::1111", resp) {
		t.Errorf("DNS Query #1 has unexpected result: %v", resp)
	}

	resp4, resp6 := DnsLookup46("one.one.one.one")

	if !isInIpList("1.1.1.1", resp4) || !isInIpList("1.0.0.1", resp4) ||
		!isInIpList("2606:4700:4700::1001", resp6) || !isInIpList("2606:4700:4700::1111", resp6) ||
		isInIpList("1.1.1.1", resp6) || isInIpList("2606:4700:4700::1001", resp4) {
		t.Errorf("DNS Query #2 has unexpected result: ip4 %v, ip6 %v", resp4, resp6)
	}
}

func isInIpList(value string, list []net.IP) bool {
	for i := range list {
		if list[i].String() == value {
			return true
		}
	}
	return false
}
