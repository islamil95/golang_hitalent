package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/islamil95/golang_hitalent/internal/service"
)

// EmployeeHandler обрабатывает HTTP-запросы по сотрудникам внутри подразделений.
type EmployeeHandler struct {
	svc *service.EmployeeService
	log *slog.Logger
}

// NewEmployeeHandler создаёт новый обработчик сотрудников.
func NewEmployeeHandler(svc *service.EmployeeService, log *slog.Logger) *EmployeeHandler {
	return &EmployeeHandler{svc: svc, log: log}
}

// Create обрабатывает POST /departments/{id}/employees/.
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Err(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	departmentID, ok := PathID(r.URL.Path, "/departments/")
	if !ok || departmentID <= 0 {
		Err(w, http.StatusBadRequest, "invalid department id")
		return
	}
	// Ожидается путь вида /departments/{id}/employees (допускается завершающий слэш).
	path := strings.TrimPrefix(r.URL.Path, "/departments/")
	path = strings.Trim(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 || (parts[1] != "employees" && !strings.HasPrefix(parts[1], "employees/")) {
		Err(w, http.StatusNotFound, "not found")
		return
	}
	var in service.CreateEmployeeInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		Err(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	emp, err := h.svc.Create(r.Context(), departmentID, in)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	JSON(w, http.StatusCreated, emp)
}

func (h *EmployeeHandler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case err == service.ErrDepartmentNotFound:
		Err(w, http.StatusNotFound, "department not found")
	case err == service.ErrValidation:
		Err(w, http.StatusBadRequest, "validation error")
	default:
		h.log.Error("service error", "err", err)
		Err(w, http.StatusInternalServerError, "internal server error")
	}
}
