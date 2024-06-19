package api

import (
	"net/http"
	"strings"
)

func (s *APIServer) dispatchProfile(w http.ResponseWriter, r *http.Request) error {
	return s.dispatch(s.getProfile)(w, r)
}

func (s *APIServer) getProfile(w http.ResponseWriter, r *http.Request) error {
	symbol, err := getSymbol(r)
	if err != nil {
		return err
	}
	symbol = strings.ToUpper(symbol)
	res, err := s.finnhubProfile(symbol)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, res)
}
