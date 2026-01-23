package main

import (
	"Book-Store/internal/config"
	"Book-Store/internal/http/handlers"
	"Book-Store/internal/http/middleware"
	"Book-Store/internal/http/router"
	"Book-Store/internal/reports"
	"Book-Store/internal/scheduler"
	"Book-Store/internal/store"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()
	memStore := store.NewMemStore()

	apiCfg := &middleware.ApiConfig{}

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

	reportStore := reports.NewReportStore(cfg.ReportOutputDirectory)
	reportHandler := &handlers.ReportHandler{
		OrderStore:  memStore,
		ReportStore: reportStore,
	}

	reportScheduler := scheduler.NewReportScheduler(memStore, reportStore, apiCfg, cfg.ReportInterval)
	reportScheduler.Start()

	metricsHandler := &handlers.MetricsHandler{
		BookStore:     memStore,
		AuthorStore:   memStore,
		CustomerStore: memStore,
		OrderStore:    memStore,
	}

	router.Router(
		apiCfg,
		bookHandler,
		authorHandler,
		customerHandler,
		orderHandler,
		reportHandler,
		metricsHandler,
		apiCfg,
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutdown signal received, stopping scheduler...")
		reportScheduler.Stop()
		os.Exit(0)
	}()

	fmt.Printf("Running Server on port :%s\n", cfg.ServerPort)
	fmt.Printf("Report generation interval: %v\n", cfg.ReportInterval)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
