package rcv

import (
	"github.com/pires/go-proxyproto"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
)

func newServerProxyProtoTcp(addr string, lncnf cnf.ServiceListener) (Server, error) {
	transparentSrv, err := newServerTransparentTcp(addr, lncnf)
	if err != nil {
		return transparentSrv, err
	}
	lnpp := &proxyproto.Listener{Listener: transparentSrv.Listener}

	return Server{
		Listener: lnpp,
		Cnf:      lncnf,
		L4Proto:  meta.ProtoL4Tcp,
	}, nil
}
