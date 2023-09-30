package main

import (
	"context"
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
				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				svc.shutdown(cancel)
			}
		}
	}
}

func (svc *service) start() {
	svc.servers = rcv.BuildServers()
	for i := range svc.servers {
		listener := svc.servers[i]
		go svc.serve(listener)
	}
	log.Info("service", "Started")
}

func (svc *service) shutdown(cancel context.CancelFunc) {
	cancel()
	for i := range svc.servers {
		server := svc.servers[i]
		if u.IsIn(string(server.Cnf.Mode), []string{"http", "https"}) {
			server.HttpServer.Close()
			time.Sleep(time.Millisecond * 500)
		}
		server.Listener.Close()
	}
	log.Info("service", "Stopped")
	os.Exit(0)
	/*
		ctx := context.Background()
		doneHTTP := httpserver.Shutdown(ctx)
		<-doneHTTP
	*/
}

func (svc *service) serve(srv rcv.Server) (err error) {
	// log.Info("service", fmt.Sprintf("Serving %s://%s", ln.Addr().Network(), ln.Addr().String()))

	if srv.Cnf.Mode == meta.ListenModeHttp {
		log.Debug("service", "Starting HTTP server")
		err = srv.HttpServer.ListenAndServe()
		if err != nil {
			log.ErrorS("service", fmt.Sprintf("Failed to start HTTP server: %v", err))
			return err
		}

	} else if srv.Cnf.Mode == meta.ListenModeHttps {
		log.Debug("service", "Starting HTTPS server")
		err = srv.HttpServer.ListenAndServeTLS(
			cnf.C.Service.Certs.ServerPublic,
			cnf.C.Service.Certs.ServerPrivate,
		)
		if err != nil {
			log.ErrorS("service", fmt.Sprintf("Failed to start HTTPS server: %v", err))
			return err
		}

	} else {
		var tempDelay time.Duration

		for {
			conn, err := srv.Listener.Accept()
			if err != nil {
				if _, ok := err.(net.Error); ok {
					if !strings.Contains(fmt.Sprintf("%v", err), "use of closed network connection") {
						if tempDelay == 0 {
							tempDelay = 1 * time.Second
						} else {
							tempDelay *= 2
						}
						if max := 5 * time.Second; tempDelay > max {
							tempDelay = max
						}
						log.Warn("service", fmt.Sprintf("Error: %v, retrying in %v", err, tempDelay))
						time.Sleep(tempDelay)
						continue
					}
				}
				return err
			}
			tempDelay = 0
			log.Debug("service", fmt.Sprintf("Accept: %s://%s", srv.Listener.Addr().Network(), srv.Listener.Addr().String()))

			go fwd.Forward(srv.Cnf, srv.L4Proto, conn)
		}
	}
	return nil
}
