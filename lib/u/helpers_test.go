package u

import (
	"fmt"
	"net"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/superstes/calamary/cnf"
)

var TestIPv4_1 = "1.1.1.1"
var TestIPv4_2 = "1.0.0.1"
var TestIPv6_1 = "2606:4700:4700::1001"
var TestIPv6_2 = "2606:4700:4700::1111"

func TestAllStrInList(t *testing.T) {
	a := []string{"test", "abc"}
	b := []string{"test", "xyz", "abc"}
	if !AllStrInList(b, a) {
		t.Errorf("Failed to compare slices #1")
	}
	if AllStrInList(a, b) {
		t.Errorf("Failed to compare slices #2")
	}
}

func TestToStr(t *testing.T) {
	e := fmt.Errorf("test")
	if ToStr(e) != "test" {
		t.Errorf("Failed convert to string #1")
	}

	if ToStr(11) != "11" {
		t.Errorf("Failed convert to string #2")
	}
}

func TestIsDomainName(t *testing.T) {
	if !IsDomainName("calamary.net") {
		t.Errorf("Failed validate domain #1")
	}
	if !IsDomainName("sub.calamary.net") {
		t.Errorf("Failed validate domain #2")
	}
	if !IsDomainName("othersub.sub.calamary.net") {
		t.Errorf("Failed validate domain #3")
	}
	if !IsDomainName("129a.calamary.net") {
		t.Errorf("Failed validate domain #4")
	}
	if IsDomainName("ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeerrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrssssssssssssssssssssssesassssssssub.calamary.net") {
		t.Errorf("Failed validate domain #5")
	}
	if IsDomainName("subddddddddddddddddddddddddddeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeerrrrrrrrrrrrrrrrrrrrr.calamary.net") {
		t.Errorf("Failed validate domain #6")
	}
	if IsDomainName(TestIPv4_1) {
		t.Errorf("Failed validate domain #7")
	}
	if IsDomainName(TestIPv6_1) {
		t.Errorf("Failed validate domain #8")
	}
}

func TestTimeout(t *testing.T) {
	testTime := uint(33)
	to := Timeout(testTime)
	comp := time.Duration(int(testTime) * int(time.Millisecond))
	if to != comp {
		t.Errorf("Failed get timeout #1")
	}
}

func TestIn(t *testing.T) {
	a := []string{"test", "abc"}
	if !IsIn("test", a) {
		t.Errorf("Failed to find in slice #1")
	}
	if IsIn("testa", a) {
		t.Errorf("Failed to find in slice #2")
	}
}

func TestInStr(t *testing.T) {
	if !IsInStr("test", "justAtest") {
		t.Errorf("Failed to find substring #1")
	}
	if IsInStr("no", "justAtest") {
		t.Errorf("Failed to find substring #2")
	}
}

func TestDnsResolveWithServer(t *testing.T) {
	dnsResolveWithServer(TestIPv4_1)
	dnsResolveWithServer(TestIPv4_1 + ":53")
	dnsResolveWithServer(TestIPv4_1 + ":853")
}

func TestIPv4(t *testing.T) {
	if !IsIPv4(TestIPv4_1) {
		t.Errorf("Failed to validate IPv4 #1")
	}
	if IsIPv4(TestIPv6_1) {
		t.Errorf("Failed to validate IPv4 #2")
	}
}

func TestFormatIPv6(t *testing.T) {
	ip := FormatIPv6(TestIPv6_1)
	if ip != ("[" + TestIPv6_1 + "]") {
		t.Errorf("Failed to format IPv6 #1")
	}

	ip = FormatIPv6(TestIPv4_1)
	if ip != TestIPv4_1 {
		t.Errorf("Failed to format IPv6 #2")
	}
}

func TestIsIP(t *testing.T) {
	valid, ip := IsIP(TestIPv4_1)
	if !valid || ip.String() != TestIPv4_1 {
		t.Errorf("Failed to check for IP #1")
	}

	valid, ip = IsIP("test.at")
	if valid {
		t.Errorf("Failed to check for IP #2")
	}

	valid, ip = IsIP(TestIPv6_1)
	if !valid || ip.String() != TestIPv6_1 {
		t.Errorf("Failed to check for IP #3")
	}
}

func TestDnsLookup(t *testing.T) {
	// NOTE: tests will fail if the IPs change.. should not be common
	cnf.C = &cnf.Config{
		Service: cnf.ServiceConfig{
			DnsNameservers: []string{TestIPv4_1},
		},
	}
	resp := DnsLookup("one.one.one.one")

	if !isInIpList(TestIPv4_1, resp) || !isInIpList(TestIPv4_2, resp) ||
		!isInIpList(TestIPv6_1, resp) || !isInIpList(TestIPv6_2, resp) {
		t.Errorf("DNS Query #1 has unexpected result: %v", resp)
	}

	resp4, resp6 := DnsLookup46("one.one.one.one")

	if !isInIpList(TestIPv4_1, resp4) || !isInIpList(TestIPv4_2, resp4) ||
		!isInIpList(TestIPv6_1, resp6) || !isInIpList(TestIPv6_2, resp6) ||
		isInIpList(TestIPv4_1, resp6) || isInIpList(TestIPv6_1, resp4) {
		t.Errorf("DNS Query #2 has unexpected result: ip4 %v, ip6 %v", resp4, resp6)
	}

	// edge-cases - input is already IP
	resp = DnsLookup(TestIPv4_1)
	if len(resp) != 1 || !isInIpList(TestIPv4_1, resp) {
		t.Errorf("DNS Query #3 has unexpected result: %v", resp)
	}

	resp = DnsLookup(TestIPv6_1)
	if len(resp) != 1 || !isInIpList(TestIPv6_1, resp) {
		t.Errorf("DNS Query #4 has unexpected result: %v", resp)
	}

	resp4, resp6 = DnsLookup46(TestIPv4_1)
	if len(resp4) != 1 || len(resp6) != 0 || !isInIpList(TestIPv4_1, resp4) {
		t.Errorf("DNS Query #5 has unexpected result: %v", resp4)
	}

	resp4, resp6 = DnsLookup46(TestIPv6_1)
	if len(resp6) != 1 || len(resp4) != 0 || !isInIpList(TestIPv6_1, resp6) {
		t.Errorf("DNS Query #6 has unexpected result: %v", resp6)
	}
}

func isInIpList(value string, list []net.IP) bool {
	for i := range list {
		if list[i].String() == value {
			return true
		}
	}
	return false
}

func TestDumpToFile(t *testing.T) {
	DumpToFile([]byte{'t', 'e', 's', 't'})
}

func TestTrustedCAs(t *testing.T) {
	_, pathToTest, _, _ := runtime.Caller(0)
	pathToTestCerts := filepath.Dir(pathToTest) + "/testdata/"

	cnf.C = &cnf.Config{
		Service: cnf.ServiceConfig{
			Certs: cnf.ServiceCertificates{
				CAPublic: pathToTestCerts + "ca.crt",
			},
		},
	}
	TrustedCAs()
}
