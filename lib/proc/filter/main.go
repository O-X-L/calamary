package filter

import (
	"net"

	"github.com/superstes/calamary/proc/parse"
)

// http://www.squid-cache.org/Doc/config/acl/

func Filter(pkg parse.ParsedPackage) bool {
	// just test filter
	_, netip, _ := net.ParseCIDR("135.181.170.219/32")
	if netip.Contains(pkg.L3.DestIP) && pkg.L4.DestPort == 443 {
		return true
	}

	// implicit deny
	return false
}
