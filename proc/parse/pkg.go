package parse

import (
	"net"

	"github.com/superstes/calamary/proc/meta"
)

type ParsedPackage struct {
	L3     *ParsedL3Package
	L4     *ParsedL4Package
	L4Tcp  *ParsedTcpPackage
	L4Udp  *ParsedUdpPackage
	L5Http *ParsedHttpPackage
	L5Dns  *ParsedDnsPackage
}

type ParsedL3Package struct {
	Proto   meta.Proto
	SrcIP   net.IP
	DestIP  net.IP
	L4Proto meta.Proto
}

type ParsedL4Package struct {
	SrcPort  int
	DestPort int
	L5Proto  meta.Proto
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
