package api

import (
	"encoding/json"
	"log"
	"net/http"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/gorilla/mux"
	"github.com/remy-z/PaperTrade/finnhubProxy/common"
)

type APIServer struct {
	listenAddr          string
	finnhubClient       *finnhub.DefaultApiService
	tree                *common.TST
	descriptionToSymbol map[string]string
	companyInfo         map[string]*CompanyProfile
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func NewAPIServer(listenAddr string, finnhubClient *finnhub.DefaultApiService) *APIServer {
	return &APIServer{
		listenAddr:          listenAddr,
		finnhubClient:       finnhubClient,
		tree:                &common.TST{},
		descriptionToSymbol: make(map[string]string),
		companyInfo:         make(map[string]*CompanyProfile),
	}
}

func (s *APIServer) Run() {
	//load finnhubData into memory
	s.initFinnhubData()

	router := mux.NewRouter()

	router.Use(CORSMiddleware)
	router.HandleFunc("/quote/{symbol}", makeHttpHandleFunc(s.dispatchQuote))
	router.HandleFunc("/search/{symbol}", makeHttpHandleFunc(s.dispatchSearch))
	router.HandleFunc("/financials/{symbol}", makeHttpHandleFunc(s.dispatchFinancials))

	log.Println("JSON API server running on port: ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// decorate our api func into a handler func so we can use with mux.HandleFunc
func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
