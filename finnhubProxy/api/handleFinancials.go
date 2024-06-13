package api

import (
	"net/http"
	"strings"
)

func (s *APIServer) dispatchFinancials(w http.ResponseWriter, r *http.Request) error {
	return s.dispatch(s.getFinancials)(w, r)
}

func (s *APIServer) getFinancials(w http.ResponseWriter, r *http.Request) error {
	symbol, err := getSymbol(r)
	if err != nil {
		return err
	}
	symbol = strings.ToUpper(symbol)
	res, err := s.finnhubFinancial(symbol)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, res.Metric)
}
