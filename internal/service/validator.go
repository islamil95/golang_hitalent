package service

import (
	"strings"
	"unicode/utf8"
)

const maxLen = 200

// ValidateDepartmentName триммит и валидирует имя подразделения (1..200 символов).
func ValidateDepartmentName(name string) (string, error) {
	s := strings.TrimSpace(name)
	if s == "" {
		return "", ErrValidation
	}
	if utf8.RuneCountInString(s) > maxLen {
		return "", ErrValidation
	}
	return s, nil
}

// ValidateEmployeeFullName валидирует full_name сотрудника (1..200 символов).
func ValidateEmployeeFullName(name string) (string, error) {
	s := strings.TrimSpace(name)
	if s == "" {
		return "", ErrValidation
	}
	if utf8.RuneCountInString(s) > maxLen {
		return "", ErrValidation
	}
	return s, nil
}

// ValidateEmployeePosition валидирует должность (1..200 символов).
func ValidateEmployeePosition(pos string) (string, error) {
	s := strings.TrimSpace(pos)
	if s == "" {
		return "", ErrValidation
	}
	if utf8.RuneCountInString(s) > maxLen {
		return "", ErrValidation
	}
	return s, nil
}
