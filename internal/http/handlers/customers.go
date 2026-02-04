package handlers

import (
	"Book-Store/internal/authentication"
	"Book-Store/internal/http/middleware"
	"Book-Store/internal/models"
	"Book-Store/internal/response"
	"Book-Store/internal/store"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CustomerHandler struct {
	Store store.CustomerStore
	Cfg   *middleware.ApiConfig
}

func (h *CustomerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	path = strings.TrimSpace(path)
	pathParts := strings.Split(path, "/")

	if len(pathParts) == 2 && pathParts[1] == "register" {
		if r.Method == http.MethodPost {
			h.registerCustomer(w, r)
			return
		}
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if len(pathParts) == 2 && pathParts[1] == "login" {
		if r.Method == http.MethodPost {
			h.loginCustomer(w, r)
			return
		}
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var id int
	var hasID bool

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
		h.createCustomer(w, r)
	case http.MethodGet:
		if hasID {
			h.getCustomer(w, r, id)
		} else {
			h.listCustomers(w, r)
		}
	case http.MethodPut:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing customer ID")
			return
		}
		h.updateCustomer(w, r, id)
	case http.MethodDelete:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing customer ID")
			return
		}
		h.deleteCustomer(w, r, id)
	default:
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CustomerHandler) registerCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if customer.Email == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Email is required!")
		return
	}

	if customer.Password == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Password is required!")
		return
	}

	if customer.Name == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Name is required!")
		return
	}

	customers, err := h.Store.ListCustomers(ctx)
	if err == nil {
		for _, c := range customers {
			if c.Email == customer.Email {
				response.RespondWithError(w, http.StatusBadRequest, "Email already exists")
				return
			}
		}
	}

	hashed_password, err := authentication.HashPassword(customer.Password)
	if err != nil {
		log.Println("Got error to create password Hash")
		response.RespondWithError(w, http.StatusInternalServerError, "Could not create customer")
		return
	}

	customer.Password = hashed_password
	createdCustomer, err := h.Store.CreateCustomer(ctx, customer)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := authentication.MakeJWT(
		createdCustomer.ID,
		h.Cfg.Token,
		time.Hour,
	)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Error creating Token")
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, struct {
		ID        int            `json:"id"`
		Name      string         `json:"name"`
		Email     string         `json:"email"`
		Token     string         `json:"token"`
		Address   models.Address `json:"address"`
		CreatedAt time.Time      `json:"created_at"`
	}{
		ID:        createdCustomer.ID,
		Name:      createdCustomer.Name,
		Email:     createdCustomer.Email,
		Token:     token,
		Address:   createdCustomer.Address,
		CreatedAt: createdCustomer.CreatedAt,
	})
}

func (h *CustomerHandler) loginCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if loginRequest.Email == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Email is required!")
		return
	}

	if loginRequest.Password == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Password is required!")
		return
	}

	customers, err := h.Store.ListCustomers(ctx)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var foundCustomer *models.Customer
	for _, c := range customers {
		if c.Email == loginRequest.Email {
			foundCustomer = &c
			break
		}
	}

	if foundCustomer == nil {
		response.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	passwordMatch, err := authentication.CheckPasswordHash(loginRequest.Password, foundCustomer.Password)
	if err != nil || !passwordMatch {
		response.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := authentication.MakeJWT(
		foundCustomer.ID,
		h.Cfg.Token,
		time.Hour,
	)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Error creating Token")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, struct {
		ID        int            `json:"id"`
		Name      string         `json:"name"`
		Email     string         `json:"email"`
		Token     string         `json:"token"`
		Address   models.Address `json:"address"`
		CreatedAt time.Time      `json:"created_at"`
	}{
		ID:        foundCustomer.ID,
		Name:      foundCustomer.Name,
		Email:     foundCustomer.Email,
		Token:     token,
		Address:   foundCustomer.Address,
		CreatedAt: foundCustomer.CreatedAt,
	})
}

func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if customer.Email == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Email is required!")
		return
	}

	if customer.Password == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Password is required!")
		return
	}

	hashed_password, err := authentication.HashPassword(customer.Password)
	if err != nil {
		log.Println("Got error to create password Hash")
		response.RespondWithError(w, http.StatusInternalServerError, "Could not create customer")
		return
	}

	customer.Password = hashed_password
	createdCustomer, err := h.Store.CreateCustomer(ctx, customer)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := authentication.MakeJWT(
		createdCustomer.ID,
		h.Cfg.Token,
		time.Hour,
	)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Error creating Token")
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, struct {
		ID        int            `json:"id"`
		Name      string         `json:"name"`
		Email     string         `json:"email"`
		Token     string         `json:"token"`
		Address   models.Address `json:"address"`
		CreatedAt time.Time      `json:"created_at"`
	}{
		ID:        createdCustomer.ID,
		Name:      createdCustomer.Name,
		Email:     createdCustomer.Email,
		Token:     token,
		Address:   createdCustomer.Address,
		CreatedAt: createdCustomer.CreatedAt,
	})
}

func (h *CustomerHandler) getCustomer(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()

	customer, err := h.Store.GetCustomer(ctx, id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) listCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	customers, err := h.Store.ListCustomers(ctx)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) updateCustomer(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	defer r.Body.Close()

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updatedCustomer, err := h.Store.UpdateCustomer(ctx, id, customer)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, updatedCustomer)
}

func (h *CustomerHandler) deleteCustomer(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()

	err := h.Store.DeleteCustomer(ctx, id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Customer deleted successfully"})
}
