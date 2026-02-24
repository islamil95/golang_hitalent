package repository

import (
	"context"
	"github.com/islamil95/golang_hitalent/internal/model"
	"gorm.io/gorm"
)

// DepartmentRepository отвечает за доступ к данным подразделений.
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository создаёт новый репозиторий подразделений.
func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// Create сохраняет новое подразделение.
func (r *DepartmentRepository) Create(ctx context.Context, d *model.Department) error {
	return r.db.WithContext(ctx).Create(d).Error
}

// GetByID загружает подразделение по ID.
func (r *DepartmentRepository) GetByID(ctx context.Context, id int) (*model.Department, error) {
	var d model.Department
	err := r.db.WithContext(ctx).First(&d, id).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Update обновляет поля подразделения.
func (r *DepartmentRepository) Update(ctx context.Context, d *model.Department) error {
	return r.db.WithContext(ctx).Save(d).Error
}

// Delete удаляет подразделение по ID.
func (r *DepartmentRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Department{}, id).Error
}

// ExistsByNameAndParent проверяет, есть ли подразделение с тем же именем под тем же родителем.
func (r *DepartmentRepository) ExistsByNameAndParent(ctx context.Context, name string, parentID *int, excludeID int) (bool, error) {
	var count int64
	var err error
	if parentID == nil {
		q := r.db.WithContext(ctx).Model(&model.Department{}).Where("name = ? AND parent_id IS NULL", name)
		if excludeID > 0 {
			q = q.Where("id != ?", excludeID)
		}
		err = q.Count(&count).Error
	} else {
		q := r.db.WithContext(ctx).Model(&model.Department{}).Where("name = ? AND parent_id = ?", name, *parentID)
		if excludeID > 0 {
			q = q.Where("id != ?", excludeID)
		}
		err = q.Count(&count).Error
	}
	return count > 0, err
}

// GetChildrenByParentID возвращает прямых потомков подразделения.
func (r *DepartmentRepository) GetChildrenByParentID(ctx context.Context, parentID int) ([]model.Department, error) {
	var list []model.Department
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&list).Error
	return list, err
}

// GetSubtreeIDs возвращает все ID дочерних подразделений (для проверки циклов).
func (r *DepartmentRepository) GetSubtreeIDs(ctx context.Context, departmentID int) ([]int, error) {
	var ids []int
	seen := map[int]bool{departmentID: true}
	queue := []int{departmentID}
	for len(queue) > 0 {
		parent := queue[0]
		queue = queue[1:]
		var children []model.Department
		if err := r.db.WithContext(ctx).Where("parent_id = ?", parent).Find(&children).Error; err != nil {
			return nil, err
		}
		for _, c := range children {
			if !seen[c.ID] {
				seen[c.ID] = true
				ids = append(ids, c.ID)
				queue = append(queue, c.ID)
			}
		}
	}
	return ids, nil
}

// ReassignEmployees переносит всех сотрудников из одного подразделения в другое.
func (r *DepartmentRepository) ReassignEmployees(ctx context.Context, fromID, toID int) error {
	return r.db.WithContext(ctx).Model(&model.Employee{}).Where("department_id = ?", fromID).Update("department_id", toID).Error
}

// ReassignChildren обновляет parent_id у всех прямых потомков fromParentID на toParentID.
func (r *DepartmentRepository) ReassignChildren(ctx context.Context, fromParentID, toParentID int) error {
	return r.db.WithContext(ctx).Model(&model.Department{}).Where("parent_id = ?", fromParentID).Update("parent_id", toParentID).Error
}

// DeleteCascade удаляет подразделение и все его потомки (сотрудники удаляются каскадно по FK).
func (r *DepartmentRepository) DeleteCascade(ctx context.Context, id int) error {
	subtree, err := r.GetSubtreeIDs(ctx, id)
	if err != nil {
		return err
	}
	for i := len(subtree) - 1; i >= 0; i-- {
		if err := r.db.WithContext(ctx).Delete(&model.Department{}, subtree[i]).Error; err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Delete(&model.Department{}, id).Error
}
