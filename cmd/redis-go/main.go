package main

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/db"
	"github.com/codecrafters-io/redis-starter-go/internal/server"
)

func main() {
	log.SetFlags(0)

	store := db.New()
	handler := commands.New(store)
	srv := server.New(handler)

	if err := srv.ListenAndServe("0.0.0.0:6379"); err != nil {
		log.Fatal(err)
	}
}
