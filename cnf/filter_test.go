package cnf

import (
	"testing"

	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/u"
)

func TestMatchFilterAction(t *testing.T) {
	if filterAction("allow") != meta.ActionAccept || filterAction("accept") != meta.ActionAccept {
		t.Error("Filter action accept")
	}
	if filterAction("deny") != meta.ActionDeny || filterAction("random") != meta.ActionDeny {
		t.Error("Filter action deny")
	}
}

func TestMatchProtoL3(t *testing.T) {
	if matchProtoL3("ip4") != meta.ProtoL3IP4 || matchProtoL3("ipv4") != meta.ProtoL3IP4 || matchProtoL3("IPv4") != meta.ProtoL3IP4 {
		t.Error("Match proto L3 ip4")
	}
	if matchProtoL3("ip6") != meta.ProtoL3IP6 || matchProtoL3("ipv6") != meta.ProtoL3IP6 || matchProtoL3("IPv6") != meta.ProtoL3IP6 {
		t.Error("Match proto L3 ip6")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match proto L3 not failing on invalid value")
		}
	}()
	matchProtoL3("random")
}

func TestMatchProtoL4(t *testing.T) {
	if matchProtoL4("tcp") != meta.ProtoL4Tcp || matchProtoL4("TCP") != meta.ProtoL4Tcp {
		t.Error("Match proto L4 tcp")
	}
	if matchProtoL4("udp") != meta.ProtoL4Udp || matchProtoL4("UDP") != meta.ProtoL4Udp {
		t.Error("Match proto L4 udp")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match proto L4 not failing on invalid value")
		}
	}()
	matchProtoL4("random")
}

func TestMatchProtoL5(t *testing.T) {
	if matchProtoL5("tls") != meta.ProtoL5Tls {
		t.Error("Match proto L5 tls")
	}
	if matchProtoL5("http") != meta.ProtoL5Http {
		t.Error("Match proto L5 http")
	}
	/*
		if matchProtoL5("dns") != meta.ProtoL5Dns {
			t.Error("Match proto L5 dns")
		}
		if matchProtoL5("ntp") != meta.ProtoL5Ntp {
			t.Error("Match proto L5 ntp")
		}
	*/
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match proto L5 not failing on invalid value")
		}
	}()
	matchProtoL5("random")
}

func TestMatchIPs(t *testing.T) {
	nets1 := matchIPs([]string{"192.168.0.1", "192.168.10.0/24"})
	if !u.AllStrInList(
		[]string{nets1[0].String(), nets1[1].String()},
		[]string{"192.168.0.1/32", "192.168.10.0/24"},
	) {
		t.Error("Match IPs #1")
	}

	nets2 := matchIPs([]string{"192.168.20.0/29", "2001:db8::1", "2001:db8::/80"})
	if !u.AllStrInList(
		[]string{nets2[0].String(), nets2[1].String(), nets2[2].String()},
		[]string{"192.168.20.0/29", "2001:db8::1/128", "2001:db8::/80"},
	) {
		t.Errorf("Match IPs #2 %v %v %v", nets2[0].String(), nets2[1].String(), nets2[2].String())
	}

	nets3 := matchIPs([]string{"192.168.40.0/24"})
	if !u.AllStrInList(
		[]string{nets3[0].String()},
		[]string{"192.168.40.0/24"},
	) {
		t.Error("Match IPs #3")
	}

	nets4 := matchIPs([]string{"192.168.50.0/24"})
	if !u.AllStrInList(
		[]string{nets4[0].String()},
		[]string{"192.168.50.0/24"},
	) {
		t.Error("Match IPs #4")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match IPs not failing on invalid value")
		}
	}()
	matchIPs([]string{"random"})
}
