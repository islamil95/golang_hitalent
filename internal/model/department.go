package model

import (
	"time"
)

// Department представляет подразделение в оргструктуре.
type Department struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"type:varchar(200);not null"`
	ParentID  *int       `json:"parent_id" gorm:"index"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	Parent    *Department `json:"-" gorm:"foreignKey:ParentID"`
	Children  []Department `json:"-" gorm:"foreignKey:ParentID"`
	Employees []Employee  `json:"-" gorm:"foreignKey:DepartmentID"`
}

// TableName задаёт имя таблицы в базе.
func (Department) TableName() string {
	return "departments"
}
