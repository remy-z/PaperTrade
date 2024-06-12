package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) dispatch(get apiFunc) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case "GET":
			return get(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return nil
		}
	}
}

func getSymbol(r *http.Request) (string, error) {
	return mux.Vars(r)["symbol"], nil
}
