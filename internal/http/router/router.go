package router

import (
	"Book-Store/internal/http/handlers"
	"Book-Store/internal/http/middleware"
	"net/http"
)

func Router(
	apiCfg *middleware.ApiConfig,
	bookHandler *handlers.BookHandler,
	authorHandler *handlers.AuthorHandler,
	customerHandler *handlers.CustomerHandler,
	orderHandler *handlers.OrderHandler,
	reportHandler *handlers.ReportHandler,
	metricsHandler *handlers.MetricsHandler,
	hitsHandler *middleware.ApiConfig,
) {
	http.Handle("/books/", apiCfg.MiddlewareMetricsInc(bookHandler))
	http.Handle("/authors/", apiCfg.MiddlewareMetricsInc(authorHandler))

	http.Handle("/customers", apiCfg.MiddlewareMetricsInc(customerHandler))
	http.Handle("/customers/", middleware.AuthMiddleware(apiCfg.Token,
		apiCfg.MiddlewareMetricsInc(customerHandler)))

	http.Handle("/orders", apiCfg.MiddlewareMetricsInc(orderHandler))
	http.Handle("/orders/", middleware.AuthMiddleware(apiCfg.Token,
		apiCfg.MiddlewareMetricsInc(orderHandler)))

	http.Handle("/reports/sales", reportHandler)

	http.Handle("/metrics", metricsHandler)

	http.Handle("/metrics/hits", hitsHandler)
}
