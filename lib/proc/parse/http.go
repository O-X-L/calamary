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
