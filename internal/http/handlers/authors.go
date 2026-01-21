package handlers

import (
	"Book-Store/internal/models"
	"Book-Store/internal/response"
	"Book-Store/internal/store"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type AuthorHandler struct {
	Store store.AuthorStore
}

func (h *AuthorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	path = strings.TrimSpace(path)
	pathParts := strings.Split(path, "/")

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
		h.createAuthor(w, r)
	case http.MethodGet:
		if hasID {
			h.getAuthor(w, r, id)
		} else {
			h.listAuthors(w, r)
		}
	case http.MethodPut:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing author ID")
			return
		}
		h.updateAuthor(w, r, id)
	case http.MethodDelete:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing author ID")
			return
		}
		h.deleteAuthor(w, r, id)
	default:
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *AuthorHandler) createAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	createdAuthor, err := h.Store.CreateAuthor(ctx, author)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, createdAuthor)
}

func (h *AuthorHandler) getAuthor(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()

	author, err := h.Store.GetAuthor(ctx, id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, author)
}

func (h *AuthorHandler) listAuthors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authors, err := h.Store.ListAuthors(ctx)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, authors)
}

func (h *AuthorHandler) updateAuthor(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	defer r.Body.Close()

	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updatedAuthor, err := h.Store.UpdateAuthor(ctx, id, author)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, updatedAuthor)
}

func (h *AuthorHandler) deleteAuthor(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()

	err := h.Store.DeleteAuthor(ctx, id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, "Author deleted successfully")
}

