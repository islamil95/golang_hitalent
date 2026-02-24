package model

import (
	"time"
)

// Employee представляет сотрудника подразделения.
type Employee struct {
	ID           int        `json:"id" gorm:"primaryKey;autoIncrement"`
	DepartmentID int        `json:"department_id" gorm:"not null;index"`
	FullName     string     `json:"full_name" gorm:"type:varchar(200);not null"`
	Position     string     `json:"position" gorm:"type:varchar(200);not null"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	Department   *Department `json:"-" gorm:"foreignKey:DepartmentID"`
}

// TableName задаёт имя таблицы в базе.
func (Employee) TableName() string {
	return "employees"
}
