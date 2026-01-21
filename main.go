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

	memStore := store.NewMemStore()

	log.Printf("Loading database from: %s", cfg.DBPath)
	if err := memStore.LoadFromFile(cfg.DBPath); err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	bookHandler := &handlers.BookHandler{
		BookStore:   memStore,
		AuthorStore: memStore,
	}
	authorHandler := &handlers.AuthorHandler{Store: memStore}
	customerHandler := &handlers.CustomerHandler{Store: memStore}
	orderHandler := &handlers.OrderHandler{Store: memStore}

	router.Router(bookHandler, authorHandler, customerHandler, orderHandler)

	fmt.Printf("Running Server on port :%s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
