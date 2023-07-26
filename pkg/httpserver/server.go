package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	notify chan error
}

func New(handler http.Handler, address string) *Server {
	httpserver := http.Server{
		Handler: handler,
		Addr:    address,
	}

	s := &Server{
		server: &httpserver,
		notify: make(chan error, 1),
	}

	s.start()

	return s
}

func (s Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.server.Shutdown(ctx)
}
