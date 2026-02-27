package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type ConnectionsHandler struct {
	uc *usecase.ConnectionUsecase
}

func NewConnectionsHandler(uc *usecase.ConnectionUsecase) *ConnectionsHandler {
	return &ConnectionsHandler{uc: uc}
}

func (h *ConnectionsHandler) List(w http.ResponseWriter, r *http.Request) {
	connections, lastID, err := h.uc.List()
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]any{
		"connections": connections,
		"last_id":     lastID,
	}

	JSONResponse(w, http.StatusOK, response)
}

func (h *ConnectionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.CreateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	conn, err := h.uc.Create(req)
	if err != nil {
		if err == usecase.ErrNameRequired || err == usecase.ErrTypeRequired {
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
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

	if err := h.uc.Delete(id); err != nil {
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

	if err := h.uc.SetLastConnected(id); err != nil {
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

	var req entity.UpdateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	conn, err := h.uc.Update(id, req)
	if err != nil {
		if err == usecase.ErrNameRequired || err == usecase.ErrTypeRequired {
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, conn)
}

func (h *ConnectionsHandler) Test(w http.ResponseWriter, r *http.Request) {
	var req entity.TestConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.uc.Test(req); err != nil {
		if err == usecase.ErrTypeRequired {
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
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

	if err := h.uc.TestExisting(id); err != nil {
		if err == usecase.ErrConnectionNotFound {
			JSONError(w, http.StatusNotFound, "connection not found")
			return
		}
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]bool{"success": true})
}
