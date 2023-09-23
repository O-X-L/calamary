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

	"github.com/superstes/calamary/log"
	"github.com/superstes/calamary/proc/fwd"
	"github.com/superstes/calamary/rcv"
)

type service struct {
	listeners []rcv.Listener
}

func (svc *service) signalHandler() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigc
	log.Info("service", fmt.Sprintf("Signal received: %v", sig))
}

func (svc *service) start() (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithCancel(context.Background())
	svc.listeners = rcv.Start()
	for i := range svc.listeners {
		listener := svc.listeners[i]
		go svc.serve(listener)
	}
	log.Info("service", "Started")
	return
}

func (svc *service) shutdown(cancel context.CancelFunc) {
	cancel()
	for i := range svc.listeners {
		listener := svc.listeners[i]
		listener.Close()
	}
	log.Info("service", "Stopped")
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
