package cnf_file

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/u"
)

func ParseRules(rawRules []cnf.RuleRaw) (rules []cnf.Rule) {
	for i := range rawRules {
		ruleRaw := rawRules[i]
		rule := cnf.Rule{
			Action: filterAction(ruleRaw.Action),
		}

		if len(ruleRaw.Match.SrcNet) > 0 {
			rule.Match.SrcNet = []*net.IPNet{}
			rule.Match.SrcNetN = []*net.IPNet{}
		}
		for i2 := range ruleRaw.Match.SrcNet {
			if ruleRaw.Match.SrcNet[i2][0] == '!' {
				rule.Match.SrcNetN = append(rule.Match.SrcNetN, matchNet(ruleRaw.Match.SrcNet[i2]))
			} else {
				rule.Match.SrcNet = append(rule.Match.SrcNet, matchNet(ruleRaw.Match.SrcNet[i2]))
			}
		}

		if len(ruleRaw.Match.DestNet) > 0 {
			rule.Match.DestNet = []*net.IPNet{}
			rule.Match.DestNetN = []*net.IPNet{}
		}
		for i2 := range ruleRaw.Match.DestNet {
			if ruleRaw.Match.DestNet[i2][0] == '!' {
				rule.Match.DestNetN = append(rule.Match.DestNetN, matchNet(ruleRaw.Match.DestNet[i2]))
			} else {
				rule.Match.DestNet = append(rule.Match.DestNet, matchNet(ruleRaw.Match.DestNet[i2]))
			}
		}

		if len(ruleRaw.Match.SrcPort) > 0 {
			rule.Match.SrcPort = []uint16{}
			rule.Match.SrcPortN = []uint16{}
		}
		for i2 := range ruleRaw.Match.SrcPort {
			rule.Match.SrcPort = []uint16{}
			rule.Match.SrcPortN = []uint16{}
			if ruleRaw.Match.SrcPort[i2][0] == '!' {
				rule.Match.SrcPortN = append(rule.Match.SrcPortN, matchPort(ruleRaw.Match.SrcPort[i2]))
			} else {
				rule.Match.SrcPort = append(rule.Match.SrcPort, matchPort(ruleRaw.Match.SrcPort[i2]))
			}
		}

		if len(ruleRaw.Match.DestPort) > 0 {
			rule.Match.DestPort = []uint16{}
			rule.Match.DestPortN = []uint16{}
		}
		for i2 := range ruleRaw.Match.DestPort {
			if ruleRaw.Match.DestPort[i2][0] == '!' {
				rule.Match.DestPortN = append(rule.Match.DestPortN, matchPort(ruleRaw.Match.DestPort[i2]))
			} else {
				rule.Match.DestPort = append(rule.Match.DestPort, matchPort(ruleRaw.Match.DestPort[i2]))
			}
		}

		if len(ruleRaw.Match.ProtoL3) > 0 {
			rule.Match.ProtoL3 = []meta.Proto{}
			rule.Match.ProtoL3N = []meta.Proto{}
		}
		for i2 := range ruleRaw.Match.ProtoL3 {
			if ruleRaw.Match.ProtoL3[i2][0] == '!' {
				rule.Match.ProtoL3N = append(rule.Match.ProtoL3N, matchProtoL3(ruleRaw.Match.ProtoL3[i2]))
			} else {
				rule.Match.ProtoL3 = append(rule.Match.ProtoL3, matchProtoL3(ruleRaw.Match.ProtoL3[i2]))
			}
		}

		if len(ruleRaw.Match.ProtoL4) > 0 {
			rule.Match.ProtoL4 = []meta.Proto{}
			rule.Match.ProtoL4N = []meta.Proto{}
		}
		for i2 := range ruleRaw.Match.ProtoL4 {
			if ruleRaw.Match.ProtoL4[i2][0] == '!' {
				rule.Match.ProtoL4N = append(rule.Match.ProtoL4N, matchProtoL4(ruleRaw.Match.ProtoL4[i2]))
			} else {
				rule.Match.ProtoL4 = append(rule.Match.ProtoL4, matchProtoL4(ruleRaw.Match.ProtoL4[i2]))
			}
		}

		if len(ruleRaw.Match.ProtoL5) > 0 {
			rule.Match.ProtoL5 = []meta.Proto{}
			rule.Match.ProtoL5N = []meta.Proto{}
		}
		for i2 := range ruleRaw.Match.ProtoL5 {
			if ruleRaw.Match.ProtoL5[i2][0] == '!' {
				rule.Match.ProtoL5N = append(rule.Match.ProtoL5N, matchProtoL5(ruleRaw.Match.ProtoL5[i2]))
			} else {
				rule.Match.ProtoL5 = append(rule.Match.ProtoL5, matchProtoL5(ruleRaw.Match.ProtoL5[i2]))
			}
		}

		if len(ruleRaw.Match.Domains) > 0 {
			rule.Match.Domains = []string{}
		}
		for i2 := range ruleRaw.Match.Domains {
			if ruleRaw.Match.Domains[i2][0] == '!' {
				rule.Match.Domains = append(rule.Match.Domains, matchDomain(ruleRaw.Match.Domains[i2]))
			} else {
				rule.Match.Domains = append(rule.Match.Domains, matchDomain(ruleRaw.Match.Domains[i2]))
			}
		}

		rules = append(rules, rule)
	}
	return
}

