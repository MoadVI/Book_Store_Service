package handlers

import (
	"Book-Store/internal/http/middleware"
	"Book-Store/internal/models"
	"Book-Store/internal/response"
	"Book-Store/internal/store"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type OrderHandler struct {
	Store store.OrderStore
}

func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	path = strings.TrimSpace(path)
	pathParts := strings.Split(path, "/")

	var (
		id    int
		hasID bool
	)

	if len(pathParts) > 1 && pathParts[1] != "" {
		idStr := strings.TrimSpace(pathParts[1])
		parsedID, err := strconv.Atoi(idStr)
		if err == nil {
			id = parsedID
			hasID = true
		}
	}

	switch r.Method {
	case http.MethodPost:
		h.createOrder(w, r)
	case http.MethodGet:
		if hasID {
			h.getOrderByID(w, r, id)
		} else if status := r.URL.Query().Get("status"); status != "" {
			h.searchOrderByStatus(w, r)
		} else if r.URL.Query().Get("start_date") != "" && r.URL.Query().Get("end_date") != "" {
			h.getOrdersInTimeRange(w, r)
		} else {
			h.listOrders(w, r)
		}
	case http.MethodPut:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Order ID required")
			return
		}
		h.changeStatus(w, r, id)
	default:
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	userID := middleware.GetUserIDFromContext(ctx)

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	order.Customer.ID = userID

	resultChan := make(chan models.Order, 1)
	errChan := make(chan error, 1)

	go func() {
		createdOrder, err := h.Store.CreateOrder(ctx, order)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- createdOrder
	}()

	select {
	case <-ctx.Done():
		response.RespondWithError(w, http.StatusRequestTimeout, "Request cancelled")
		return
	case err := <-errChan:
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	case createdOrder := <-resultChan:
		response.RespondWithJSON(w, http.StatusCreated, createdOrder)
	}
}

func (h *OrderHandler) getOrderByID(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()

	order, err := h.Store.GetOrder(ctx, id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}
	response.RespondWithJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) changeStatus(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	status := r.URL.Query().Get("status")
	if status == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Status query parameter required")
		return
	}

	switch status {
	case "completed":
		_, err := h.Store.CompleteOrder(ctx, id)
		if err != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Order completed"})
	case "cancelled":
		_, err := h.Store.CancelOrder(ctx, id)
		if err != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Order cancelled"})
	default:
		response.RespondWithError(w, http.StatusBadRequest, "Invalid status")
	}
}

func (h *OrderHandler) searchOrderByStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status := r.URL.Query().Get("status")

	orders, err := h.Store.SearchOrderByStatus(ctx, status)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.RespondWithJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) listOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.Store.ListOrders(ctx)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.RespondWithJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) getOrdersInTimeRange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid start_date format. 2024-01-01T00:00:00Z")
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid end_date format.2024-12-31T23:59:59Z")
		return
	}

	if startDate.After(endDate) {
		response.RespondWithError(w, http.StatusBadRequest, "start_date must be before end_date")
		return
	}

	orders, err := h.Store.GetOrdersInTimeRange(ctx, startDate, endDate)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, orders)
}
