package handler

import (
	"net/http"
	"strconv"
	"strings"
)

// PathID извлекает целочисленный id из пути вида "/departments/123" или "/departments/123/employees".
func PathID(path, prefix string) (int, bool) {
	path = strings.TrimPrefix(path, prefix)
	path = strings.Trim(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return 0, false
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}
	return id, true
}

// QueryInt читает query-параметр как int или возвращает defaultVal при отсутствии/ошибке.
func QueryInt(r *http.Request, key string, defaultVal int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// QueryBool читает query-параметр как bool; defaultVal, если параметр не передан.
func QueryBool(r *http.Request, key string, defaultVal bool) bool {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultVal
	}
	return strings.EqualFold(s, "true") || s == "1"
}

// ClampDepth ограничивает глубину значением в диапазоне 1..5.
func ClampDepth(d int) int {
	if d < 1 {
		return 1
	}
	if d > 5 {
		return 5
	}
	return d
}
