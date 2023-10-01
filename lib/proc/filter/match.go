package filter

import (
	"fmt"
	"net"
	"strings"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

func matchProtoL3(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	// protocols layer 3
	if rule.Match.ProtoL3 != nil && len(rule.Match.ProtoL3) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("Proto L3: %v vs %v", rule.Match.ProtoL3, pkt.L3.Proto))
		return anyProtoMatch(rule.Match.ProtoL3, pkt.L3.Proto)
	}
	if rule.Match.ProtoL3N != nil && len(rule.Match.ProtoL3N) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!Proto L3: %v vs %v", rule.Match.ProtoL3N, pkt.L3.Proto))
		return negMatch(anyProtoMatch(rule.Match.ProtoL3N, pkt.L3.Proto))
	}
	return meta.MatchNeutral
}

func matchProtoL4(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	// protocols layer 4
	if rule.Match.ProtoL4 != nil && len(rule.Match.ProtoL4) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("Proto L4: %v vs %v", rule.Match.ProtoL4, pkt.L4.Proto))
		return anyProtoMatch(rule.Match.ProtoL4, pkt.L4.Proto)
	}
	if rule.Match.ProtoL4N != nil && len(rule.Match.ProtoL4N) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!Proto L4: %v vs %v", rule.Match.ProtoL4N, pkt.L4.Proto))
		return negMatch(anyProtoMatch(rule.Match.ProtoL4N, pkt.L4.Proto))
	}
	return meta.MatchNeutral
}

func matchProtoL5(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	result := meta.MatchNeutral
	// protocols layer 5
	if rule.Match.ProtoL5 != nil && pkt.L5.Proto != meta.ProtoNone {
		ruleDebug(pkt, rid, fmt.Sprintf("Proto L5: %v vs %v", rule.Match.ProtoL5, pkt.L5.Proto))
		result = anyProtoMatch(rule.Match.ProtoL5, pkt.L5.Proto)
		if result == meta.MatchNegative {
			return result
		}
	}
	if rule.Match.ProtoL5N != nil && pkt.L5.Proto != meta.ProtoNone {
		ruleDebug(pkt, rid, fmt.Sprintf("!Proto L5: %v vs %v", rule.Match.ProtoL5N, pkt.L5.Proto))
		result = anyProtoMatch(rule.Match.ProtoL5N, pkt.L5.Proto)
		if result == meta.MatchNegative {
			return result
		}
	}
	// todo: make sure 'none' is OK for 'pkt.L5.Encrypted'
	if rule.Match.Encrypted != meta.OptBoolNone && pkt.L5.Encrypted != meta.OptBoolNone {
		ruleDebug(pkt, rid, fmt.Sprintf("Encrypted: %v vs %v", rule.Match.Encrypted, pkt.L5.Encrypted))
		if rule.Match.Encrypted != pkt.L5.Encrypted {
			return meta.MatchNegative
		} else {
			result = meta.MatchPositive
		}
	}
	return result
}

// save result to handle if excluded subnet is inside included subnet
func matchSourceNetwork(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	result := meta.MatchNeutral
	// source network
	if rule.Match.SrcNet != nil && len(rule.Match.SrcNet) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("SNet: %v vs %v", rule.Match.SrcNet, pkt.L3.SrcIP))
		result = anyNetMatch(rule.Match.SrcNet, pkt.L3.SrcIP)
	}
	if rule.Match.SrcNetN != nil && len(rule.Match.SrcNetN) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!SNet: %v vs %v", rule.Match.SrcNetN, pkt.L3.SrcIP))
		return negMatch(anyNetMatch(rule.Match.SrcNetN, pkt.L3.SrcIP))
	}
	return result
}

func matchDestinationNetwork(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	result := meta.MatchNeutral
	// destination network
	if rule.Match.DestNet != nil && len(rule.Match.DestNet) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("DNet: %v vs %v", rule.Match.DestNet, pkt.L3.DestIP))
		result = anyNetMatch(rule.Match.DestNet, pkt.L3.DestIP)
	}
	if rule.Match.DestNetN != nil && len(rule.Match.DestNetN) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!DNet: %v vs %v", rule.Match.DestNetN, pkt.L3.DestIP))
		return negMatch(anyNetMatch(rule.Match.DestNetN, pkt.L3.DestIP))
	}
	return result
}

func matchSourcePort(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	// source port
	if rule.Match.SrcPort != nil && len(rule.Match.SrcPort) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("SPort: %v vs %v", rule.Match.SrcPort, pkt.L4.SrcPort))
		return anyPortMatch(rule.Match.SrcPort, pkt.L4.SrcPort)
	}
	if rule.Match.SrcPortN != nil && len(rule.Match.SrcPortN) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!SPort: %v vs %v", rule.Match.SrcPortN, pkt.L4.SrcPort))
		return negMatch(anyPortMatch(rule.Match.SrcPortN, pkt.L4.SrcPort))
	}
	return meta.MatchNeutral
}

func matchDestinationPort(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	// destination port
	if rule.Match.DestPort != nil && len(rule.Match.DestPort) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("DPort: %v vs %v", rule.Match.DestPort, pkt.L4.DestPort))
		return anyPortMatch(rule.Match.DestPort, pkt.L4.DestPort)
	}
	if rule.Match.DestPortN != nil && len(rule.Match.DestPortN) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!DPort: %v vs %v", rule.Match.DestPortN, pkt.L4.DestPort))
		return negMatch(anyPortMatch(rule.Match.DestPortN, pkt.L4.DestPort))
	}
	return meta.MatchNeutral
}

func matchDomain(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	if rule.Match.Domains != nil && len(rule.Match.Domains) > 0 {
		if pkt.L5.Proto == meta.ProtoL5Tls {
			return anyDomainMatch(rule.Match.Domains, pkt.L5.TlsSni)
		}
		// NOTE: domains from plain http host-headers are ignored by design as they can be modified easily
		//   no important dataflow should use plain HTTP anyway - just move to HTTPS already..
	}
	return meta.MatchNeutral
}

func anyProtoMatch(list []meta.Proto, single meta.Proto) meta.Match {
	for i := range list {
		if list[i] == single {
			return meta.MatchPositive
		}
	}
	return meta.MatchNegative
}

func anyPortMatch(list []uint16, single uint16) meta.Match {
	for i := range list {
		if list[i] == single {
			return meta.MatchPositive
		}
	}
	return meta.MatchNegative
}

func anyNetMatch(nets []*net.IPNet, ip net.IP) meta.Match {
	for i := range nets {
		if nets[i].Contains(ip) {
			return meta.MatchPositive
		}
	}
	return meta.MatchNegative
}

func anyDomainMatch(domains []string, domain string) meta.Match {
	for _, matchDomain := range domains {
		if strings.HasPrefix(matchDomain, "*.") {
			matchDomain = strings.Replace(matchDomain, "*.", "", 1)
			if strings.HasSuffix(domain, matchDomain) {
				return meta.MatchPositive
			}
		} else {
			if matchDomain == domain {
				return meta.MatchPositive
			}
		}
	}
	return meta.MatchNegative
}

func negMatch(match meta.Match) meta.Match {
	if match == meta.MatchPositive {
		return meta.MatchNegative
	} else if match == meta.MatchNegative {
		return meta.MatchPositive
	}
	return meta.MatchNegative
}
