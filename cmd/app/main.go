package service

import (
	"net/http"
	"sync"
	"json"

	"https://github.com/gorilla/mux"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	r.HadleFunc("/subscription", handler).Methods("POST")
	r.HadleFunc("/subscription", handler).Methods("GET")
	r.HadleFunc("/subscription", handler).Methods("UPDATE")
	r.HadleFunc("/subscription", handler).Methods("DELETE")
}