package cnf

import (
	"net"
	"time"
)

const (
	VERSION         string = "1.0"
	LOG_TIME_FORMAT string = "2006-01-02 15:04:05"
	UDP_TTL                = 30 * time.Second
	UDP_BUFFER_SIZE int    = 4096
	BYTES_HDR_L5    int    = 5
	BYTES_TLS_REC   int    = 5
	BYTES_TLS_HS    int    = 4
	BYTES_TLS_EXT   int    = 4
)

var ConfigFileAbs string = "/etc/calamary/config.yml"

var NetForwardDeny []*net.IPNet

func InitNetForwardDeny() {
	_, localhost1, _ := net.ParseCIDR("127.0.0.0/8")
	_, localhost2, _ := net.ParseCIDR("::1/128")
	_, localhost3, _ := net.ParseCIDR("::/128")
	_, linklocal1, _ := net.ParseCIDR("169.254.0.0/16")
	_, linklocal2, _ := net.ParseCIDR("fe80::/10")
	_, linklocal3, _ := net.ParseCIDR("fc00::/7")
	_, multicast1, _ := net.ParseCIDR("224.0.0.0/4")
	_, multicast2, _ := net.ParseCIDR("ff00::/8")
	_, broadcast1, _ := net.ParseCIDR("255.255.255.255/32")
	_, blackhole1, _ := net.ParseCIDR("100::/64")

	NetForwardDeny = []*net.IPNet{
		localhost1,
		localhost2,
		localhost3,
		linklocal1,
		linklocal2,
		linklocal3,
		multicast1,
		multicast2,
		broadcast1,
		blackhole1,
	}
}
