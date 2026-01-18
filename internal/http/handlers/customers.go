package handlers

import (
	"Book-Store/internal/models"
	"Book-Store/internal/response"
	"Book-Store/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CustomerHandler struct {
	Store store.CustomerStore
}

func (h *CustomerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	var id int
	var hasID bool

	if len(pathParts) > 1 && pathParts[1] != "" {
		parsedID, err := strconv.Atoi(pathParts[1])
		if err == nil {
			id = parsedID
			hasID = true
		}
	}

	switch r.Method {
	case http.MethodPost:
		h.createCustomer(w, r)
	case http.MethodGet:
		if hasID {
			h.getCustomerByID(w, id)
		} else {
			h.listCustomers(w)
		}
	case http.MethodPut:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing Customer ID")
			return
		}
		h.updateCustomer(w, r, id)
	case http.MethodDelete:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing Customer ID")
			return
		}
		h.deleteCustomer(w, id)
	default:
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}

func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		fmt.Println("JSON Decode error in createCustomer")
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	created_customer, err := h.Store.CreateCustomer(customer)
	if err != nil {
		fmt.Printf("Failed to create Customer: %v\n", err)
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, created_customer)
}

func (h *CustomerHandler) getCustomerByID(w http.ResponseWriter, id int) {
	customer, err := h.Store.GetCustomer(id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) listCustomers(w http.ResponseWriter) {
	customers, err := h.Store.ListCustomers()
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) updateCustomer(w http.ResponseWriter, r *http.Request, id int) {
	defer r.Body.Close()

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updated_customer, err := h.Store.UpdateCustomer(id, customer)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, updated_customer)
}

func (h *CustomerHandler) deleteCustomer(w http.ResponseWriter, id int) {
	err := h.Store.DeleteCustomer(id)
	if err != nil {
		fmt.Printf("Got Error in deleteCustomer %v", err)
		response.RespondWithError(w, http.StatusInternalServerError, "Customer not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, "Customer deleted successfully")
}
