package u

import (
	"strings"
)

func IsIPv4(address string) bool {
	return strings.Contains(address, ".")
}
