package repository

import (
	"context"
	"github.com/islamil95/golang_hitalent/internal/model"
	"gorm.io/gorm"
)

// EmployeeRepository отвечает за доступ к данным сотрудников.
type EmployeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository создаёт новый репозиторий сотрудников.
func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// Create сохраняет нового сотрудника.
func (r *EmployeeRepository) Create(ctx context.Context, e *model.Employee) error {
	return r.db.WithContext(ctx).Create(e).Error
}

// ListByDepartmentID возвращает сотрудников подразделения с сортировкой по created_at (или по full_name).
func (r *EmployeeRepository) ListByDepartmentID(ctx context.Context, departmentID int, orderBy string) ([]model.Employee, error) {
	var list []model.Employee
	q := r.db.WithContext(ctx).Where("department_id = ?", departmentID)
	switch orderBy {
	case "full_name":
		q = q.Order("full_name ASC")
	default:
		q = q.Order("created_at ASC")
	}
	err := q.Find(&list).Error
	return list, err
}
