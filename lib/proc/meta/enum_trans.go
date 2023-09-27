package meta

import (
	"crypto/tls"
	"fmt"
	"strings"
)

func cleanRaw(configRaw string) (configClean string) {
	configClean = strings.ReplaceAll(configRaw, " ", "")
	configClean = strings.ReplaceAll(configClean, "!", "")
	return
}

func MatchEncrypted(configEncrypted string) OptBool {
	switch strings.ToLower(configEncrypted) {
	case "true", "yes", "y", "1":
		return OptBoolTrue
	case "false", "no", "n", "0":
		return OptBoolFalse
	default:
		return OptBoolNone
	}
}

func RuleAction(configAction string) Action {
	switch strings.ToLower(configAction) {
	case "accept", "allow":
		return ActionAccept
	default:
		return ActionDeny
	}
}

func MatchProtoL3(configProto string) Proto {
	configProto = cleanRaw(configProto)
	switch strings.ToLower(configProto) {
	case "ip4", "ipv4":
		return ProtoL3IP4
	case "ip6", "ipv6":
		return ProtoL3IP6
	default:
		panic(fmt.Sprintf("protoL3 '%v' not found", configProto))
	}
}

func MatchProtoL4(configProto string) Proto {
	configProto = cleanRaw(configProto)
	switch strings.ToLower(configProto) {
	case "tcp":
		return ProtoL4Tcp
	case "udp":
		return ProtoL4Udp
	default:
		panic(fmt.Sprintf("protoL4 '%v' not found or not yet supported", configProto))
	}
}

func MatchProtoL5(configProto string) Proto {
	configProto = cleanRaw(configProto)
	switch strings.ToLower(configProto) {
	case "tls":
		return ProtoL5Tls
	case "http":
		return ProtoL5Http
	/*
		case "dns":
			return ProtoL5Dns
		case "ntp":
			return ProtoL5Ntp
	*/
	default:
		panic(fmt.Sprintf("protoL5 '%v' not found or not yet supported", configProto))
	}
}

// reverse functions for logging only
func RevRuleAction(action Action) string {
	switch action {
	case ActionAccept:
		return "accept"
	case ActionDeny:
		return "deny"
	default:
		return "unknown"
	}
}

func RevProto(proto Proto) string {
	switch proto {
	case ProtoNone:
		return "none"
	case ProtoL3IP4:
		return "IPv4"
	case ProtoL3IP6:
		return "IPv6"
	case ProtoL4Tcp:
		return "TCP"
	case ProtoL4Udp:
		return "UDP"
	case ProtoL5Tls:
		return "TLS"
	case ProtoL5Http:
		return "HTTP"
	case ProtoL5Dns:
		return "DNS"
	case ProtoL5Ntp:
		return "NTP"
	default:
		return "unknown"
	}
}

func RevTlsVersion(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLSv1.0"
	case tls.VersionTLS11:
		return "TLSv1.1"
	case tls.VersionTLS12:
		return "TLSv1.2"
	case tls.VersionTLS13:
		return "TLSv1.3"
	default:
		return "none"
	}
}
