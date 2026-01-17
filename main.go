package main

import (
	"Book-Store/internal/config"
	"Book-Store/internal/http/handlers"
	"Book-Store/internal/http/router"
	"Book-Store/internal/store"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	bookStore := store.NewMemStore()

	log.Printf("Loading database from: %s", cfg.DBPath)
	if err := bookStore.LoadFromFile(cfg.DBPath); err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	bookHandler := &handlers.BookHandler{Store: bookStore}

	router.Router(bookHandler)

	fmt.Printf("Running Server on port :%s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}

