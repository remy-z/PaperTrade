package api

import (
	"net/http"
	"strings"
)

func (s *APIServer) dispatchQuote(w http.ResponseWriter, r *http.Request) error {
	return s.dispatch(s.getQuote)(w, r)
}

func (s *APIServer) getQuote(w http.ResponseWriter, r *http.Request) error {
	symbol, err := getSymbol(r)
	if err != nil {
		return err
	}
	symbol = strings.ToUpper(symbol)
	res, err := s.finnhubQuote(symbol)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, res)
}
