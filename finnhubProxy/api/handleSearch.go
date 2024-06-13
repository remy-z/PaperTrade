package api

import (
	"fmt"
	"net/http"
)

func (s *APIServer) dispatchSearch(w http.ResponseWriter, r *http.Request) error {
	return s.dispatch(s.getSearch)(w, r)
}

// write a json object of format [{symbol: description,},] to a search get request
// with all terms that match the search term in ticker or description
func (s *APIServer) getSearch(w http.ResponseWriter, r *http.Request) error {
	term, err := getSymbol(r)
	if err != nil {
		return err
	}
	symbols := findTermWithSubstring(s.companyInfo, term)
	descriptions := findTermWithSubstring(s.descriptionToSymbol, term)
	result := make(map[string]string)
	fmt.Println()
	for _, symbol := range symbols {
		result[symbol] = s.companyInfo[symbol].Description
	}
	for _, d := range descriptions {
		result[s.descriptionToSymbol[d]] = d
	}
	return writeJSON(w, http.StatusOK, result)
}
