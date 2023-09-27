package meta

import (
	"testing"
)

func TestMatchFilterAction(t *testing.T) {
	if RuleAction("allow") != ActionAccept || RuleAction("accept") != ActionAccept {
		t.Error("Filter action accept")
	}
	if RuleAction("deny") != ActionDeny || RuleAction("random") != ActionDeny {
		t.Error("Filter action deny")
	}
}

func TestMatchEncryption(t *testing.T) {
	if MatchEncrypted("true") != OptBoolTrue || MatchEncrypted("Yes") != OptBoolTrue || MatchEncrypted("1") != OptBoolTrue {
		t.Error("Match encryption true")
	}
	if MatchEncrypted("false") != OptBoolFalse || MatchEncrypted("No") != OptBoolFalse || MatchEncrypted("0") != OptBoolFalse {
		t.Error("Match encryption false")
	}
	if MatchEncrypted("random") != OptBoolNone || MatchEncrypted("") != OptBoolNone {
		t.Error("Match encryption none")
	}
}

func TestMatchProtoL3(t *testing.T) {
	if MatchProtoL3("ip4") != ProtoL3IP4 || MatchProtoL3("ipv4") != ProtoL3IP4 || MatchProtoL3("IPv4") != ProtoL3IP4 {
		t.Error("Match proto L3 ip4")
	}
	if MatchProtoL3("ip6") != ProtoL3IP6 || MatchProtoL3("ipv6") != ProtoL3IP6 || MatchProtoL3("IPv6") != ProtoL3IP6 {
		t.Error("Match proto L3 ip6")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match proto L3 not failing on invalid value")
		}
	}()
	MatchProtoL3("random")
}

func TestMatchProtoL4(t *testing.T) {
	if MatchProtoL4("tcp") != ProtoL4Tcp || MatchProtoL4("TCP") != ProtoL4Tcp {
		t.Error("Match proto L4 tcp")
	}
	if MatchProtoL4("udp") != ProtoL4Udp || MatchProtoL4("UDP") != ProtoL4Udp {
		t.Error("Match proto L4 udp")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match proto L4 not failing on invalid value")
		}
	}()
	MatchProtoL4("random")
}

func TestMatchProtoL5(t *testing.T) {
	if MatchProtoL5("tls") != ProtoL5Tls {
		t.Error("Match proto L5 tls")
	}
	if MatchProtoL5("http") != ProtoL5Http {
		t.Error("Match proto L5 http")
	}
	/*
		if MatchProtoL5("dns") != ProtoL5Dns {
			t.Error("Match proto L5 dns")
		}
		if MatchProtoL5("ntp") != ProtoL5Ntp {
			t.Error("Match proto L5 ntp")
		}
	*/
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Match proto L5 not failing on invalid value")
		}
	}()
	MatchProtoL5("random")
}
