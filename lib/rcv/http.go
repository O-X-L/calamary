package rcv

import (
	"crypto/tls"
	"fmt"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
)

func newServerHttpTcp(addr string, lncnf cnf.ServiceListener) (Server, error) {
	transparentSrv, err := newServerTransparentTcp(addr, lncnf)
	if err != nil {
		return transparentSrv, err
	}
	return Server{
		Listener: transparentSrv.Listener,
		Cnf:      lncnf,
		L4Proto:  meta.ProtoL4Tcp,
	}, nil
}

func newServerHttpsTcp(addr string, lncnf cnf.ServiceListener) (Server, error) {
	keyPair, err := tls.LoadX509KeyPair(
		cnf.C.Service.Certs.ServerPublic,
		cnf.C.Service.Certs.ServerPrivate,
	)
	if err != nil {
		return Server{}, fmt.Errorf("Failed to load certificates")
	}
	tlsCnf := &tls.Config{
		MinVersion:   tls.VersionTLS10,
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{keyPair},
	}
	ln, err := tls.Listen(
		"tcp",
		fmt.Sprintf("%v:%v", addr, lncnf.Port),
		tlsCnf,
	)
	if err != nil {
		return Server{}, err
	}
	return Server{
		Listener: ln,
		Cnf:      lncnf,
		L4Proto:  meta.ProtoL4Tcp,
	}, nil
}
