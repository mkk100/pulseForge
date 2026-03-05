package httpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"pulseforge/internal/service"
)

type createUserReq struct {
	UserName string `json:"userName"`
}

type createPostReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
}

func NewMux(userService *service.UserService, postService *service.PostService) *http.ServeMux {
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
		id, err := userService.GetUserIDByName(r.Context(), name)
		if err != nil {
			log.Printf("user lookup failed: %v", err)
			http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
			return
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

		userID, err := userService.CreateUser(r.Context(), req.UserName)
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
		switch r.Method {
		case http.MethodPost:
			var req createPostReq
			decoded := json.NewDecoder(r.Body)
			decoded.DisallowUnknownFields()
			if err := decoded.Decode(&req); err != nil {
				http.Error(w, "invalid json body", http.StatusBadRequest)
				return
			}

			postID, err := postService.CreatePost(r.Context(), service.CreatePostInput{
				Title:       req.Title,
				Description: req.Description,
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

		case http.MethodGet:
			limit := 10
			if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
				parsedLimit, err := strconv.Atoi(limitParam)
				if err != nil || parsedLimit <= 0 {
					http.Error(w, "invalid limit", http.StatusBadRequest)
					return
				}
				limit = parsedLimit
			}

			posts, err := postService.ListRecentPosts(r.Context(), limit)
			if err != nil {
				log.Printf("post retrieval failed: %v", err)
				http.Error(w, "failed to retrieve posts", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"posts": posts,
			})

		default:
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		}
	})

	return mux
}
