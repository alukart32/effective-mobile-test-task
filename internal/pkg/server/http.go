package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env/v7"
)

// A httpServer defines wrapper for the HTTP httpServer.
type httpServer struct {
	srv             *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// Http returns a new HTTP server. It runs after creation.
func Http(cfg HttpConfig, h http.Handler) (*httpServer, error) {
	if len(cfg.ADDR) == 0 || cfg.Empty() {
		opts := env.Options{RequiredIfNoDef: true}
		if err := env.Parse(&cfg, opts); err != nil {
			return nil, fmt.Errorf("read config: %v", err)
		}
	}

	srv := &http.Server{
		Addr:         cfg.ADDR,
		Handler:      h,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	s := httpServer{
		srv:             srv,
		notify:          make(chan error, 1),
		shutdownTimeout: cfg.ShutdownTimeout,
	}

	go func() {
		defer close(s.notify)
		s.notify <- s.srv.ListenAndServe()
	}()

	return &s, nil
}

// Notify throws a server error.
func (s *httpServer) Notify() <-chan error {
	return s.notify
}

// Shutdown gracefully stops the server during timeout.
func (s *httpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

// HttpConfig represents the http server configuration.
type HttpConfig struct {
	//Port specifies the listening port of the server.
	ADDR string `env:"HTTP_SERVER_ADDRESS" envDefault:":8080"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	ReadTimeout time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"500ms"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. A zero or negative value means
	// there will be no timeout.
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"1s"`

	// ShutdownTimeout is the maximum duration before timing out
	// stops the running server. A zero or negative value means
	// there will be no timeout.
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"1s"`
}

// Empty checks on being empty.
func (c HttpConfig) Empty() bool {
	return len(c.ADDR) == 0 &&
		c.ReadTimeout == 0 &&
		c.WriteTimeout == 0 &&
		c.ShutdownTimeout == 0
}
