package handlers

import (
	"Book-Store/internal/reports"
	"Book-Store/internal/response"
	"Book-Store/internal/store"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReportHandler struct {
	OrderStore  store.OrderStore
	ReportStore *reports.ReportStore
}

func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	path = strings.TrimSpace(path)
	pathParts := strings.Split(path, "/")

	if len(pathParts) > 1 && pathParts[1] == "generate" {
		if r.Method == http.MethodPost {
			h.generateReport(w, r)
			return
		}
	}

	if len(pathParts) > 1 && pathParts[1] == "latest" {
		if r.Method == http.MethodGet {
			h.getLatestReport(w)
			return
		}
	}

	if r.Method == http.MethodGet {
		h.listReports(w, r)
		return
	}

	response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func (h *ReportHandler) generateReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	report, err := reports.GenerateSalesReport(ctx, h.OrderStore)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.ReportStore.SaveReport(report); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, report)
}

func (h *ReportHandler) getLatestReport(w http.ResponseWriter) {
	report, err := h.ReportStore.GetLatestReport()
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, report)
}

func (h *ReportHandler) listReports(w http.ResponseWriter, r *http.Request) {
	since := r.URL.Query().Get("since")

	if since != "" {
		daysAgo, err := strconv.Atoi(since)
		if err == nil {
			since := time.Now().AddDate(0, 0, -daysAgo)
			reports := h.ReportStore.GetReportsSince(since)
			response.RespondWithJSON(w, http.StatusOK, reports)
			return
		}
	}

	allReports := h.ReportStore.ListAllReports()
	response.RespondWithJSON(w, http.StatusOK, allReports)
}
