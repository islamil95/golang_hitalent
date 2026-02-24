package service

import (
	"context"
	"github.com/islamil95/golang_hitalent/internal/model"
	"github.com/islamil95/golang_hitalent/internal/repository"
)

// DepartmentService содержит бизнес-логику для подразделений.
type DepartmentService struct {
	depRepo *repository.DepartmentRepository
	empRepo *repository.EmployeeRepository
}

// NewDepartmentService создаёт новый сервис подразделений.
func NewDepartmentService(depRepo *repository.DepartmentRepository, empRepo *repository.EmployeeRepository) *DepartmentService {
	return &DepartmentService{depRepo: depRepo, empRepo: empRepo}
}

// Create создаёт подразделение.
func (s *DepartmentService) Create(ctx context.Context, in CreateDepartmentInput) (*model.Department, error) {
	name, err := ValidateDepartmentName(in.Name)
	if err != nil {
		return nil, err
	}
	if in.ParentID != nil {
		_, err := s.depRepo.GetByID(ctx, *in.ParentID)
		if err != nil {
			return nil, ErrDepartmentNotFound
		}
	}
	exists, err := s.depRepo.ExistsByNameAndParent(ctx, name, in.ParentID, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateName
	}
	d := &model.Department{Name: name, ParentID: in.ParentID}
	if err := s.depRepo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetByID возвращает подразделение с сотрудниками и поддеревом до указанной глубины.
func (s *DepartmentService) GetByID(ctx context.Context, id int, depth int, includeEmployees bool, employeeOrder string) (*DepartmentDetailResponse, error) {
	d, err := s.depRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDepartmentNotFound
	}
	res := &DepartmentDetailResponse{
		Department: toDepartmentDTO(d),
	}
	if includeEmployees {
		emps, _ := s.empRepo.ListByDepartmentID(ctx, id, employeeOrder)
		res.Employees = make([]EmployeeDTO, 0, len(emps))
		for i := range emps {
			res.Employees = append(res.Employees, toEmployeeDTO(&emps[i]))
		}
	}
	if depth > 0 {
		children, _ := s.getChildren(ctx, d.ID, depth-1, includeEmployees, employeeOrder)
		res.Children = children
	}
	return res, nil
}

func (s *DepartmentService) getChildren(ctx context.Context, parentID int, depth int, includeEmployees bool, employeeOrder string) ([]DepartmentDetailResponse, error) {
	deps, err := s.depRepo.GetChildrenByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	out := make([]DepartmentDetailResponse, 0, len(deps))
	for i := range deps {
		child, err := s.GetByID(ctx, deps[i].ID, depth, includeEmployees, employeeOrder)
		if err != nil {
			continue
		}
		out = append(out, *child)
	}
	return out, nil
}

// Update обновляет подразделение (name и/или parent_id).
func (s *DepartmentService) Update(ctx context.Context, id int, in UpdateDepartmentInput) (*model.Department, error) {
	d, err := s.depRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDepartmentNotFound
	}
	if in.Name != nil {
		name, err := ValidateDepartmentName(*in.Name)
		if err != nil {
			return nil, err
		}
		d.Name = name
	}
	// Check name uniqueness under current parent (after potential parent_id change)
	parentID := d.ParentID
	if in.ParentID != nil {
		parentID = in.ParentID
	}
	exists, err := s.depRepo.ExistsByNameAndParent(ctx, d.Name, parentID, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateName
	}
	if in.ParentID != nil {
		if *in.ParentID == id {
			return nil, ErrSelfParent
		}
		subtree, err := s.depRepo.GetSubtreeIDs(ctx, id)
		if err != nil {
			return nil, err
		}
		for _, sid := range subtree {
			if sid == *in.ParentID {
				return nil, ErrCycle
			}
		}
		_, err = s.depRepo.GetByID(ctx, *in.ParentID)
		if err != nil {
			return nil, ErrDepartmentNotFound
		}
		d.ParentID = in.ParentID
	}
	if err := s.depRepo.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

// Delete удаляет подразделение (режимы cascade или reassign).
func (s *DepartmentService) Delete(ctx context.Context, id int, mode string, reassignToID *int) error {
	_, err := s.depRepo.GetByID(ctx, id)
	if err != nil {
		return ErrDepartmentNotFound
	}
	switch mode {
	case "cascade":
		return s.depRepo.DeleteCascade(ctx, id)
	case "reassign":
		if reassignToID == nil {
			return ErrReassignIDRequired
		}
		_, err := s.depRepo.GetByID(ctx, *reassignToID)
		if err != nil {
			return ErrDepartmentNotFound
		}
		if *reassignToID == id {
			return ErrValidation
		}
		if err := s.depRepo.ReassignEmployees(ctx, id, *reassignToID); err != nil {
			return err
		}
		// Режим reassign: переносим сотрудников и привязку дочерних подразделений к новому родителю, затем удаляем исходное подразделение.
		return s.deleteReassign(ctx, id, *reassignToID)
	default:
		return ErrValidation
	}
}

func (s *DepartmentService) deleteReassign(ctx context.Context, id, reassignToID int) error {
	if err := s.depRepo.ReassignChildren(ctx, id, reassignToID); err != nil {
		return err
	}
	return s.depRepo.Delete(ctx, id)
}
