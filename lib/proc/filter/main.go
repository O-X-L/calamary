package filter

import (
	"fmt"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/metrics"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

// http://www.squid-cache.org/Doc/config/acl/

func Filter(pkt parse.ParsedPacket) bool {
	for rid := range *cnf.RULES {
		rule := (*cnf.RULES)[rid]

		if cnf.Metrics() {
			metrics.RuleHits.WithLabelValues(fmt.Sprintf("%v", rid)).Inc()
		}

		// go to next rule if match is defined and packet missed it
		if matchProtoL3(pkt, rule, rid) == meta.MatchNegative ||
			matchProtoL4(pkt, rule, rid) == meta.MatchNegative ||
			matchProtoL5(pkt, rule, rid) == meta.MatchNegative {
			continue
		}
		if matchSourceNetwork(pkt, rule, rid) == meta.MatchNegative ||
			matchDestinationNetwork(pkt, rule, rid) == meta.MatchNegative {
			continue
		}
		if matchSourcePort(pkt, rule, rid) == meta.MatchNegative ||
			matchDestinationPort(pkt, rule, rid) == meta.MatchNegative {
			continue
		}
		if matchDomain(pkt, rule, rid) == meta.MatchNegative {
			continue
		}

		actionStr := meta.RevRuleAction(rule.Action)
		if cnf.Metrics() {
			metrics.RuleMatches.WithLabelValues(fmt.Sprintf("%v", rid)).Inc()
			metrics.RuleActions.WithLabelValues(actionStr).Inc()
		}
		ruleDebug(pkt, rid, fmt.Sprintf("Applying action '%v'", actionStr))
		return applyAction(rule.Action)

	}

	// implicit deny
	log.ConnDebug("filter", parse.PktSrc(pkt), parse.PktDest(pkt), "No rule matched - implicit deny")
	return applyAction(meta.ActionDeny)
}

func ruleDebug(pkt parse.ParsedPacket, rule_id int, msg string) {
	if cnf.C.Service.Debug {
		log.ConnDebug("filter", parse.PktSrc(pkt), parse.PktDest(pkt), fmt.Sprintf("Rule %v - %s", rule_id, msg))
	}
}

func applyAction(action meta.Action) bool {
	if action == meta.ActionAccept {
		return true
	}
	return false
}
