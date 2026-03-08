package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	httpapi "pulseforge/internal/http"
	"pulseforge/internal/repo"
	"pulseforge/internal/service"
)

type createUserResponse struct {
	UserID int64 `json:"userId"`
}

type createPostResponse struct {
	PostID int64 `json:"postId"`
}

type loginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"userId"`
}

type listPostsResponse struct {
	Posts []service.Post `json:"posts"`
}

func TestCreateAndListPostsIntegration(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL must be set for integration tests")
	}

	ctx := context.Background()
	pool, err := repo.NewPool(ctx, dsn)
	if err != nil {
		t.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, `TRUNCATE posts, users RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("reset db: %v", err)
	}

	userRepo := repo.NewUserRepo(pool)
	postRepo := repo.NewPostRepo(pool)
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)
	server := httptest.NewServer(httpapi.NewMux(userService, postService, "test-secret"))
	defer server.Close()

	userBody := bytes.NewBufferString(`{"userName":"thomas"}`)
	userResp, err := http.Post(server.URL+"/users", "application/json", userBody)
	if err != nil {
		t.Fatalf("create user request: %v", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusCreated {
		t.Fatalf("create user status = %d, want %d", userResp.StatusCode, http.StatusCreated)
	}

	var createdUser createUserResponse
	if err := json.NewDecoder(userResp.Body).Decode(&createdUser); err != nil {
		t.Fatalf("decode create user response: %v", err)
	}
	if createdUser.UserID == 0 {
		t.Fatal("create user returned userId = 0")
	}

	loginBody := bytes.NewBufferString(`{"userName":"thomas"}`)
	loginResp, err := http.Post(server.URL+"/login", "application/json", loginBody)
	if err != nil {
		t.Fatalf("login request: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginResp.StatusCode, http.StatusOK)
	}

	var loggedInUser loginResponse
	if err := json.NewDecoder(loginResp.Body).Decode(&loggedInUser); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if loggedInUser.Token == "" {
		t.Fatal("login returned empty token")
	}

	postPayload, err := json.Marshal(map[string]any{
		"title":       "hello",
		"description": "world",
	})
	if err != nil {
		t.Fatalf("marshal post payload: %v", err)
	}

	postReq, err := http.NewRequest(http.MethodPost, server.URL+"/posts", bytes.NewReader(postPayload))
	if err != nil {
		t.Fatalf("build create post request: %v", err)
	}
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Authorization", "Bearer "+loggedInUser.Token)

	postResp, err := http.DefaultClient.Do(postReq)
	if err != nil {
		t.Fatalf("create post request: %v", err)
	}
	defer postResp.Body.Close()

	if postResp.StatusCode != http.StatusCreated {
		t.Fatalf("create post status = %d, want %d", postResp.StatusCode, http.StatusCreated)
	}

	var createdPost createPostResponse
	if err := json.NewDecoder(postResp.Body).Decode(&createdPost); err != nil {
		t.Fatalf("decode create post response: %v", err)
	}
	if createdPost.PostID == 0 {
		t.Fatal("create post returned postId = 0")
	}

	listResp, err := http.Get(server.URL + "/posts?limit=10")
	if err != nil {
		t.Fatalf("list posts request: %v", err)
	}
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("list posts status = %d, want %d", listResp.StatusCode, http.StatusOK)
	}

	var listedPosts listPostsResponse
	if err := json.NewDecoder(listResp.Body).Decode(&listedPosts); err != nil {
		t.Fatalf("decode list posts response: %v", err)
	}
	if len(listedPosts.Posts) != 1 {
		t.Fatalf("list posts count = %d, want 1", len(listedPosts.Posts))
	}

	post := listedPosts.Posts[0]
	if post.ID != createdPost.PostID {
		t.Fatalf("post id = %d, want %d", post.ID, createdPost.PostID)
	}
	if post.Title != "hello" {
		t.Fatalf("post title = %q, want %q", post.Title, "hello")
	}
	if post.Description != "world" {
		t.Fatalf("post description = %q, want %q", post.Description, "world")
	}
	if post.UserID != createdUser.UserID {
		t.Fatalf("post userId = %d, want %d", post.UserID, createdUser.UserID)
	}
}