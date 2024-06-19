package api

import (
	"fmt"
	"net/http"

	"github.com/remy-z/PaperTrade/finnhubProxy/common"
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
	result := make([]map[string]string, 0)
	for _, symbol := range symbols {
		m := make(map[string]string)
		m["symbol"], m["description"] = symbol, s.companyInfo[symbol].Description
		result = append(result, m)
	}
	for _, d := range descriptions {
		m := make(map[string]string)
		m["symbol"], m["description"] = s.descriptionToSymbol[d], d
		result = append(result, m)
	}
	fmt.Println(result)
	common.HeapSort(result, s.wilshire)
	return writeJSON(w, http.StatusOK, result)
}

func (s *APIServer) searchScore() {

}

// if len(term) >= 2, display all tickers with matching prefix, then all tickers not matching prefix
func (s *APIServer) sortSearchResult() {

}
