package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TabsHandler struct {
	db *sql.DB
}

func NewTabsHandler(db *sql.DB) *TabsHandler {
	return &TabsHandler{db: db}
}

type Tab struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	ConnectionID int    `json:"connectionId"`
	TableName    string `json:"tableName,omitempty"`
	Query        string `json:"query,omitempty"`
}

func (h *TabsHandler) Get(w http.ResponseWriter, r *http.Request) {
	connectionIDStr := chi.URLParam(r, "id")
	connectionID, err := strconv.Atoi(connectionIDStr)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	key := "tabs_" + strconv.Itoa(connectionID)
	var value string
	err = h.db.QueryRow("SELECT value FROM app_state WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		JSONResponse(w, http.StatusOK, map[string]interface{}{"tabs": []Tab{}})
		return
	} else if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to get tabs")
		return
	}

	var tabs []Tab
	if err := json.Unmarshal([]byte(value), &tabs); err != nil {
		JSONResponse(w, http.StatusOK, map[string]interface{}{"tabs": []Tab{}})
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{"tabs": tabs})
}

func (h *TabsHandler) Save(w http.ResponseWriter, r *http.Request) {
	connectionIDStr := chi.URLParam(r, "id")
	connectionID, err := strconv.Atoi(connectionIDStr)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	var req struct {
		Tabs []Tab `json:"tabs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tabsJSON, err := json.Marshal(req.Tabs)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to marshal tabs")
		return
	}

	key := "tabs_" + strconv.Itoa(connectionID)
	_, err = h.db.Exec(
		"INSERT INTO app_state (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		key, string(tabsJSON), string(tabsJSON),
	)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to save tabs")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{"tabs": req.Tabs})
}
