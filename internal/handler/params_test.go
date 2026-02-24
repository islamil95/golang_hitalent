package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathID(t *testing.T) {
	id, ok := PathID("/departments/42", "/departments/")
	assert.True(t, ok)
	assert.Equal(t, 42, id)

	id, ok = PathID("/departments/42/employees", "/departments/")
	assert.True(t, ok)
	assert.Equal(t, 42, id)

	_, ok = PathID("/departments/", "/departments/")
	assert.False(t, ok)

	_, ok = PathID("/departments/abc", "/departments/")
	assert.False(t, ok)
}

func TestClampDepth(t *testing.T) {
	assert.Equal(t, 1, ClampDepth(0))
	assert.Equal(t, 1, ClampDepth(1))
	assert.Equal(t, 3, ClampDepth(3))
	assert.Equal(t, 5, ClampDepth(5))
	assert.Equal(t, 5, ClampDepth(10))
}

func TestQueryInt(t *testing.T) {
	r := httptest.NewRequest("GET", "http://test/?depth=3", nil)
	assert.Equal(t, 3, QueryInt(r, "depth", 1))
	assert.Equal(t, 1, QueryInt(r, "missing", 1))
}

func TestQueryBool(t *testing.T) {
	r := httptest.NewRequest("GET", "http://test/?include_employees=true", nil)
	assert.True(t, QueryBool(r, "include_employees", false))
	assert.False(t, QueryBool(r, "missing", false))
}
