package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LayoutHandler struct {
	db *sql.DB
}

func NewLayoutHandler(db *sql.DB) *LayoutHandler {
	return &LayoutHandler{db: db}
}

func (h *LayoutHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	var value string
	err := h.db.QueryRow("SELECT value FROM app_state WHERE key = ?", "layout_"+key).Scan(&value)
	if err == sql.ErrNoRows {
		value = ""
	} else if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to get layout")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"layout": value})
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

	_, err := h.db.Exec(
		"INSERT INTO app_state (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		"layout_"+key, req.Layout, req.Layout,
	)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to save layout")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"layout": req.Layout})
}
