package server

import (
	"context"
	"net/http"
)

type Logger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

type HTTPServer struct {
	Server *http.Server
	done   chan struct{}
	log    Logger
}

func NewServer(addr string, mainHandler http.Handler, logger Logger) *HTTPServer {
	mux := http.NewServeMux()

	// main service logic (app layer)
	mux.Handle("/", mainHandler)

	// auxiliary service logic
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return &HTTPServer{
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		done: make(chan struct{}),
		log:  logger,
	}
}

func (s *HTTPServer) Start() {
	go func() {
		defer close(s.done)
		s.log.Info("HTTP server starting ...")

		err := s.Server.ListenAndServe()
		if err != nil {
			s.log.Info("HTTP server stopped: ", err)
		}
	}()
	<-s.done
}

func (s *HTTPServer) Stop(ctx context.Context) {
	err := s.Server.Shutdown(ctx)
	if err != nil {
		s.log.Error("HTTP server shutdown error: ", err)
		return
	}
	s.done <- struct{}{}

	s.log.Info("HTTP server stops")
}
