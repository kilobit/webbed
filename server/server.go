/* Copyright 2021 Kilobit Labs Inc. */

package server

import _ "fmt"
import _ "errors"
import "context"
import "log"
import "time"
import "io/ioutil"
import "net"
import _ "net/url"
import "net/http"

type HTTPListener struct {
	net.Listener
}

func NewHTTPListener(ctx context.Context, addr string) (*HTTPListener, error) {

	config := &net.ListenConfig{
		Control:   nil,
		KeepAlive: 30 * time.Second,
	}

	ln, err := config.Listen(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	return &HTTPListener{ln}, nil
}

type Server struct {
	ctx   context.Context
	log   *log.Logger
	Error error
	http.Server
}

type ServerOpt func(*Server)

func ServerOptLogger(lg *log.Logger) ServerOpt {
	return ServerOpt(func(srv *Server) {
		srv.log = lg
	})
}

func NewHTTPServer(ctx context.Context, addr string, handler http.Handler, opts ...ServerOpt) *Server {

	srv := &Server{
		ctx,
		log.New(ioutil.Discard, "", 0),
		nil,
		http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 19, // 500 KB
		},
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (srv Server) ListenAndServe() error {

	ln, err := NewHTTPListener(srv.ctx, srv.Addr)
	if err != nil {
		return err
	}

	return srv.Serve(ln)
}

func (srv Server) ListenAndServeTLS(certFile, keyFile string) error {

	ln, err := NewHTTPListener(srv.ctx, srv.Addr)
	if err != nil {
		return err
	}

	return srv.ServeTLS(ln, certFile, keyFile)
}

// Starts the server in a goroutine and returns a wait channel to the
// caller.
//
func (srv Server) Start() <-chan bool {

	done := make(chan bool)

	srv.log.Printf("Starting server on %s.", srv.Addr)

	go func() {
		srv.Error = srv.ListenAndServe()

		srv.log.Printf("Server on %s has ended.", srv.Addr)

		done <- true
	}()

	return done
}

