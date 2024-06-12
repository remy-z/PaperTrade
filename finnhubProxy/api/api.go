package api

import (
	"encoding/json"
	"log"
	"net/http"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr    string
	finnhubClient *finnhub.DefaultApiService
}

func NewAPIServer(listenAddr string, finnhubClient *finnhub.DefaultApiService) *APIServer {
	return &APIServer{
		listenAddr:    listenAddr,
		finnhubClient: finnhubClient,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.Use(CORSMiddleware)
	router.HandleFunc("/quote/{symbol}", makeHttpHandleFunc(s.dispatchQuote))
	router.HandleFunc("/search/{symbol}", makeHttpHandleFunc(s.dispatchSearch))

	log.Println("JSON API server running on port: ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

// decorate our api func into a handler func so we can use with mux.HandleFunc
func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
