package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

type Metrics interface {
	GetHitsByGenre(genre string) int64
	GetAllHits() map[string]int64
	ResetHits()
}

type ApiConfig struct {
	booksHits     atomic.Int64
	authorsHits   atomic.Int64
	customersHits atomic.Int64
	ordersHits    atomic.Int64
	Token         string
}

func (h *ApiConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")

	if genre != "" {
		fmt.Fprintf(w, "%s hits: %d", genre, h.GetHitsByGenre(genre))
		return
	}

	for k, v := range h.GetAllHits() {
		fmt.Fprintf(w, "%s: %d\n", k, v)
	}
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case strings.HasPrefix(path, "/books"):
			cfg.booksHits.Add(1)
		case strings.HasPrefix(path, "/authors"):
			cfg.authorsHits.Add(1)
		case strings.HasPrefix(path, "/customers"):
			cfg.customersHits.Add(1)
		case strings.HasPrefix(path, "/orders"):
			cfg.ordersHits.Add(1)
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) GetHitsByGenre(genre string) int64 {
	switch strings.ToLower(genre) {
	case "books":
		return cfg.booksHits.Load()
	case "authors":
		return cfg.authorsHits.Load()
	case "customers":
		return cfg.customersHits.Load()
	case "orders":
		return cfg.ordersHits.Load()
	default:
		return 0
	}
}

func (cfg *ApiConfig) GetAllHits() map[string]int64 {
	return map[string]int64{
		"books":     cfg.booksHits.Load(),
		"authors":   cfg.authorsHits.Load(),
		"customers": cfg.customersHits.Load(),
		"orders":    cfg.ordersHits.Load(),
	}
}

func (cfg *ApiConfig) ResetHits() {
	cfg.booksHits.Store(0)
	cfg.authorsHits.Store(0)
	cfg.customersHits.Store(0)
	cfg.ordersHits.Store(0)
}
