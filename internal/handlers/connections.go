package handlers

import (
	"datafrost/internal/db"
	"datafrost/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ConnectionsHandler struct {
	store *db.ConnectionStore
}

func NewConnectionsHandler(store *db.ConnectionStore) *ConnectionsHandler {
	return &ConnectionsHandler{store: store}
}

func (h *ConnectionsHandler) List(w http.ResponseWriter, r *http.Request) {
	connections, err := h.store.List()
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	lastID, _ := h.store.GetLastConnected()

	response := map[string]any{
		"connections": connections,
		"last_id":     lastID,
	}

	JSONResponse(w, http.StatusOK, response)
}

func (h *ConnectionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.URL == "" {
		JSONError(w, http.StatusBadRequest, "name and url are required")
		return
	}

	conn, err := h.store.Create(req)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, conn)
}

func (h *ConnectionsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.store.Delete(id); err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ConnectionsHandler) SetLastConnected(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.store.SetLastConnected(id); err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
