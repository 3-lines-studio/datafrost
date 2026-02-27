package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type TabsHandler struct {
	uc *usecase.AppStateUsecase
}

func NewTabsHandler(uc *usecase.AppStateUsecase) *TabsHandler {
	return &TabsHandler{uc: uc}
}

func (h *TabsHandler) Get(w http.ResponseWriter, r *http.Request) {
	connectionIDStr := chi.URLParam(r, "id")
	connectionID, err := strconv.Atoi(connectionIDStr)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	tabs, err := h.uc.GetTabs(connectionID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to get tabs")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]any{"tabs": tabs})
}

func (h *TabsHandler) Save(w http.ResponseWriter, r *http.Request) {
	connectionIDStr := chi.URLParam(r, "id")
	connectionID, err := strconv.Atoi(connectionIDStr)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	var req struct {
		Tabs []entity.Tab `json:"tabs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.uc.SaveTabs(connectionID, req.Tabs); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to save tabs")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]any{"tabs": req.Tabs})
}
