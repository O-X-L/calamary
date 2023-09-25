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

	"github.com/superstes/calamary/cnf/cnf_file"
	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/fwd"
	"github.com/superstes/calamary/rcv"
)

type service struct {
	listeners []rcv.Listener
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
				cnf_file.Load()

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
	svc.listeners = rcv.Start()
	for i := range svc.listeners {
		listener := svc.listeners[i]
		go svc.serve(listener)
	}
	log.Info("service", "Started")
}

func (svc *service) shutdown(cancel context.CancelFunc) {
	cancel()
	for i := range svc.listeners {
		listener := svc.listeners[i]
		listener.Close()
	}
	log.Info("service", "Stopped")
	os.Exit(0)
	/*
		ctx := context.Background()
		doneHTTP := httpserver.Shutdown(ctx)
		<-doneHTTP
	*/
}

func (svc *service) serve(ln rcv.Listener) error {
	// log.Info("service", fmt.Sprintf("Serving %s://%s", ln.Addr().Network(), ln.Addr().String()))
	var tempDelay time.Duration
	for {
		conn, err := ln.Accept()
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
		log.Debug("service", fmt.Sprintf("Accept: %s://%s", ln.Addr().Network(), ln.Addr().String()))

		go fwd.Forward(
			ln.Addr().Network(),
			conn,
		)
	}
}
