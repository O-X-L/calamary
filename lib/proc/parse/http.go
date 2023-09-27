package parse

import (
	"net/http"
	"strings"

	"github.com/superstes/calamary/cnf"
)

func hdrL5Http(hdr [cnf.L5HDRLEN]byte) bool {
	s := string(hdr[:])
	return strings.HasPrefix(http.MethodGet, s[:3]) ||
		strings.HasPrefix(http.MethodPost, s[:4]) ||
		strings.HasPrefix(http.MethodPut, s[:3]) ||
		strings.HasPrefix(http.MethodDelete, s) ||
		strings.HasPrefix(http.MethodOptions, s) ||
		strings.HasPrefix(http.MethodPatch, s) ||
		strings.HasPrefix(http.MethodHead, s[:4]) ||
		strings.HasPrefix(http.MethodConnect, s) ||
		strings.HasPrefix(http.MethodTrace, s)
}
