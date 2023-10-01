package rcv

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/superstes/calamary/cnf"
	proc_http "github.com/superstes/calamary/proc/http"
	"github.com/superstes/calamary/proc/meta"
)

func newServerHttpTcp(addr string, lncnf cnf.ServiceListener) (Server, error) {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", proc_http.HandleRequest)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf("%s:%v", addr, lncnf.Port),
		Handler: httpMux,
	}
	return Server{
		HttpServer: httpSrv,
		Cnf:        lncnf,
		L4Proto:    meta.ProtoL4Tcp,
	}, nil
}

func newServerHttpsTcp(addr string, lncnf cnf.ServiceListener) (Server, error) {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", proc_http.HandleRequestTls)

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
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf("%s:%v", addr, lncnf.Port),
		Handler:      httpMux,
		TLSConfig:    tlsCnf,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	return Server{
		HttpServer: httpSrv,
		Cnf:        lncnf,
		L4Proto:    meta.ProtoL4Tcp,
	}, nil
}
