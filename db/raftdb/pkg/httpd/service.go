package httpd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type Store interface {
	// data operations
	Get(key string) (string, error)

	Set(key string, value string) error

	Delete(key string) error

	// node operations
	Join(nodeID string, httpAddr string, addr string) error

	LeaderAPIAddr() string
}

// Service http service
type Service struct {
	ln    net.Listener
	store Store
	addr  string

	logger *log.Logger
}

func New(addr string, store Store) *Service {
	return &Service{
		addr:  addr,
		store: store,
		// todo: logger
	}
}

func (s *Service) Start() error {
	srv := http.Server{
		Handler: s,
	}
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.ln = ln

	http.Handle("/", s)

	go func() {
		err := srv.Serve(s.ln)
		if err != nil {
			// TODO: log error
			fmt.Printf("http server error: %s", err)
		}
	}()

	return nil
}

// ServeHTTP simple router
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// key means operations on data resource
	if strings.HasPrefix(r.URL.Path, "/key") {
		// todo
	} else if r.URL.Path == "/join" {
		// todo
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
