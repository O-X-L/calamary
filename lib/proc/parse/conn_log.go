package parse

import "github.com/superstes/calamary/log"

func LogConnDebug(pkg string, pkt ParsedPacket, msg string) {
	log.Conn("DEBUG", pkg, PktSrc(pkt), PktDest(pkt), msg)
}

func LogConnInfo(pkg string, pkt ParsedPacket, msg string) {
	log.Conn("INFO", pkg, PktSrc(pkt), PktDest(pkt), msg)
}

func LogConnWarn(pkg string, pkt ParsedPacket, msg string) {
	log.Conn("WARN", pkg, PktSrc(pkt), PktDest(pkt), msg)
}

func LogConnError(pkg string, pkt ParsedPacket, msg string) {
	log.Conn("ERROR", pkg, PktSrc(pkt), PktDest(pkt), msg)
}
