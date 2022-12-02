package httpd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/ensn1to/experiment/tree/master/db/raftdb/pkg/store"
)

type Store interface {
	// data operations
	Get(key string, lvl store.ConsistencyLevel) (string, error)

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

// Stop stop the service.
func (s *Service) Stop() {
	s.ln.Close()
}

// ServeHTTP simple router
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// key means operations on data resource
	if strings.HasPrefix(r.URL.Path, "/key") {
		// todo
	} else if r.URL.Path == "/join" {
		s.joinRequestHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// 301 redirect
func (s *Service) FormRedirect(r *http.Request, host string) string {
	protocol := "http"
	rq := r.URL.RawQuery
	if rq != "" {
		rq = fmt.Sprintf("?%s", rq)
	}

	return fmt.Sprintf("%s://%s%s%s", protocol, host, r.URL.Path, rq)
}

func level(req *http.Request) (store.ConsistencyLevel, error) {
	q := req.URL.Query()
	lvl := strings.TrimSpace(q.Get("level"))

	switch strings.ToLower(lvl) {
	case "default":
		return store.Default, nil
	case "stale":
		return store.Stale, nil
	case "consistent":
		return store.Consistent, nil
	default:
		return store.Default, nil
	}
}
