package http

import (
	"net/http"

	"github.com/3-lines-studio/datafrost/internal/usecase"
)

type AdapterHandler struct {
	uc *usecase.AdapterUsecase
}

func NewAdapterHandler(uc *usecase.AdapterUsecase) *AdapterHandler {
	return &AdapterHandler{uc: uc}
}

func (h *AdapterHandler) List(w http.ResponseWriter, r *http.Request) {
	adapters := h.uc.List()
	JSONResponse(w, http.StatusOK, adapters)
}
