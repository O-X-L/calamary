package cnf_file

import (
	"testing"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
)

func TestCleanRaw(t *testing.T) {
	if cleanRaw(" te st ") != "test" {
		t.Error("Clean raw #1 failed")
	}
	if cleanRaw("!test ") != "test" {
		t.Error("Clean raw #2 failed")
	}
}

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

func TestMatchNet(t *testing.T) {
	nets1 := matchNet("192.168.1.24")
	if nets1.String() != "192.168.1.24/32" {
		t.Error("Match IPs #1")
	}

	nets2 := matchNet("192.168.2.0/24")
	if nets2.String() != "192.168.2.0/24" {
		t.Error("Match IPs #2")
	}

	nets3 := matchNet("2001:db8::1")
	if nets3.String() != "2001:db8::1/128" {
		t.Error("Match IPs #3")
	}

	nets4 := matchNet("2001:db8::/80")
	if nets4.String() != "2001:db8::/80" {
		t.Error("Match IPs #4")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match IPs not failing on invalid value")
		}
	}()
	matchNet("random")
}

func TestMatchPort(t *testing.T) {
	if matchPort("80") != 80 || matchPort("!80") != 80 || matchPort("! 80") != 80 {
		t.Error("Match port #1")
	}
	if matchPort("50000") != 50000 || matchPort("!50000") != 50000 {
		t.Error("Match port #2")
	}
	if matchPort("65000") != 65000 || matchPort("!65000") != 65000 {
		t.Error("Match port #3")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match port not failing on invalid value #1")
		}
	}()
	matchPort("66000")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match port not failing on invalid value #2")
		}
	}()
	matchPort("random")
}

func TestParseRules(t *testing.T) {
	ParseRules([]cnf.RuleRaw{
		cnf.RuleRaw{
			Match: cnf.MatchRaw{
				ProtoL3:  []string{"ip4"},
				ProtoL4:  []string{"tcp"},
				DestPort: []string{"80"},
				Domains:  []string{"superstes.eu"},
				DestNet:  []string{"!192.168.0.0/16", "!172.16.0.0/12", "!10.0.0.0/8"},
			},
			Action: "allow",
		},
	})
	ParseRules([]cnf.RuleRaw{
		cnf.RuleRaw{
			Match: cnf.MatchRaw{
				ProtoL3:  []string{"ip4", "ipv6"},
				ProtoL4:  []string{"tcp", "udp"},
				DestPort: []string{"80", "443", "8443"},
				Domains:  []string{"superstes.eu", "oxl.at"},
				DestNet:  []string{"8.8.8.8", "1.1.1.1"},
			},
			Action: "deny",
		},
	})
}
