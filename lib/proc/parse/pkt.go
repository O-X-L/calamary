package parse

import (
	"net"

	"github.com/superstes/calamary/proc/meta"
)

type ParsedPacket struct {
	L3     *ParsedL3
	L4     *ParsedL4
	L5     *ParsedL5
	L4Tcp  *ParsedTcp
	L4Udp  *ParsedUdp
	L5Http *ParsedHttp
	L5Dns  *ParsedDns
}

type ParsedL3 struct {
	Proto  meta.Proto
	SrcIP  net.IP
	DestIP net.IP
}

type ParsedL4 struct {
	Proto    meta.Proto
	SrcPort  uint16
	DestPort uint16
}

type ParsedL5 struct {
	Proto     meta.Proto
	Encrypted meta.OptBool
}

type ParsedTcp struct {
}

type ParsedUdp struct {
}

type ParsedHttp struct {
	// plaintext or intercepted HTTP
	/*
		method
	*/
	Headers  string
	MimeType string
	AuthUser string
	AuthPwd  string
}

type ParsedTls struct {
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

type ParsedDns struct {
	Record     string
	RecordType string
	Ttl        int
}
