package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type ThemeHandler struct {
	db *sql.DB
}

func NewThemeHandler(db *sql.DB) *ThemeHandler {
	return &ThemeHandler{db: db}
}

func (h *ThemeHandler) Get(w http.ResponseWriter, r *http.Request) {
	var value string
	err := h.db.QueryRow("SELECT value FROM app_state WHERE key = 'theme'").Scan(&value)
	if err == sql.ErrNoRows {
		value = "light"
	} else if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to get theme")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"theme": value})
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

	_, err := h.db.Exec(
		"INSERT INTO app_state (key, value) VALUES ('theme', ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		req.Theme, req.Theme,
	)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to save theme")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"theme": req.Theme})
}
