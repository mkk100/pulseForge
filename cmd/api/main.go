package main

import (
	"context"
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
	postRepo := db.NewPostRepo(pool)

	addr := ":8080"
	log.Print("Listening on ", addr)

	server := &http.Server{
		Addr: addr,
		Handler: newMux(userRepo, postRepo),
	}
	log.Fatal(server.ListenAndServe())
}
