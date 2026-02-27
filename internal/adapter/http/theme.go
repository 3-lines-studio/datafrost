package http

import (
	"encoding/json"
	"net/http"

	"github.com/3-lines-studio/datafrost/internal/usecase"
)

type ThemeHandler struct {
	uc *usecase.AppStateUsecase
}

func NewThemeHandler(uc *usecase.AppStateUsecase) *ThemeHandler {
	return &ThemeHandler{uc: uc}
}

func (h *ThemeHandler) Get(w http.ResponseWriter, r *http.Request) {
	theme, err := h.uc.GetTheme()
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to get theme")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"theme": theme})
}

func (h *ThemeHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Theme string `json:"theme"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Theme != "light" && req.Theme != "dark" {
		JSONError(w, http.StatusBadRequest, "Invalid theme value")
		return
	}

	if err := h.uc.SetTheme(req.Theme); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to save theme")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"theme": req.Theme})
}
