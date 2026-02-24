package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/islamil95/golang_hitalent/internal/service"
)

// DepartmentHandler обрабатывает HTTP-запросы по подразделениям.
type DepartmentHandler struct {
	svc *service.DepartmentService
	log *slog.Logger
}

// NewDepartmentHandler создаёт новый обработчик подразделений.
func NewDepartmentHandler(svc *service.DepartmentService, log *slog.Logger) *DepartmentHandler {
	return &DepartmentHandler{svc: svc, log: log}
}

// Create обрабатывает POST /departments/.
func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Err(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var in service.CreateDepartmentInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		Err(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	dep, err := h.svc.Create(r.Context(), in)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	JSON(w, http.StatusCreated, dep)
}

// GetByID обрабатывает GET /departments/{id}.
func (h *DepartmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Err(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id, ok := PathID(r.URL.Path, "/departments/")
	if !ok || id <= 0 {
		Err(w, http.StatusBadRequest, "invalid department id")
		return
	}
	depth := ClampDepth(QueryInt(r, "depth", 1))
	includeEmployees := QueryBool(r, "include_employees", true)
	orderBy := "created_at"
	if QueryBool(r, "sort_employees_by_name", false) {
		orderBy = "full_name"
	}
	res, err := h.svc.GetByID(r.Context(), id, depth, includeEmployees, orderBy)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	JSON(w, http.StatusOK, res)
}

// Update обрабатывает PATCH /departments/{id}.
func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		Err(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id, ok := PathID(r.URL.Path, "/departments/")
	if !ok || id <= 0 {
		Err(w, http.StatusBadRequest, "invalid department id")
		return
	}
	var in service.UpdateDepartmentInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		Err(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	dep, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	JSON(w, http.StatusOK, dep)
}

// Delete обрабатывает DELETE /departments/{id}.
func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		Err(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id, ok := PathID(r.URL.Path, "/departments/")
	if !ok || id <= 0 {
		Err(w, http.StatusBadRequest, "invalid department id")
		return
	}
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "cascade"
	}
	var reassignToID *int
	if mode == "reassign" {
		sid := r.URL.Query().Get("reassign_to_department_id")
		if sid == "" {
			Err(w, http.StatusBadRequest, "reassign_to_department_id is required when mode=reassign")
			return
		}
		rid, err := parseInt(sid)
		if err != nil {
			Err(w, http.StatusBadRequest, "invalid reassign_to_department_id")
			return
		}
		reassignToID = &rid
	}
	err := h.svc.Delete(r.Context(), id, mode, reassignToID)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *DepartmentHandler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case err == service.ErrDepartmentNotFound:
		Err(w, http.StatusNotFound, "department not found")
	case err == service.ErrDuplicateName:
		Err(w, http.StatusConflict, "department with this name already exists under the same parent")
	case err == service.ErrSelfParent:
		Err(w, http.StatusBadRequest, "department cannot be its own parent")
	case err == service.ErrCycle:
		Err(w, http.StatusConflict, "would create cycle in department tree")
	case err == service.ErrReassignIDRequired:
		Err(w, http.StatusBadRequest, "reassign_to_department_id is required when mode=reassign")
	case err == service.ErrValidation:
		Err(w, http.StatusBadRequest, "validation error")
	default:
		h.log.Error("service error", "err", err)
		Err(w, http.StatusInternalServerError, "internal server error")
	}
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
