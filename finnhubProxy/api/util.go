package api

import (
	"net/http"
	"strings"

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

func findTermWithSubstring[V any](terms map[string]V, substring string) []string {
	var matchingTerms []string

	for key := range terms {
		if strings.Contains(key, substring) {
			matchingTerms = append(matchingTerms, key)
		}
	}
	return matchingTerms
}
