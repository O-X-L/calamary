package filter

import (
	"fmt"
	"net"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

// http://www.squid-cache.org/Doc/config/acl/

func Filter(pkg parse.ParsedPackage) bool {
	for rid := range *cnf.RULES {
		rule := (*cnf.RULES)[rid]

		// protocols layer 3
		if rule.Match.ProtoL3 != nil && len(rule.Match.ProtoL3) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("Proto L3: %v vs %v", rule.Match.ProtoL3, pkg.L3.Proto))
			if !anyProtoMatch(rule.Match.ProtoL3, pkg.L3.Proto) {
				continue
			}
		}
		if rule.Match.ProtoL3N != nil && len(rule.Match.ProtoL3N) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("!Proto L3: %v vs %v", rule.Match.ProtoL3N, pkg.L3.Proto))
			if anyProtoMatch(rule.Match.ProtoL3N, pkg.L3.Proto) {
				continue
			}
		}

		// protocols layer 4
		if rule.Match.ProtoL4 != nil && len(rule.Match.ProtoL4) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("Proto L4: %v vs %v", rule.Match.ProtoL4, pkg.L4.Proto))
			if !anyProtoMatch(rule.Match.ProtoL4, pkg.L4.Proto) {
				continue
			}
		}
		if rule.Match.ProtoL4N != nil && len(rule.Match.ProtoL4N) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("!Proto L4: %v vs %v", rule.Match.ProtoL4N, pkg.L4.Proto))
			if anyProtoMatch(rule.Match.ProtoL4N, pkg.L4.Proto) {
				continue
			}
		}
		// protocols layer 5
		if rule.Match.ProtoL5 != nil && pkg.L5.Proto != meta.ProtoNone {
			ruleDebug(pkg, rid, fmt.Sprintf("Proto L5: %v vs %v", rule.Match.ProtoL5, pkg.L5.Proto))
			if !anyProtoMatch(rule.Match.ProtoL5, pkg.L5.Proto) {
				continue
			}
		}
		if rule.Match.ProtoL5N != nil && pkg.L5.Proto != meta.ProtoNone {
			ruleDebug(pkg, rid, fmt.Sprintf("!Proto L5: %v vs %v", rule.Match.ProtoL5N, pkg.L5.Proto))
			if anyProtoMatch(rule.Match.ProtoL5N, pkg.L5.Proto) {
				continue
			}
		}
		if rule.Match.Encrypted != meta.OptBoolNone && pkg.L5.Encrypted != meta.OptBoolNone {
			ruleDebug(pkg, rid, fmt.Sprintf("Encrypted: %v vs %v", rule.Match.Encrypted, pkg.L5.Encrypted))
			if rule.Match.Encrypted != pkg.L5.Encrypted {
				continue
			}
		}
		// source port
		if rule.Match.SrcPort != nil && len(rule.Match.SrcPort) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("SPort: %v vs %v", rule.Match.SrcPort, pkg.L4.SrcPort))
			if !anyPortMatch(rule.Match.SrcPort, pkg.L4.SrcPort) {
				continue
			}
		}
		if rule.Match.SrcPortN != nil && len(rule.Match.SrcPortN) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("!SPort: %v vs %v", rule.Match.SrcPortN, pkg.L4.SrcPort))
			if anyPortMatch(rule.Match.SrcPortN, pkg.L4.SrcPort) {
				continue
			}
		}
		// destination port
		if rule.Match.DestPort != nil && len(rule.Match.DestPort) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("DPort: %v vs %v", rule.Match.DestPort, pkg.L4.DestPort))
			if !anyPortMatch(rule.Match.DestPort, pkg.L4.DestPort) {
				continue
			}
		}
		if rule.Match.DestPortN != nil && len(rule.Match.DestPortN) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("!DPort: %v vs %v", rule.Match.DestPortN, pkg.L4.DestPort))
			if anyPortMatch(rule.Match.DestPortN, pkg.L4.DestPort) {
				continue
			}
		}
		// source network
		if rule.Match.SrcNet != nil && len(rule.Match.SrcNet) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("SNet: %v vs %v", rule.Match.SrcNet, pkg.L3.SrcIP))
			if !anyNetMatch(rule.Match.SrcNet, pkg.L3.SrcIP) {
				continue
			}
		}
		if rule.Match.SrcNetN != nil && len(rule.Match.SrcNetN) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("!SNet: %v vs %v", rule.Match.SrcNetN, pkg.L3.SrcIP))
			if anyNetMatch(rule.Match.SrcNetN, pkg.L3.SrcIP) {
				continue
			}
		}
		// destination network
		if rule.Match.DestNet != nil && len(rule.Match.DestNet) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("DNet: %v vs %v", rule.Match.DestNet, pkg.L3.DestIP))
			if !anyNetMatch(rule.Match.DestNet, pkg.L3.DestIP) {
				continue
			}
		}
		if rule.Match.DestNetN != nil && len(rule.Match.DestNetN) > 0 {
			ruleDebug(pkg, rid, fmt.Sprintf("!DNet: %v vs %v", rule.Match.DestNetN, pkg.L3.DestIP))
			if anyNetMatch(rule.Match.DestNetN, pkg.L3.DestIP) {
				continue
			}
		}

		ruleDebug(pkg, rid, fmt.Sprintf("Applying action '%v'", reverseFilterAction(rule.Action)))
		return applyAction(rule.Action)

	}

	// implicit deny
	log.ConnDebug("filter", parse.PkgSrc(pkg), parse.PkgDest(pkg), "No rule matched - implicit deny")
	return applyAction(meta.ActionDeny)
}

func ruleDebug(pkg parse.ParsedPackage, rule_id int, msg string) {
	log.ConnDebug("filter", parse.PkgSrc(pkg), parse.PkgDest(pkg), fmt.Sprintf("Rule %v - %s", rule_id, msg))
}

func applyAction(action meta.Action) bool {
	if action == meta.ActionAccept {
		return true
	}
	return false
}

func anyProtoMatch(list []meta.Proto, single meta.Proto) bool {
	for i := range list {
		if list[i] == single {
			return true
		}
	}
	return false
}

func anyPortMatch(list []uint16, single uint16) bool {
	for i := range list {
		if list[i] == single {
			return true
		}
	}
	return false
}

func allNetMatch(nets []*net.IPNet, ip net.IP) bool {
	for i := range nets {
		if !nets[i].Contains(ip) {
			return false
		}
	}
	return true
}

func anyNetMatch(nets []*net.IPNet, ip net.IP) bool {
	for i := range nets {
		if nets[i].Contains(ip) {
			return true
		}
	}
	return false
}

func reverseFilterAction(action meta.Action) string {
	// todo: merge with config-parser function to keep action matching in one place
	switch action {
	case meta.ActionAccept:
		return "accept"
	case meta.ActionDeny:
		return "deny"
	default:
		return "unknown"
	}
}
