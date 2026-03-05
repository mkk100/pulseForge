package main

import (
	"encoding/json"
	"log"
	"net/http"

	"pulseforge/internal/db"
)

type createUserReq struct {
	UserName string `json:"userName"`
}

type createPostReq struct {
	PostTitle       string `json:"postTitle"`
	PostDescription string `json:"postDescription"`
	UserID          int64  `json:"userID"`
}

func newMux(userRepo *db.UserRepo, postRepo *db.PostRepo) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok": true,
		})
	})

	mux.HandleFunc("/users/id", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		id, err := userRepo.GetUserIDByName(r.Context(), name)
		if err != nil {
			log.Print("retrieval failed")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"userId":   id,
			"userName": name,
		})
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

		userID, err := userRepo.CreateUser(r.Context(), req.UserName)
		if err != nil {
			log.Printf("create user failed: %v", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"userId": userID})
	})

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		var req createPostReq
		decoded := json.NewDecoder(r.Body)
		decoded.DisallowUnknownFields()
		if err := decoded.Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		postID, err := postRepo.CreatePost(r.Context(), db.Post{
			Title:       req.PostTitle,
			Description: req.PostDescription,
			UserID:      req.UserID,
		})
		if err != nil {
			log.Printf("create post failed: %v", err)
			http.Error(w, "failed to create post", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"postId": postID})
	})

	mux.HandleFunc("/retrievePosts", func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		posts, err := postRepo.ListRecentPosts(r.Context(), limit)
		if err != nil {
			log.Print("post retrieval failed")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"posts": posts,
		})
	})

	return mux
}
