package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter,r *http.Request){
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // w writes back, r receives the request
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok": true,
		})
	})
	addr := ":8080"
	log.Printf("Listening on ", addr)

	server := &http.Server{
		Addr: addr,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}