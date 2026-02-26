package handlers

import (
	"datafrost/internal/db"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SavedQueriesHandler struct {
	store *db.SavedQueriesStore
}

func NewSavedQueriesHandler(store *db.SavedQueriesStore) *SavedQueriesHandler {
	return &SavedQueriesHandler{store: store}
}

func (h *SavedQueriesHandler) List(w http.ResponseWriter, r *http.Request) {
	connectionIDStr := chi.URLParam(r, "id")
	connectionID, err := strconv.ParseInt(connectionIDStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	queries, err := h.store.ListByConnection(connectionID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to fetch saved queries")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{"queries": queries})
}

func (h *SavedQueriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	connectionIDStr := chi.URLParam(r, "id")
	connectionID, err := strconv.ParseInt(connectionIDStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Query == "" {
		JSONError(w, http.StatusBadRequest, "Name and query are required")
		return
	}

	query, err := h.store.Create(connectionID, req.Name, req.Query)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to create saved query")
		return
	}

	JSONResponse(w, http.StatusCreated, query)
}

func (h *SavedQueriesHandler) Update(w http.ResponseWriter, r *http.Request) {
	queryIDStr := chi.URLParam(r, "queryId")
	queryID, err := strconv.ParseInt(queryIDStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid query ID")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Query == "" {
		JSONError(w, http.StatusBadRequest, "Name and query are required")
		return
	}

	query, err := h.store.Update(queryID, req.Name, req.Query)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to update saved query")
		return
	}

	if query == nil {
		JSONError(w, http.StatusNotFound, "Query not found")
		return
	}

	JSONResponse(w, http.StatusOK, query)
}

func (h *SavedQueriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	queryIDStr := chi.URLParam(r, "queryId")
	queryID, err := strconv.ParseInt(queryIDStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid query ID")
		return
	}

	if err := h.store.Delete(queryID); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to delete saved query")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
