package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_NotFound(t *testing.T) {
	// Use nil handlers; router only checks path and method, so we need minimal deps.
	// For 404 we only need the router to route to "not found" path.
	dep := &DepartmentHandler{}
	emp := &EmployeeHandler{}
	router := NewRouter(dep, emp)

	req := httptest.NewRequest("GET", "http://test/unknown", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	dep := &DepartmentHandler{}
	emp := &EmployeeHandler{}
	router := NewRouter(dep, emp)

	req := httptest.NewRequest("PUT", "http://test/departments/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
