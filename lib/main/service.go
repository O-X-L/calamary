package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/cnf/cnf_file"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/fwd"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/rcv"
	"github.com/superstes/calamary/u"
)

type service struct {
	servers []rcv.Server
}

func (svc *service) signalHandler() {
	signalCh := make(chan os.Signal, 1024)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		select {
		case s := <-signalCh:
			switch s {
			case syscall.SIGHUP:
				log.Warn("service", "Received reload signal")
				cnf_file.Load(false, false)

			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
				log.Warn("service", "Received shutdown signal")
				svc.shutdown()
			}
		}
	}
}

func (svc *service) start() {
	svc.servers = rcv.BuildServers()
	for _, srv := range svc.servers {
		go svc.serve(srv)
	}
	log.Info("service", "Started")
}

func (svc *service) shutdown() {
	for _, srv := range svc.servers {
		srv.Listener.Close()
	}
	log.Info("service", "Stopped")
	os.Exit(0)
}

func (svc *service) serve(srv rcv.Server) {
	for {
		conn, err := srv.Listener.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				if !strings.Contains(fmt.Sprintf("%v", err), "use of closed network connection") {
					// todo: retries
					time.Sleep(u.Timeout(cnf.DefaultConnectRetryWait))
					continue
				}
			}
			return
		}
		log.Debug("service", fmt.Sprintf("Accept: %s://%s", srv.Listener.Addr().Network(), srv.Listener.Addr().String()))

		if isModeHttp(srv.Cnf.Mode) {
			go fwd.ForwardHttp(srv.Cnf, srv.L4Proto, conn)
		} else {
			go fwd.Forward(srv.Cnf, srv.L4Proto, conn)
		}
	}
}

func isModeHttp(mode meta.ListenMode) bool {
	return mode == meta.ListenModeHttp || mode == meta.ListenModeHttps
}
