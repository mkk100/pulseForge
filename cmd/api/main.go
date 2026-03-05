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

type createUserReq struct {
    UserName string `json:"userName"`
}

type createPostReq struct {
	postID int64 `json:"postID"`
	postTitle string `json:"postTitle"`
	postDescription string `json:"postDescription"`
	userID int64 `json:"userID`
}

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
		id, err := userRepo.GetUserIDByName(r.Context(), "thomas")
		fmt.Print(id, err)
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req createUserReq
		decoded := json.NewDecoder(r.Body)
		decoded.DisallowUnknownFields()
		if err := decoded.Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
		}
		userId, err := userRepo.CreateUser(r.Context(), req.UserName)
		if err != nil {
			log.Printf("create user failed: %v", err)
    		http.Error(w, "failed to gycreate user", http.StatusInternalServerError)
    		return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"userId": userId})
	})

	mux.HandleFunc("/posts",func(w http.ResponseWriter, r *http.Request){
		var req createPostReq
		decoded := json.NewDecoder(r.Body)
		decoded.DisallowUnknownFields()
		if err := decoded.Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		fmt.Println(req)
	})

	addr := ":8080"
	log.Print("Listening on ", addr)

	server := &http.Server{
		Addr: addr,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}