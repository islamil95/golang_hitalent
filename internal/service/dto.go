package service

import (
	"time"

	"github.com/islamil95/golang_hitalent/internal/model"
)

// CreateDepartmentInput — тело запроса для POST /departments/.
type CreateDepartmentInput struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

// UpdateDepartmentInput — тело запроса для PATCH /departments/{id}.
type UpdateDepartmentInput struct {
	Name     *string `json:"name"`
	ParentID *int    `json:"parent_id"`
}

// CreateEmployeeInput — тело запроса для POST /departments/{id}/employees/.
type CreateEmployeeInput struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at"`
}

// DepartmentDTO — представление подразделения в ответе API.
type DepartmentDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

// EmployeeDTO — представление сотрудника в ответе API.
type EmployeeDTO struct {
	ID           int        `json:"id"`
	DepartmentID int        `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// DepartmentDetailResponse — ответ для GET /departments/{id}.
type DepartmentDetailResponse struct {
	Department DepartmentDTO              `json:"department"`
	Employees  []EmployeeDTO              `json:"employees,omitempty"`
	Children   []DepartmentDetailResponse `json:"children,omitempty"`
}

// toDepartmentDTO преобразует модель подразделения в DTO.
func toDepartmentDTO(d *model.Department) DepartmentDTO {
	return DepartmentDTO{
		ID:        d.ID,
		Name:      d.Name,
		ParentID:  d.ParentID,
		CreatedAt: d.CreatedAt,
	}
}

// toEmployeeDTO преобразует модель сотрудника в DTO.
func toEmployeeDTO(e *model.Employee) EmployeeDTO {
	return EmployeeDTO{
		ID:           e.ID,
		DepartmentID: e.DepartmentID,
		FullName:     e.FullName,
		Position:     e.Position,
		HiredAt:      e.HiredAt,
		CreatedAt:    e.CreatedAt,
	}
}