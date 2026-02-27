package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/3-lines-studio/datafrost/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type SavedQueriesHandler struct {
	uc *usecase.SavedQueryUsecase
}

func NewSavedQueriesHandler(uc *usecase.SavedQueryUsecase) *SavedQueriesHandler {
	return &SavedQueriesHandler{uc: uc}
}

func (h *SavedQueriesHandler) List(w http.ResponseWriter, r *http.Request) {
	connectionIDStr := chi.URLParam(r, "id")
	connectionID, err := strconv.ParseInt(connectionIDStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	queries, err := h.uc.ListByConnection(connectionID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to fetch saved queries")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]any{"queries": queries})
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

	query, err := h.uc.Create(connectionID, req.Name, req.Query)
	if err != nil {
		if err == usecase.ErrNameRequired || err == usecase.ErrQueryRequired {
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
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

	query, err := h.uc.Update(queryID, req.Name, req.Query)
	if err != nil {
		if err == usecase.ErrQueryNotFound {
			JSONError(w, http.StatusNotFound, "Query not found")
			return
		}
		if err == usecase.ErrNameRequired || err == usecase.ErrQueryRequired {
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to update saved query")
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

	if err := h.uc.Delete(queryID); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to delete saved query")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
