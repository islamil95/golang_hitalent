package service

import (
	"context"
	"github.com/islamil95/golang_hitalent/internal/model"
	"github.com/islamil95/golang_hitalent/internal/repository"
)

// EmployeeService содержит бизнес-логику для сотрудников.
type EmployeeService struct {
	depRepo *repository.DepartmentRepository
	empRepo *repository.EmployeeRepository
}

// NewEmployeeService создаёт новый сервис сотрудников.
func NewEmployeeService(depRepo *repository.DepartmentRepository, empRepo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{depRepo: depRepo, empRepo: empRepo}
}

// Create создаёт сотрудника в указанном подразделении.
func (s *EmployeeService) Create(ctx context.Context, departmentID int, in CreateEmployeeInput) (*model.Employee, error) {
	_, err := s.depRepo.GetByID(ctx, departmentID)
	if err != nil {
		return nil, ErrDepartmentNotFound
	}
	fullName, err := ValidateEmployeeFullName(in.FullName)
	if err != nil {
		return nil, err
	}
	position, err := ValidateEmployeePosition(in.Position)
	if err != nil {
		return nil, err
	}
	e := &model.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      in.HiredAt,
	}
	if err := s.empRepo.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}
