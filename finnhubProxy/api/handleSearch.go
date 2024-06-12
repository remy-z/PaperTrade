package api

import "net/http"

func (s *APIServer) dispatchSearch(w http.ResponseWriter, r *http.Request) error {
	return s.dispatch(s.getSearch)(w, r)
}

func (s *APIServer) getSearch(w http.ResponseWriter, r *http.Request) error {
	symbol, err := getSymbol(r)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, symbol)
}
