package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pulseforge/internal/db"
)

func main(){
	pool, err := db.NewPool(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil { log.Fatal(err) }
	defer pool.Close()
	userRepo := db.NewUserRepo(pool)

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

	mux.HandleFunc("/users/id", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HIT %s %s", r.Method, r.URL.Path)
		id, err := userRepo.GetUserIDByName(r.Context(), "alice")
		fmt.Print(id, err)
	})

	addr := ":8080"
	log.Print("Listening on ", addr)

	server := &http.Server{
		Addr: addr,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}