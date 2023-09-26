package parse

import (
	"net"

	"github.com/superstes/calamary/proc/meta"
)

type ParsedPackage struct {
	L3     *ParsedL3Package
	L4     *ParsedL4Package
	L5     *ParsedL5Package
	L4Tcp  *ParsedTcpPackage
	L4Udp  *ParsedUdpPackage
	L5Http *ParsedHttpPackage
	L5Dns  *ParsedDnsPackage
}

type ParsedL3Package struct {
	Proto  meta.Proto
	SrcIP  net.IP
	DestIP net.IP
}

type ParsedL4Package struct {
	Proto    meta.Proto
	SrcPort  uint16
	DestPort uint16
}

type ParsedL5Package struct {
	Proto     meta.Proto
	Encrypted meta.OptBool
}

type ParsedTcpPackage struct {
}

type ParsedUdpPackage struct {
}

type ParsedHttpPackage struct {
	// plaintext or intercepted HTTP
	/*
		method
	*/
	Headers  string
	MimeType string
	AuthUser string
	AuthPwd  string
}

type ParsedTlsPackage struct {
	Valid           bool
	Trusted         bool
	Expired         bool
	Fingerprint     string
	Issuer          string
	Sni             string
	CommonName      string
	Ou              string
	Org             string
	Country         string
	SubjectAltNames []string
}

type ParsedDnsPackage struct {
	Record     string
	RecordType string
	Ttl        int
}
