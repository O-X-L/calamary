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
	tlsCnf := &tls.Config{
		MinVersion: tls.VersionTLS11,
		NextProtos: []string{"h2", "http/1.1"},
		/*
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		*/
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
