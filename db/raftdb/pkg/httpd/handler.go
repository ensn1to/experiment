package httpd

import "net/http"

func (s *Service) keyRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	}
}
