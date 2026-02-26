package handlers

import (
	"datafrost/internal/adapters"
	"datafrost/internal/db"
	"datafrost/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TablesHandler struct {
	store   *db.ConnectionStore
	factory *adapters.Factory
}

func NewTablesHandler(store *db.ConnectionStore) *TablesHandler {
	return &TablesHandler{
		store:   store,
		factory: adapters.NewFactory(),
	}
}

func (h *TablesHandler) List(w http.ResponseWriter, r *http.Request) {
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

	adapter, err := h.factory.GetAdapter(conn.Type)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "failed to get adapter: "+err.Error())
		return
	}
	defer adapter.Close()

	if err := adapter.Connect(conn.Credentials); err != nil {
		JSONError(w, http.StatusBadRequest, "failed to connect: "+err.Error())
		return
	}

	tables, err := adapter.ListTables()
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, tables)
}

func (h *TablesHandler) GetData(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	tableName := chi.URLParam(r, "name")
	if tableName == "" {
		JSONError(w, http.StatusBadRequest, "table name is required")
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

	adapter, err := h.factory.GetAdapter(conn.Type)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "failed to get adapter: "+err.Error())
		return
	}
	defer adapter.Close()

	if err := adapter.Connect(conn.Credentials); err != nil {
		JSONError(w, http.StatusBadRequest, "failed to connect: "+err.Error())
		return
	}

	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	limit := 25
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	page := 1
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	var filters []models.Filter
	filtersStr := r.URL.Query().Get("filters")
	if filtersStr != "" {
		if err := json.Unmarshal([]byte(filtersStr), &filters); err != nil {
			JSONError(w, http.StatusBadRequest, "invalid filters parameter")
			return
		}
	}

	offset := (page - 1) * limit

	result, err := adapter.GetTableData(tableName, limit, offset, filters)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, result)
}

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryHandler struct {
	store   *db.ConnectionStore
	factory *adapters.Factory
}

func NewQueryHandler(store *db.ConnectionStore) *QueryHandler {
	return &QueryHandler{
		store:   store,
		factory: adapters.NewFactory(),
	}
}

func (h *QueryHandler) Execute(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Query == "" {
		JSONError(w, http.StatusBadRequest, "query is required")
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

	adapter, err := h.factory.GetAdapter(conn.Type)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "failed to get adapter: "+err.Error())
		return
	}
	defer adapter.Close()

	if err := adapter.Connect(conn.Credentials); err != nil {
		JSONError(w, http.StatusBadRequest, "failed to connect: "+err.Error())
		return
	}

	result, err := adapter.ExecuteQuery(req.Query)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, result)
}
