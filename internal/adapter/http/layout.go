package http

import (
	"encoding/json"
	"net/http"

	"github.com/3-lines-studio/datafrost/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type LayoutHandler struct {
	uc *usecase.AppStateUsecase
}

func NewLayoutHandler(uc *usecase.AppStateUsecase) *LayoutHandler {
	return &LayoutHandler{uc: uc}
}

func (h *LayoutHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	layout, err := h.uc.GetLayout(key)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to get layout")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"layout": layout})
}

func (h *LayoutHandler) Save(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	var req struct {
		Layout string `json:"layout"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.uc.SaveLayout(key, req.Layout); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to save layout")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"layout": req.Layout})
}
