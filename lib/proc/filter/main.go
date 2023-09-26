package filter

import (
	"fmt"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

// http://www.squid-cache.org/Doc/config/acl/

func Filter(pkt parse.ParsedPacket) bool {
	for rid := range *cnf.RULES {
		rule := (*cnf.RULES)[rid]

		// go to next rule if match is defined and packet missed it
		if !matchProtoL3(pkt, rule, rid) || !matchProtoL4(pkt, rule, rid) || !matchProtoL5(pkt, rule, rid) {
			continue
		}
		if !matchSourceNetwork(pkt, rule, rid) || !matchDestinationNetwork(pkt, rule, rid) {
			continue
		}
		if !matchSourcePort(pkt, rule, rid) || !matchDestinationPort(pkt, rule, rid) {
			continue
		}

		ruleDebug(pkt, rid, fmt.Sprintf("Applying action '%v'", reverseFilterAction(rule.Action)))
		return applyAction(rule.Action)

	}

	// implicit deny
	log.ConnDebug("filter", parse.PkgSrc(pkt), parse.PkgDest(pkt), "No rule matched - implicit deny")
	return applyAction(meta.ActionDeny)
}

func ruleDebug(pkt parse.ParsedPacket, rule_id int, msg string) {
	if cnf.C.Service.Debug {
		log.ConnDebug("filter", parse.PkgSrc(pkt), parse.PkgDest(pkt), fmt.Sprintf("Rule %v - %s", rule_id, msg))
	}
}

func applyAction(action meta.Action) bool {
	if action == meta.ActionAccept {
		return true
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
