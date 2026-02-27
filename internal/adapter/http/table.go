package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type TablesHandler struct {
	uc *usecase.TableUsecase
}

func NewTablesHandler(uc *usecase.TableUsecase) *TablesHandler {
	return &TablesHandler{uc: uc}
}

func (h *TablesHandler) List(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	tables, err := h.uc.ListTables(id)
	if err != nil {
		if err == usecase.ErrConnectionNotFound {
			JSONError(w, http.StatusNotFound, "connection not found")
			return
		}
		JSONError(w, http.StatusBadRequest, "failed to connect: "+err.Error())
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

	var filters []entity.Filter
	filtersStr := r.URL.Query().Get("filters")
	if filtersStr != "" {
		if err := json.Unmarshal([]byte(filtersStr), &filters); err != nil {
			JSONError(w, http.StatusBadRequest, "invalid filters parameter")
			return
		}
	}

	offset := (page - 1) * limit

	result, err := h.uc.GetTableData(id, tableName, limit, offset, filters)
	if err != nil {
		if err == usecase.ErrConnectionNotFound {
			JSONError(w, http.StatusNotFound, "connection not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, result)
}

func (h *TablesHandler) GetSchema(w http.ResponseWriter, r *http.Request) {
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

	schema, err := h.uc.GetTableSchema(id, tableName)
	if err != nil {
		if err == usecase.ErrConnectionNotFound {
			JSONError(w, http.StatusNotFound, "connection not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, schema)
}

type QueryHandler struct {
	uc *usecase.QueryUsecase
}

func NewQueryHandler(uc *usecase.QueryUsecase) *QueryHandler {
	return &QueryHandler{uc: uc}
}

func (h *QueryHandler) Execute(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req entity.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.uc.Execute(id, req.Query)
	if err != nil {
		if err == usecase.ErrConnectionNotFound {
			JSONError(w, http.StatusNotFound, "connection not found")
			return
		}
		if err == usecase.ErrQueryRequired {
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, result)
}
