package u

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"slices"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
)

func ToStr(data any) string {
	return fmt.Sprintf("%v", data)
}

func AllStrInList(list []string, check []string) bool {
	for i := range check {
		if !slices.Contains(list, check[i]) {
			return false
		}
	}
	return true
}

func IsIPv4(address string) bool {
	return strings.Contains(address, ".")
}

func IsDomainName(s string) bool {
	// source: https://github.com/golang/go/blob/go1.20.5/src/net/dnsclient.go#L72-L75
	if s == "." {
		return true
	}
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	nonNumeric := false // true once we've seen a letter or hyphen
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			nonNumeric = true
			partlen++
		case '0' <= c && c <= '9':
			partlen++
		case c == '-':
			if last == '.' {
				return false
			}
			partlen++
			nonNumeric = true
		case c == '.':
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return nonNumeric
}

func Timeout(timeout uint) time.Duration {
	return time.Duration(int(timeout) * int(time.Millisecond))
}

func IsIn(value string, list []string) bool {
	for i := range list {
		if list[i] == value {
			return true
		}
	}
	return false
}

// just as shorter version
func IsInStr(searchFor string, searchIn string) bool {
	return strings.Contains(searchIn, searchFor)
}

func dnsResolveWithServer(srv string) *net.Resolver {
	if !strings.Contains(srv, ":") {
		srv = srv + ":53"
	}
	return &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: Timeout(cnf.C.Service.Timeout.DnsLookup),
			}
			return d.DialContext(ctx, network, srv)
		},
	}
}

func DnsLookup(dns string) (ips []net.IP) {
	var err error
	for _, srv := range cnf.C.Service.DnsNameservers {
		ips, err = dnsResolveWithServer(srv).LookupIP(
			context.Background(), "ip", dns,
		)
		if err != nil {
			log.Debug("util", fmt.Sprintf("Failed to lookup DNS '%s' via server %s: %v", dns, srv, err))
			continue
		}
		if len(ips) > 0 {
			break
		}
	}
	if len(ips) == 0 {
		log.ErrorS("util", fmt.Sprintf("Failed to lookup DNS '%s'", dns))
		return
	}
	log.Debug("util", fmt.Sprintf("DNS '%s' resolved to: %v", dns, ips))
	return ips
}

func DnsLookup46(dns string) (ip4 []net.IP, ip6 []net.IP) {
	for _, ip := range DnsLookup(dns) {
		if IsIPv4(ip.String()) {
			ip4 = append(ip4, ip)
		} else {
			ip6 = append(ip6, ip)
		}
	}
	return
}

func FormatIPv6(ip string) string {
	if !IsIPv4(ip) && !strings.Contains(ip, "[") {
		return fmt.Sprintf("[%v]", ip)
	}
	return ip
}
