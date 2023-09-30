package u

import (
	"fmt"
	"strings"
	"time"

	"slices"
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