func cleanRaw(configRaw string) (configClean string) {
	configClean = strings.ReplaceAll(configRaw, " ", "")
	configClean = strings.ReplaceAll(configClean, "!", "")
	return
}

func filterAction(configAction string) meta.Action {
	switch strings.ToLower(configAction) {
	case "accept", "allow":
		return meta.ActionAccept
	default:
		return meta.ActionDeny
	}
}

func matchProtoL3(configProto string) meta.Proto {
	configProto = cleanRaw(configProto)
	switch strings.ToLower(configProto) {
	case "ip4", "ipv4":
		return meta.ProtoL3IP4
	case "ip6", "ipv6":
		return meta.ProtoL3IP6
	default:
		panic(fmt.Sprintf("protoL3 '%v' not found", configProto))
	}
}

func matchProtoL4(configProto string) meta.Proto {
	configProto = cleanRaw(configProto)
	switch strings.ToLower(configProto) {
	case "tcp":
		return meta.ProtoL4Tcp
	case "udp":
		return meta.ProtoL4Udp
	default:
		panic(fmt.Sprintf("protoL4 '%v' not found or not yet supported", configProto))
	}
}

func matchProtoL5(configProto string) meta.Proto {
	configProto = cleanRaw(configProto)
	switch strings.ToLower(configProto) {
	case "tls":
		return meta.ProtoL5Tls
	case "http":
		return meta.ProtoL5Http
	/*
		case "dns":
			return meta.ProtoL5Dns
		case "ntp":
			return meta.ProtoL5Ntp
	*/
	default:
		panic(fmt.Sprintf("protoL5 '%v' not found or not yet supported", configProto))
	}
}

/*
func matchIPsRaw(configIPs interface{}) (networks []*net.IPNet) {
	switch reflect.TypeOf(configIPs).Kind() {
	case reflect.String:
		return matchIPs(strings.Split(fmt.Sprintf("%v", configIPs), ","))
	case reflect.Slice, reflect.Array:
		return matchIPs(configIPs)
	default:
		panic(fmt.Errorf("IP-match '%v' neither of type string nor array/slice", configIPs))
	}
}
*/

func matchNet(ip string) *net.IPNet {
	ip = cleanRaw(ip)
	// todo: allow users to provide single ip or list
	if strings.Contains(ip, "/") {
		_, netip, err := net.ParseCIDR(ip)
		if err != nil {
			panic(fmt.Sprintf("IP-network '%s' could not be parsed (must be valid CIDR-notation)", ip))
		}
		return netip

	} else if strings.Contains(ip, ".") {
		netip := net.ParseIP(ip)
		if netip == nil {
			panic(fmt.Sprintf("IPv4 '%s' could not be parsed", ip))
		}
		return &net.IPNet{
			IP:   netip,
			Mask: net.CIDRMask(32, 32),
		}

	} else {
		netip := net.ParseIP(ip)
		if netip == nil {
			panic(fmt.Sprintf("IPv6 '%s' could not be parsed", ip))
		}
		return &net.IPNet{
			IP:   netip,
			Mask: net.CIDRMask(128, 128),
		}
	}
}

func matchPort(configPort string) uint16 {
	configPort = cleanRaw(configPort)
	port, err := strconv.ParseUint(configPort, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("Port '%s' could not be parsed", configPort))
	}
	if port > 65535 {
		panic(fmt.Sprintf("Port '%s' outside of valid range", configPort))
	}
	return uint16(port)
}

func matchDomain(configDomain string) string {
	configDomain = cleanRaw(configDomain)
	validateDomain := configDomain
	if strings.HasPrefix(configDomain, ".") {
		validateDomain = strings.Replace(configDomain, ".", "", 1)
		configDomain = strings.Replace(configDomain, ".", "*.", 1)
	}
	if !u.IsDomainName(validateDomain) {
		panic(fmt.Sprintf("Domain '%s' is not valid", validateDomain))
	}
	return configDomain
}
