package rcv

import (
	"fmt"
	"net"
	"net/http"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
)

type Listener interface {
	Accept() (net.Conn, error)
	Addr() net.Addr
	Close() error
}

type Server struct {
	Listener   Listener
	HttpServer *http.Server
	Cnf        cnf.ServiceListener
	L4Proto    meta.Proto
}

type ServerInfo struct {
	Tproxy  bool
	L4Proto string
	Mode    meta.ListenMode
}

func serverNotImplemented(addr string, lncnf cnf.ServiceListener) (Server, error) {
	return Server{}, fmt.Errorf("Listener type '%s' is not yet implemented!", lncnf.Mode)
}

func serverNotSupported(addr string, lncnf cnf.ServiceListener) (Server, error) {
	return Server{}, fmt.Errorf("Protocol is not supported by listener type '%s'!", lncnf.Mode)
}
