package parse

import (
	"net/http"
	"strings"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
)

func parseHttp(pkt ParsedPacket, hdr [cnf.BYTES_HDR_L5]byte) {
	s := string(hdr[:])
	if strings.HasPrefix(http.MethodGet, s[:3]) ||
		strings.HasPrefix(http.MethodPost, s[:4]) ||
		strings.HasPrefix(http.MethodPut, s[:3]) ||
		strings.HasPrefix(http.MethodDelete, s) ||
		strings.HasPrefix(http.MethodOptions, s) ||
		strings.HasPrefix(http.MethodPatch, s) ||
		strings.HasPrefix(http.MethodHead, s[:4]) ||
		strings.HasPrefix(http.MethodConnect, s) ||
		strings.HasPrefix(http.MethodTrace, s) {

		pkt.L5.Proto = meta.ProtoL5Http
		// todo: plain-http parsing
		pkt.L5Http = &ParsedHttp{}
	}
}

func ParseHttpPacket(pkt ParsedPacket, req *http.Request) *ParsedHttp {
	if pkt.L5.Encrypted == meta.OptBoolTrue {
		// todo: implement tls-interception
		return &ParsedHttp{}
	}

	host, port := SplitHttpHost(req.Host, pkt.L5.Encrypted)

	// todo: parse further useful information for fitlering
	return &ParsedHttp{
		Host:       host,
		Port:       port,
		Method:     req.Method,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
		Url:        req.URL,
	}
}
