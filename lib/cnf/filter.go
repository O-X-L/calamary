package cnf

import (
	"fmt"
	"net"
	"strings"

	"github.com/superstes/calamary/proc/meta"
)

func filterAction(configAction string) meta.Action {
	switch strings.ToLower(configAction) {
	case "accept", "allow":
		return meta.ActionAccept
	default:
		return meta.ActionDeny
	}
}

func matchProtoL3(configProto string) meta.Proto {
	switch strings.ToLower(configProto) {
	case "ip4", "ipv4":
		return meta.ProtoL3IP4
	case "ip6", "ipv6":
		return meta.ProtoL3IP6
	default:
		panic(fmt.Errorf("protoL3 '%v' not found", configProto))
	}
}

func matchProtoL4(configProto string) meta.Proto {
	switch strings.ToLower(configProto) {
	case "tcp":
		return meta.ProtoL4Tcp
	case "udp":
		return meta.ProtoL4Udp
	default:
		panic(fmt.Errorf("protoL4 '%v' not found or not yet supported", configProto))
	}
}

func matchProtoL5(configProto string) meta.Proto {
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
		panic(fmt.Errorf("protoL5 '%v' not found or not yet supported", configProto))
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

func matchIPs(configIPs []string) (networks []*net.IPNet) {
	// todo: allow users to provide single ip or list
	for i := range configIPs {
		ip := configIPs[i]
		if strings.Contains(ip, "/") {
			_, netip, err := net.ParseCIDR(ip)
			if err != nil {
				panic(fmt.Errorf("IP-network '%v' could not be parsed (must be valid CIDR-notation)", ip))
			}
			networks = append(
				networks,
				netip,
			)

		} else if strings.Contains(ip, ".") {
			netip := net.ParseIP(ip)
			if netip == nil {
				panic(fmt.Errorf("IPv4 '%v' could not be parsed", ip))
			}
			networks = append(
				networks,
				&net.IPNet{
					IP:   netip,
					Mask: net.CIDRMask(32, 32),
				},
			)

		} else {
			netip := net.ParseIP(ip)
			if netip == nil {
				panic(fmt.Errorf("IPv6 '%v' could not be parsed", ip))
			}
			networks = append(
				networks,
				&net.IPNet{
					IP:   netip,
					Mask: net.CIDRMask(128, 128),
				},
			)
		}
	}

	return
}
