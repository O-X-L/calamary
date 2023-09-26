package cnf_file

import (
	"net"
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

func TestMatchEncryption(t *testing.T) {
	if matchEncrypted(cnf.RuleRaw{}.Match.Encypted) != meta.OptBoolNone {
		t.Error("Match encryption default-none")
	}
	if matchEncrypted("true") != meta.OptBoolTrue || matchEncrypted("Yes") != meta.OptBoolTrue || matchEncrypted("1") != meta.OptBoolTrue {
		t.Error("Match encryption true")
	}
	if matchEncrypted("false") != meta.OptBoolFalse || matchEncrypted("No") != meta.OptBoolFalse || matchEncrypted("0") != meta.OptBoolFalse {
		t.Error("Match encryption false")
	}
	if matchEncrypted("random") != meta.OptBoolNone || matchEncrypted("") != meta.OptBoolNone {
		t.Error("Match encryption none")
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

func TestMatchDomain(t *testing.T) {
	if matchDomain(" superstes.eu") != "superstes.eu" {
		t.Error("Match domain #1")
	}
	if matchDomain("www.google.com") != "www.google.com" {
		t.Error("Match domain #2")
	}
	if matchDomain(".test.at") != "*.test.at" {
		t.Error("Match domain #3")
	}
}

func TestVars(t *testing.T) {
	testVarFP := cnf.Var{ // could be false-positive match for 'test1'
		Name:  "test11",
		Value: []string{"test.com"},
	}
	testVar1 := cnf.Var{
		Name:  "test1",
		Value: []string{"192.168.0.0/16", "172.16.0.0/12"},
	}
	testVar2 := cnf.Var{
		Name:  "test2",
		Value: []string{"80", "443"},
	}
	cnf.C = &cnf.Config{
		Vars: []cnf.Var{testVarFP, testVar1, testVar2},
	}

	vf, vn, v := usedVar("no-var")
	if vf == true {
		t.Error("Var #1")
	}
	vf, vn, v = usedVar("$non-existent-var")
	if vf == true {
		t.Error("Var #2")
	}
	vf, vn, v = usedVar("$test1")
	if vf == false || vn == true || v.Name != testVar1.Name {
		t.Error("Var #3")
	}
	vf, vn, v = usedVar("!$test1")
	if vf == false || vn == false || v.Name != testVar1.Name {
		t.Error("Var #4")
	}
	vf, vn, v = usedVar("$test2")
	if vf == false || vn == true || v.Name != testVar2.Name {
		t.Error("Var #5")
	}
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

	cnf.C = &cnf.Config{
		Vars: []cnf.Var{
			cnf.Var{
				Name:  "test1",
				Value: []string{"192.168.0.0/16", "172.16.0.0/8"},
			},
		},
	}
	ParseRules([]cnf.RuleRaw{
		cnf.RuleRaw{
			Match: cnf.MatchRaw{
				DestNet: []string{"$test1"},
				SrcNet:  []string{"!$test1"},
			},
			Action: "deny",
		},
	})
}

func TestRuleHasMatches(t *testing.T) {
	r := ruleHasMatches(cnf.Rule{})
	if r == true {
		t.Error("Rule #1")
	}
	r2 := cnf.Rule{}
	r2.Match.Encrypted = meta.OptBoolNone
	r = ruleHasMatches(r2)
	if r == true {
		t.Error("Rule #2")
	}
	r3 := cnf.Rule{}
	_, r3Net, _ := net.ParseCIDR("192.168.0.0/16")
	r3.Match.SrcNet = []*net.IPNet{r3Net}
	r = ruleHasMatches(r3)
	if r == false {
		t.Error("Rule #3")
	}
	r4 := cnf.Rule{}
	r4.Match.ProtoL4 = []meta.Proto{meta.ProtoL4Tcp}
	r = ruleHasMatches(r4)
	if r == false {
		t.Error("Rule #4")
	}
}
