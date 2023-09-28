package filter

import (
	"net"
	"testing"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

func testNet(netIn string) *net.IPNet {
	_, netOut, _ := net.ParseCIDR(netIn)
	return netOut
}

func TestAnyPortMatch(t *testing.T) {
	if anyPortMatch([]uint16{80, 443}, uint16(80)) == meta.MatchNegative {
		t.Error("AnyPort #1")
	}
	if anyPortMatch([]uint16{389, 34938, 1022}, uint16(34938)) == meta.MatchNegative {
		t.Error("AnyPort #2")
	}
	if anyPortMatch([]uint16{389, 34938, 1022}, uint16(20930)) == meta.MatchPositive {
		t.Error("AnyPort #3")
	}
	if anyPortMatch([]uint16{65000, 443}, uint16(1000)) == meta.MatchPositive {
		t.Error("AnyPort #3")
	}
}

func TestAnyNetMatch(t *testing.T) {
	net1 := testNet("192.168.0.0/16")
	ip := net.ParseIP("192.168.0.1")
	if anyNetMatch([]*net.IPNet{net1}, ip) == meta.MatchNegative {
		t.Error("AnyNet #1")
	}
	ip = net.ParseIP("192.167.0.1")
	if anyNetMatch([]*net.IPNet{net1}, ip) == meta.MatchPositive {
		t.Error("AnyNet #2")
	}
	net2 := testNet("10.0.0.0/8")
	ip = net.ParseIP("10.255.0.1")
	if anyNetMatch([]*net.IPNet{net1, net2}, ip) == meta.MatchNegative {
		t.Error("AnyNet #3")
	}
	ip = net.ParseIP("172.16.0.1")
	if anyNetMatch([]*net.IPNet{net1, net2}, ip) == meta.MatchPositive {
		t.Error("AnyNet #4")
	}
	net3 := testNet("2001:db8::/80")
	ip = net.ParseIP("192.168.251.48")
	if anyNetMatch([]*net.IPNet{net1, net2, net3}, ip) == meta.MatchNegative {
		t.Error("AnyNet #5")
	}
	ip = net.ParseIP("2001:db8::1:9")
	if anyNetMatch([]*net.IPNet{net1, net2, net3}, ip) == meta.MatchNegative {
		t.Error("AnyNet #6")
	}
	ip = net.ParseIP("2002:db8::1:9")
	if anyNetMatch([]*net.IPNet{net1, net2, net3}, ip) == meta.MatchPositive {
		t.Error("AnyNet #7")
	}
}

func TestAnyProtoMatch(t *testing.T) {
	if anyProtoMatch([]meta.Proto{meta.ProtoL4Tcp, meta.ProtoL4Udp}, meta.ProtoL4Udp) == meta.MatchNegative {
		t.Error("AnyProto #1")
	}
	if anyProtoMatch([]meta.Proto{meta.ProtoL5Dns, meta.ProtoL5Tls}, meta.ProtoL5Http) == meta.MatchPositive {
		t.Error("AnyProto #2")
	}
	if anyProtoMatch([]meta.Proto{meta.ProtoL5Tls}, meta.ProtoL5Http) == meta.MatchPositive {
		t.Error("AnyProto #3")
	}
}

func TestMatchSourcePort(t *testing.T) {
	cnf.C = &cnf.Config{}
	cnf.C.Service.Debug = false

	pkt := parse.ParsedPacket{
		L4: &parse.ParsedL4{
			SrcPort: uint16(51029),
		},
	}
	rule1 := cnf.Rule{}
	rule1.Match.SrcPort = []uint16{51029}
	if matchSourcePort(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Source-Ports #1")
	}
	rule1.Match.SrcPortN = []uint16{3000}
	if matchSourcePort(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Source-Ports #2")
	}
	rule1.Match.SrcPort = []uint16{4000}
	if matchSourcePort(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match Source-Ports #3")
	}
}

func TestMatchDestinationPort(t *testing.T) {
	cnf.C = &cnf.Config{}
	cnf.C.Service.Debug = false

	pkt := parse.ParsedPacket{
		L4: &parse.ParsedL4{
			DestPort: uint16(80),
		},
	}
	rule1 := cnf.Rule{}
	rule1.Match.DestPort = []uint16{80}
	if matchDestinationPort(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Destination-Ports #1")
	}
	rule1.Match.DestPortN = []uint16{3000}
	if matchDestinationPort(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Destination-Ports #2")
	}
	rule1.Match.DestPort = []uint16{4000}
	if matchDestinationPort(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match Destination-Ports #3")
	}
}

func TestMatchDestinationNetwork(t *testing.T) {
	cnf.C = &cnf.Config{}
	cnf.C.Service.Debug = false

	pkt := parse.ParsedPacket{
		L3: &parse.ParsedL3{
			DestIP: net.ParseIP("10.0.0.1"),
		},
	}
	rule1 := cnf.Rule{}
	rule1.Match.DestNetN = []*net.IPNet{testNet("172.16.0.0/24")}
	if matchDestinationNetwork(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Destination-Networks #1")
	}
	rule1.Match.DestNet = []*net.IPNet{testNet("192.168.0.0/24")}
	if matchDestinationNetwork(pkt, rule1, 1) != meta.MatchPositive {
		// note: edge-case behavior - might not be expected (because of DestNetN)
		t.Error("Match Destination-Networks #2")
	}
	rule1.Match.DestNetN = []*net.IPNet{}
	if matchDestinationNetwork(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match Destination-Networks #3")
	}
	rule1.Match.DestNet = []*net.IPNet{testNet("10.0.0.0/24")}
	if matchDestinationNetwork(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Destination-Networks #4")
	}
	rule1.Match.DestNetN = []*net.IPNet{testNet("192.168.0.0/29")}
	if matchDestinationNetwork(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Destination-Networks #5")
	}
	rule1.Match.DestNetN = []*net.IPNet{testNet("10.0.0.0/29")}
	if matchDestinationNetwork(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match Destination-Networks #6")
	}
}

func TestMatchSourceNetwork(t *testing.T) {
	cnf.C = &cnf.Config{}
	cnf.C.Service.Debug = false

	pkt := parse.ParsedPacket{
		L3: &parse.ParsedL3{
			SrcIP: net.ParseIP("10.0.0.1"),
		},
	}
	rule1 := cnf.Rule{}
	rule1.Match.SrcNetN = []*net.IPNet{testNet("172.16.0.0/24")}
	if matchSourceNetwork(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Source-Networks #1")
	}
	rule1.Match.SrcNet = []*net.IPNet{testNet("192.168.0.0/24")}
	if matchSourceNetwork(pkt, rule1, 1) != meta.MatchPositive {
		// note: edge-case behavior - might not be expected (because of DestNetN)
		t.Error("Match Source-Networks #2")
	}
	rule1.Match.SrcNetN = []*net.IPNet{}
	if matchSourceNetwork(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match Source-Networks #3")
	}
	rule1.Match.SrcNet = []*net.IPNet{testNet("10.0.0.0/24")}
	if matchSourceNetwork(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Source-Networks #4")
	}
	rule1.Match.SrcNetN = []*net.IPNet{testNet("192.168.0.0/29")}
	if matchSourceNetwork(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Source-Networks #5")
	}
	rule1.Match.SrcNetN = []*net.IPNet{testNet("10.0.0.0/29")}
	if matchSourceNetwork(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match Source-Networks #6")
	}
}

func TestMatchProtoL3(t *testing.T) {
	cnf.C = &cnf.Config{}
	cnf.C.Service.Debug = false

	pkt := parse.ParsedPacket{
		L3: &parse.ParsedL3{
			Proto: meta.ProtoL3IP4,
		},
	}
	rule1 := cnf.Rule{}
	rule1.Match.ProtoL3 = []meta.Proto{meta.ProtoL3IP4}
	if matchProtoL3(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match ProtoL3 #1")
	}
	rule1.Match.ProtoL3 = []meta.Proto{meta.ProtoL3IP6}
	if matchProtoL3(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match ProtoL3 #2")
	}
}

func TestMatchProtoL4(t *testing.T) {
	cnf.C = &cnf.Config{}
	cnf.C.Service.Debug = false

	pkt := parse.ParsedPacket{
		L4: &parse.ParsedL4{
			Proto: meta.ProtoL4Tcp,
		},
	}
	rule1 := cnf.Rule{}
	rule1.Match.ProtoL4N = []meta.Proto{meta.ProtoL4Udp}
	if matchProtoL4(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match ProtoL4 #1")
	}
	rule1.Match.ProtoL4 = []meta.Proto{meta.ProtoL4Tcp}
	if matchProtoL4(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match ProtoL4 #2")
	}
	rule1.Match.ProtoL4 = []meta.Proto{meta.ProtoL4Udp}
	if matchProtoL4(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match ProtoL4 #3")
	}
}

func TestMatchProtoL5(t *testing.T) {
	cnf.C = &cnf.Config{}
	cnf.C.Service.Debug = false

	pkt := parse.ParsedPacket{
		L5: &parse.ParsedL5{
			Proto:     meta.ProtoL5Http,
			Encrypted: meta.OptBoolNone,
		},
	}
	rule1 := cnf.Rule{}
	rule1.Match.ProtoL5 = []meta.Proto{meta.ProtoL5Tls}
	if matchProtoL5(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match ProtoL5 #1")
	}
	rule1.Match.ProtoL5 = []meta.Proto{meta.ProtoL5Tls, meta.ProtoL5Http}
	if matchProtoL5(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match ProtoL5 #2")
	}
}

func TestMatchDomain(t *testing.T) {
	cnf.C = &cnf.Config{}
	cnf.C.Service.Debug = false

	pkt := parse.ParsedPacket{
		L5: &parse.ParsedL5{
			Proto:     meta.ProtoL5Tls,
			Encrypted: meta.OptBoolTrue,
			TlsSni:    "random.xxx",
		},
	}
	rule1 := cnf.Rule{}
	rule1.Match.Domains = []string{"superstes.eu", "*.calamary.net"}
	if matchDomain(pkt, rule1, 1) != meta.MatchNegative {
		t.Error("Match Domain #1")
	}
	pkt.L5.TlsSni = "superstes.eu"
	if matchDomain(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Domain #2")
	}
	pkt.L5.TlsSni = "calamary.net"
	if matchDomain(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Domain #3")
	}
	pkt.L5.TlsSni = "test.calamary.net"
	if matchDomain(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Domain #4")
	}
	pkt.L5.TlsSni = "abc.test.calamary.net"
	if matchDomain(pkt, rule1, 1) != meta.MatchPositive {
		t.Error("Match Domain #5")
	}
}
