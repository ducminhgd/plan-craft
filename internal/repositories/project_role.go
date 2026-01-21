package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ducminhgd/plan-craft/internal"
	"github.com/ducminhgd/plan-craft/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProjectRoleRepository is the repository for project role entities
type ProjectRoleRepository struct {
	db *gorm.DB
}

// NewProjectRoleRepository creates a new project role repository
func NewProjectRoleRepository(db *gorm.DB) *ProjectRoleRepository {
	return &ProjectRoleRepository{db: db}
}

// Create creates a new project role and returns it with database-generated fields populated
func (r *ProjectRoleRepository) Create(ctx context.Context, projectRole *entities.ProjectRole) (*entities.ProjectRole, error) {
	err := r.db.WithContext(ctx).Create(projectRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "project_role", "method", "Create", "error", err)
			return nil, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "project_role", "method", "Create", "error", err)
			return nil, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			internal.Logger.Error("duplicated key", "repository", "project_role", "method", "Create", "error", err)
			return nil, entities.ErrDuplicatedKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "project_role", "method", "Create", "error", err)
			return nil, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "project_role", "method", "Create", "error", err)
			return nil, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to create project role", "repository", "project_role", "method", "Create", "error", err)
		return nil, err
	}
	return projectRole, nil
}

// GetOne gets a project role by ID
func (r *ProjectRoleRepository) GetOne(ctx context.Context, id uint) (*entities.ProjectRole, error) {
	var projectRole entities.ProjectRole
	err := r.db.WithContext(ctx).Model(&entities.ProjectRole{}).First(&projectRole, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "project_role", "method", "GetOne", "error", err)
			return nil, entities.ErrRecordNotFound
		}
		internal.Logger.Error("failed to get project role", "repository", "project_role", "method", "GetOne", "error", err)
		return nil, err
	}
	return &projectRole, err
}

// GetMany gets multiple project roles by query parameters
func (r *ProjectRoleRepository) GetMany(ctx context.Context, qParams *entities.ProjectRoleQueryParams) ([]*entities.ProjectRole, int64, error) {
	var (
		projectRoles []*entities.ProjectRole
		count        int64 = 0
	)
	q := r.db.WithContext(ctx).Model(&entities.ProjectRole{})

	if qParams == nil {
		qParams = &entities.ProjectRoleQueryParams{}
	}

	if len(qParams.ID_In) > 0 {
		q = q.Where("id IN @ID_In", sql.Named("ID_In", qParams.ID_In))
	}
	if qParams.ProjectID != 0 {
		q = q.Where("project_id = @ProjectID", sql.Named("ProjectID", qParams.ProjectID))
	}
	if len(qParams.ProjectID_In) > 0 {
		q = q.Where("project_id IN ?", qParams.ProjectID_In)
	}
	if qParams.Name != "" {
		q = q.Where("name = @Name", sql.Named("Name", qParams.Name))
	}
	if qParams.Name_Like != "" {
		q = q.Where("name LIKE ?", "%"+qParams.Name_Like+"%")
	}
	if qParams.Level != entities.RoleLevelUnknown {
		q = q.Where("level = @Level", sql.Named("Level", qParams.Level))
	}
	if len(qParams.Level_In) > 0 {
		q = q.Where("level IN ?", qParams.Level_In)
	}
	if qParams.Headcount_Gte != nil {
		q = q.Where("headcount >= @Headcount_Gte", sql.Named("Headcount_Gte", *qParams.Headcount_Gte))
	}
	if qParams.Headcount_Lte != nil {
		q = q.Where("headcount <= @Headcount_Lte", sql.Named("Headcount_Lte", *qParams.Headcount_Lte))
	}
	if qParams.CreatedAt_Gte != nil {
		q = q.Where("created_at >= @CreatedAt_Gte", sql.Named("CreatedAt_Gte", qParams.CreatedAt_Gte))
	}
	if qParams.CreatedAt_Lte != nil {
		q = q.Where("created_at <= @CreatedAt_Lte", sql.Named("CreatedAt_Lte", qParams.CreatedAt_Lte))
	}
	if qParams.UpdatedAt_Gte != nil {
		q = q.Where("updated_at >= @UpdatedAt_Gte", sql.Named("UpdatedAt_Gte", qParams.UpdatedAt_Gte))
	}
	if qParams.UpdatedAt_Lte != nil {
		q = q.Where("updated_at <= @UpdatedAt_Lte", sql.Named("UpdatedAt_Lte", qParams.UpdatedAt_Lte))
	}

	q = q.Session(&gorm.Session{})
	result := q.Count(&count)
	if result.Error != nil {
		internal.Logger.Error("failed to count project roles", "repository", "project_role", "method", "GetMany", "error", result.Error)
		return nil, 0, result.Error
	}

	// Apply sorting params
	if qParams.QueryParams != nil {
		if qParams.Sorts != nil {
			for _, sort := range qParams.Sorts {
				q = sort.Apply(q, entities.ProjectRoleAllowedSortField)
			}
		}
		if qParams.Pagination != nil {
			q = qParams.Pagination.Apply(q)
		}
	}

	// Execute query
	result = q.Find(&projectRoles)
	if result.Error != nil {
		internal.Logger.Error("failed to get project roles", "repository", "project_role", "method", "GetMany", "error", result.Error)
		return nil, count, result.Error
	}
	return projectRoles, count, nil
}

// Update updates a project role and returns it with updated database fields
func (r *ProjectRoleRepository) Update(ctx context.Context, projectRole *entities.ProjectRole) (int64, error) {
	result := r.db.WithContext(ctx).Model(projectRole).Clauses(clause.Returning{}).Where("id = ?", projectRole.ID).Select("*").Updates(&projectRole)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "project_role", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "project_role", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "project_role", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "project_role", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "project_role", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrCheckConstraintViolated
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			internal.Logger.Error("duplicated key", "repository", "project_role", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrDuplicatedKey
		}
		internal.Logger.Error("failed to update project role", "repository", "project_role", "method", "Update", "error", err)
		return result.RowsAffected, err
	}
	return result.RowsAffected, nil
}

// Delete deletes a project role by ID
func (r *ProjectRoleRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entities.ProjectRole{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "project_role", "method", "Delete", "error", err)
			return entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "project_role", "method", "Delete", "error", err)
			return entities.ErrForeignKeyViolated
		}
		internal.Logger.Error("failed to delete project role", "repository", "project_role", "method", "Delete", "error", err)
		return err
	}
	return nil
}

// GetByProjectNameAndLevel gets a project role by project ID, name, and level
func (r *ProjectRoleRepository) GetByProjectNameAndLevel(ctx context.Context, projectID uint, name string, level uint) (*entities.ProjectRole, error) {
	var projectRole entities.ProjectRole
	err := r.db.WithContext(ctx).Model(&entities.ProjectRole{}).
		Where("project_id = ? AND name = ? AND level = ?", projectID, name, level).
		First(&projectRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "project_role", "method", "GetByProjectNameAndLevel", "error", err)
			return nil, entities.ErrRecordNotFound
		}
		internal.Logger.Error("failed to get project role", "repository", "project_role", "method", "GetByProjectNameAndLevel", "error", err)
		return nil, err
	}
	return &projectRole, nil
}
