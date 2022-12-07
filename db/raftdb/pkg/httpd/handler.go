package httpd

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/raft"
)

func (s *Service) keyRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.get(r, w)
	case "POST":
		s.post(r, w)
	case "DELETE":
		s.delete(r, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (s *Service) getKey(r *http.Request, w http.ResponseWriter) string {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		return ""
	}

	return parts[2]
}

func (s *Service) get(r *http.Request, w http.ResponseWriter) {
	k := s.getKey(r, w)
	if k == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lvl, err := level(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v, err := s.store.Get(k, lvl)
	if err != nil {
		// 容错：在去leader读一遍
		if errors.Is(err, raft.ErrNotLeader) {
			leader := s.store.LeaderAPIAddr()
			if leader == "" {
				http.Error(w, err.Error(), http.StatusTemporaryRedirect)
				return
			}

			redirect := s.FormRedirect(r, leader)
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			return
		}
	}

	b, err := json.Marshal(map[string]string{k: v})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := io.WriteString(w, string(b)); err != nil {
		s.logger.Printf("faile to writestring: %s", err.Error())
	}
}

func (s *Service) post(r *http.Request, w http.ResponseWriter) {
	// Read the value from the POST body.
	m := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for k, v := range m {
		if err := s.store.Set(k, v); err != nil {
			if errors.Is(err, raft.ErrNotLeader) {
				leader := s.store.LeaderAPIAddr()
				if leader == "" {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}

				redirect := s.FormRedirect(r, leader)
				http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *Service) delete(r *http.Request, w http.ResponseWriter) {
	k := s.getKey(r, w)
	if k == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.store.Delete(k); err != nil {
		if errors.Is(err, raft.ErrNotLeader) {
			leader := s.store.LeaderAPIAddr()
			if leader == "" {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}

			redirect := s.FormRedirect(r, leader)
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// s.store.Delete(k)
}

func (s *Service) joinRequestHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(m) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	httpAddr, ok := m["httpAddr"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	raftAddr, ok := m["raftAddr"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nodeID, ok := m["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.store.Join(nodeID, httpAddr, raftAddr); err != nil {
		if err == raft.ErrNotLeader {
			leader := s.store.LeaderAPIAddr()
			if leader == "" {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}

			// redirect to the real leader addr
			redirect := s.FormRedirect(r, leader)
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
