/* Copyright 2020 Kilobit Labs Inc. */

package informed

import _ "fmt"
import _ "errors"
import "context"
import "time"
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
	ctx context.Context
	http.Server
}

func NewHTTPServer(ctx context.Context, handler http.Handler) (*Server, error) {

	srv := http.Server{
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 19, // 500 KB
	}

	return &Server{ctx, srv}, nil
}

func (srv Server) ListenAndServe() error {

	ln, err := NewHTTPListener(srv.ctx, ":0")
	if err != nil {
		return err
	}

	return srv.Serve(ln)
}

func (srv Server) ListenAndServeTLS(certFile, keyFile string) error {

	ln, err := NewHTTPListener(srv.ctx, ":0")
	if err != nil {
		return err
	}

	return srv.ServeTLS(ln, certFile, keyFile)
}
