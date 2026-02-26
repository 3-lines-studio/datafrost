package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/3-lines-studio/datafrost/internal/adapters"
	"github.com/3-lines-studio/datafrost/internal/db"
	"github.com/3-lines-studio/datafrost/internal/models"

	"github.com/go-chi/chi/v5"
)

type ConnectionsHandler struct {
	store   *db.ConnectionStore
	factory *adapters.Factory
}

func NewConnectionsHandler(store *db.ConnectionStore) *ConnectionsHandler {
	return &ConnectionsHandler{
		store:   store,
		factory: adapters.NewFactory(),
	}
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

func (h *ConnectionsHandler) ListAdapters(w http.ResponseWriter, r *http.Request) {
	adapters := h.factory.ListAdapters()
	JSONResponse(w, http.StatusOK, adapters)
}

func (h *ConnectionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		JSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	if req.Type == "" {
		JSONError(w, http.StatusBadRequest, "type is required")
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

func (h *ConnectionsHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req models.UpdateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		JSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	if req.Type == "" {
		JSONError(w, http.StatusBadRequest, "type is required")
		return
	}

	conn, err := h.store.Update(id, req)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, conn)
}

func (h *ConnectionsHandler) Test(w http.ResponseWriter, r *http.Request) {
	var req models.TestConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Type == "" {
		JSONError(w, http.StatusBadRequest, "type is required")
		return
	}

	if err := h.factory.TestConnection(req.Type, req.Credentials); err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *ConnectionsHandler) TestExisting(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	conn, err := h.store.GetByID(id)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if conn == nil {
		JSONError(w, http.StatusNotFound, "connection not found")
		return
	}

	if err := h.factory.TestConnection(conn.Type, conn.Credentials); err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]bool{"success": true})
}
