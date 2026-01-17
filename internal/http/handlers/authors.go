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

type AuthorHandler struct {
	Store store.AuthorStore
}

func (h *AuthorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		h.createAuthor(w, r)
	case http.MethodGet:
		if hasID {
			h.getAuthorByID(w, id)
		} else {
			h.listAuthors(w)
		}
	case http.MethodPut:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing Author ID")
			return
		}
		h.updateAuthor(w, r, id)
	case http.MethodDelete:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing Author ID")
			return
		}
		h.deleteAuthor(w, id)
	default:
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}

func (h *AuthorHandler) createAuthor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		fmt.Println("JSON Decode Error in createAuthor")
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	created_author, err := h.Store.CreateAuthor(author)
	if err != nil {
		fmt.Printf("Fialed tp create Author: %v\n", err)
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to save Author")
		return
	}
	response.RespondWithJSON(w, http.StatusCreated, created_author)
}

func (h *AuthorHandler) getAuthorByID(w http.ResponseWriter, id int) {
	author, lookUpError := h.Store.GetAuthor(id)
	if lookUpError != nil {
		response.RespondWithError(w, http.StatusNotFound, "Author does not exist")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, author)
}

func (h *AuthorHandler) listAuthors(w http.ResponseWriter) {
	authors, err := h.Store.ListAuthors()
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, authors)
}

func (h *AuthorHandler) updateAuthor(w http.ResponseWriter, r *http.Request, id int) {
	defer r.Body.Close()
	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updated_author, err := h.Store.UpdateAuthor(id, author)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, updated_author)
}

func (h *AuthorHandler) deleteAuthor(w http.ResponseWriter, id int) {
	delete_err := h.Store.DeleteAuthor(id)
	if delete_err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	msg := fmt.Sprintf("Author with ID: %d deleted successfully", id)
	response.RespondWithJSON(w, http.StatusOK, msg)
}

