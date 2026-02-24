package service

import "errors"

var (
	ErrDepartmentNotFound   = errors.New("department not found")
	ErrEmployeeNotFound     = errors.New("employee not found")
	ErrDuplicateName        = errors.New("department with this name already exists under the same parent")
	ErrSelfParent           = errors.New("department cannot be its own parent")
	ErrCycle                = errors.New("would create cycle in department tree")
	ErrValidation           = errors.New("validation error")
	ErrReassignIDRequired   = errors.New("reassign_to_department_id is required when mode=reassign")
)
