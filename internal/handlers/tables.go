package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/3-lines-studio/datafrost/internal/adapters"
	"github.com/3-lines-studio/datafrost/internal/db"
	"github.com/3-lines-studio/datafrost/internal/models"

	"github.com/go-chi/chi/v5"
)

type TablesHandler struct {
	store *db.ConnectionStore
	cache *adapters.AdapterCache
}

func NewTablesHandler(store *db.ConnectionStore, cache *adapters.AdapterCache) *TablesHandler {
	return &TablesHandler{
		store: store,
		cache: cache,
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

	adapter, err := h.cache.Get(conn.ID, conn.Type, conn.Credentials)
	if err != nil {
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

func (h *TablesHandler) Tree(w http.ResponseWriter, r *http.Request) {
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

	adapter, err := h.cache.Get(conn.ID, conn.Type, conn.Credentials)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "failed to connect: "+err.Error())
		return
	}

	if treeLister, ok := adapter.(adapters.TreeLister); ok {
		nodes, err := treeLister.ListTree()
		if err != nil {
			JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		JSONResponse(w, http.StatusOK, nodes)
		return
	}

	// Fallback: wrap flat table list into a single schema node for compatibility.
	tables, err := adapter.ListTables()
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	children := make([]models.TreeNode, 0, len(tables))
	for _, t := range tables {
		children = append(children, models.TreeNode{Name: t.Name, Type: t.Type, FullName: t.FullName})
	}

	JSONResponse(w, http.StatusOK, []models.TreeNode{
		{
			Name:     "tables",
			Type:     "schema",
			FullName: "tables",
			Children: children,
		},
	})
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

	adapter, err := h.cache.Get(conn.ID, conn.Type, conn.Credentials)
	if err != nil {
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
	store *db.ConnectionStore
	cache *adapters.AdapterCache
}

func NewQueryHandler(store *db.ConnectionStore, cache *adapters.AdapterCache) *QueryHandler {
	return &QueryHandler{
		store: store,
		cache: cache,
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

	adapter, err := h.cache.Get(conn.ID, conn.Type, conn.Credentials)
	if err != nil {
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
