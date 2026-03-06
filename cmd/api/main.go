package main

import (
	"context"
	"log"
	"net/http"
	"os"
	httpapi "pulseforge/internal/http"
	"pulseforge/internal/repo"
	"pulseforge/internal/service"
)

func main(){
	pool, err := repo.NewPool(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil { log.Fatal(err) }
	defer pool.Close()
	userRepo := repo.NewUserRepo(pool)
	postRepo := repo.NewPostRepo(pool)
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)
	jwtSecret := os.Getenv("JWT_SECRET")

	addr := ":8080"
	log.Print("Listening on ", addr)

	server := &http.Server{
		Addr: addr,
		Handler: httpapi.NewMux(userService, postService, jwtSecret),
	}
	log.Fatal(server.ListenAndServe())
}
