package httpd

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/raft"
)

func (s *Service) keyRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	}
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
