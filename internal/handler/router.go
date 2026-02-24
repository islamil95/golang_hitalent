package handler

import (
	"net/http"
	"strings"
)

// Router маршрутизирует запросы к обработчикам подразделений и сотрудников.
type Router struct {
	dep *DepartmentHandler
	emp *EmployeeHandler
}

// NewRouter создаёт новый роутер.
func NewRouter(dep *DepartmentHandler, emp *EmployeeHandler) *Router {
	return &Router{dep: dep, emp: emp}
}

// ServeHTTP реализует интерфейс http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	segments := strings.Split(path, "/")

	switch {
	case len(segments) == 1 && segments[0] == "departments" && req.Method == http.MethodPost:
		r.dep.Create(w, req)
	case len(segments) == 2 && segments[0] == "departments" && segments[1] != "":
		// /departments/{id} — GET, PATCH, DELETE.
		switch req.Method {
		case http.MethodGet:
			r.dep.GetByID(w, req)
		case http.MethodPatch:
			r.dep.Update(w, req)
		case http.MethodDelete:
			r.dep.Delete(w, req)
		default:
			Err(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	case len(segments) >= 3 && segments[0] == "departments" && segments[2] == "employees" && req.Method == http.MethodPost:
		// /departments/{id}/employees.
		r.emp.Create(w, req)
	default:
		Err(w, http.StatusNotFound, "not found")
	}
}
