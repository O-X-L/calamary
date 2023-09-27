package parse

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/meta"
	tls_dissector "github.com/superstes/calamary/proc/parse/tls"
)

func parseTls(pkt ParsedPacket, conn net.Conn, connIo io.ReadWriter, hdr [cnf.L5HDRLEN]byte) (isTls meta.OptBool, tlsVersion uint16, sni string) {
	isTlsRaw := hdr[0] == tls_dissector.Handshake
	tlsVersion = binary.BigEndian.Uint16(hdr[1:3])
	if isTlsRaw {
		isTls = meta.OptBoolTrue
	} else {
		isTls = meta.OptBoolFalse
	}

	if isTlsRaw {
		buf := new(bytes.Buffer)
		r := io.TeeReader(connIo, buf)
		record, err := tls_dissector.ReadRecord(r)
		if err != nil {
			return
		}

		clientHello := tls_dissector.ClientHelloMsg{}
		if err = clientHello.Decode(record.Opaque); err != nil {
			return
		}

		for _, ext := range clientHello.Extensions {
			if ext.Type() == tls_dissector.ExtServerName {
				snExtension := ext.(*tls_dissector.ServerNameExtension)
				sni = snExtension.Name
				break
			}
		}
	}
	log.ConnDebug("parse", PktSrc(pkt), PktDest(pkt), fmt.Sprintf(
		"TLS information: IsTls=%v, TlsVersion=%v, TlsSni=%s", isTlsRaw, tlsVersion, sni,
	))
	return
}
