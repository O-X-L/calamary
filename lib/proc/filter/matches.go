package filter

import (
	"net"

	"github.com/superstes/calamary/proc/meta"
)

type Match struct {
	srcIP     []net.IPNet
	srcIPN    []net.IPNet
	destIP    []net.IPNet
	destIPN   []net.IPNet
	srcPort   []int
	srcPortN  []int
	destPort  []int
	destPortN []int
	protoL3   []meta.Proto
	protoL3N  []meta.Proto
	protoL4   []meta.Proto
	protoL4N  []meta.Proto
	protoL5   []meta.Proto
	protoL5N  []meta.Proto
}
