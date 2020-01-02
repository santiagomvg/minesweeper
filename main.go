package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods("POST").Path("/action").Handler(newEndpoint(handleGameAction))
	r.Methods("POST").Path("/newGame").Handler(newEndpoint(handleRestartAction))
	//	r.NotFoundHandler = notFoundHandler{}
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
	log.Println("ListenAndServe:8080")
}

type endpointFunc func(w http.ResponseWriter, r *http.Request) error

type ServiceAPIHandler struct {
	serve endpointFunc
}

func (sah ServiceAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h := w.Header()
	h.Set("Content-Type", "application/json; charset=UTF-8")
	h.Set("Cache-Control", "no-store")
	h.Set("Pragma", "no-cache")

	//panic handler
	defer func() {
		if r := recover(); r != nil {
			handleError(w, r)
		}
	}()
	if err := sah.serve(w, r); err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err interface{}) {

	w.WriteHeader(500)

	type servicePanicError struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}
	lErr := json.NewEncoder(w).Encode(&servicePanicError{
		Status: 500,
		Error:  fmt.Sprintf("fatal error - %v", err),
	})

	if lErr != nil {
		log.Fatal(lErr)
	}

}

func newEndpoint(f endpointFunc) http.Handler {
	return ServiceAPIHandler{
		serve: f,
	}
}

type notFoundHandler struct{}

func (h notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
