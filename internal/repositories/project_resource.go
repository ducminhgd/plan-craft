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

// ProjectResourceRepository is the repository for project resource entities
type ProjectResourceRepository struct {
	db *gorm.DB
}

// NewProjectResourceRepository creates a new project resource repository
func NewProjectResourceRepository(db *gorm.DB) *ProjectResourceRepository {
	return &ProjectResourceRepository{db: db}
}

// Create creates a new project resource and returns it with database-generated fields populated
func (r *ProjectResourceRepository) Create(ctx context.Context, projectResource *entities.ProjectResource) (*entities.ProjectResource, error) {
	err := r.db.WithContext(ctx).Create(projectResource).Error
	if err != nil {
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "project_resource", "method", "Create", "error", err)
			return nil, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "project_resource", "method", "Create", "error", err)
			return nil, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			internal.Logger.Error("duplicated key", "repository", "project_resource", "method", "Create", "error", err)
			return nil, entities.ErrDuplicatedKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "project_resource", "method", "Create", "error", err)
			return nil, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "project_resource", "method", "Create", "error", err)
			return nil, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to create project resource", "repository", "project_resource", "method", "Create", "error", err)
		return nil, err
	}
	return projectResource, nil
}

// GetOne gets a project resource by ID
func (r *ProjectResourceRepository) GetOne(ctx context.Context, id uint) (*entities.ProjectResource, error) {
	var projectResource entities.ProjectResource
	err := r.db.WithContext(ctx).Model(&entities.ProjectResource{}).First(&projectResource, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "project_resource", "method", "GetOne", "error", err)
			return nil, entities.ErrRecordNotFound
		}
		internal.Logger.Error("failed to get project resource", "repository", "project_resource", "method", "GetOne", "error", err)
		return nil, err
	}
	return &projectResource, err
}

// GetMany gets multiple project resources by query parameters
func (r *ProjectResourceRepository) GetMany(ctx context.Context, qParams *entities.ProjectResourceQueryParams) ([]*entities.ProjectResource, int64, error) {
	var (
		projectResources []*entities.ProjectResource
		count            int64 = 0
	)
	q := r.db.WithContext(ctx).Model(&entities.ProjectResource{})

	if qParams == nil {
		qParams = &entities.ProjectResourceQueryParams{}
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
	if qParams.HumanResourceID != 0 {
		q = q.Where("human_resource_id = @HumanResourceID", sql.Named("HumanResourceID", qParams.HumanResourceID))
	}
	if len(qParams.HumanResourceID_In) > 0 {
		q = q.Where("human_resource_id IN ?", qParams.HumanResourceID_In)
	}
	if qParams.Role != "" {
		q = q.Where("role = @Role", sql.Named("Role", qParams.Role))
	}
	if qParams.Role_Like != "" {
		q = q.Where("role LIKE ?", "%"+qParams.Role_Like+"%")
	}
	if qParams.Allocation_Gte != nil {
		q = q.Where("allocation >= @Allocation_Gte", sql.Named("Allocation_Gte", *qParams.Allocation_Gte))
	}
	if qParams.Allocation_Lte != nil {
		q = q.Where("allocation <= @Allocation_Lte", sql.Named("Allocation_Lte", *qParams.Allocation_Lte))
	}
	if qParams.Status != entities.ProjectResourceStatusUnknown {
		q = q.Where("status = @Status", sql.Named("Status", qParams.Status))
	}
	if len(qParams.Status_In) > 0 {
		q = q.Where("status IN ?", qParams.Status_In)
	}
	if qParams.StartDate_Gte != nil {
		q = q.Where("start_date >= @StartDate_Gte", sql.Named("StartDate_Gte", qParams.StartDate_Gte))
	}
	if qParams.StartDate_Lte != nil {
		q = q.Where("start_date <= @StartDate_Lte", sql.Named("StartDate_Lte", qParams.StartDate_Lte))
	}
	if qParams.EndDate_Gte != nil {
		q = q.Where("end_date >= @EndDate_Gte", sql.Named("EndDate_Gte", qParams.EndDate_Gte))
	}
	if qParams.EndDate_Lte != nil {
		q = q.Where("end_date <= @EndDate_Lte", sql.Named("EndDate_Lte", qParams.EndDate_Lte))
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
		internal.Logger.Error("failed to count project resources", "repository", "project_resource", "method", "GetMany", "error", result.Error)
		return nil, 0, result.Error
	}

	// Apply sorting params
	if qParams.QueryParams != nil {
		if qParams.Sorts != nil {
			for _, sort := range qParams.Sorts {
				q = sort.Apply(q, entities.ProjectResourceAllowedSortField)
			}
		}
		if qParams.Pagination != nil {
			q = qParams.Pagination.Apply(q)
		}
	}

	// Execute query
	result = q.Find(&projectResources)
	if result.Error != nil {
		internal.Logger.Error("failed to get project resources", "repository", "project_resource", "method", "GetMany", "error", result.Error)
		return nil, count, result.Error
	}
	return projectResources, count, nil
}

// Update updates a project resource and returns it with updated database fields
func (r *ProjectResourceRepository) Update(ctx context.Context, projectResource *entities.ProjectResource) (int64, error) {
	result := r.db.WithContext(ctx).Model(projectResource).Clauses(clause.Returning{}).Where("id = ?", projectResource.ID).Select("*").Updates(&projectResource)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "project_resource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "project_resource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "project_resource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "project_resource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "project_resource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to update project resource", "repository", "project_resource", "method", "Update", "error", err)
		return result.RowsAffected, err
	}
	return result.RowsAffected, nil
}

// Delete deletes a project resource by ID
func (r *ProjectResourceRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entities.ProjectResource{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "project_resource", "method", "Delete", "error", err)
			return entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "project_resource", "method", "Delete", "error", err)
			return entities.ErrForeignKeyViolated
		}
		internal.Logger.Error("failed to delete project resource", "repository", "project_resource", "method", "Delete", "error", err)
		return err
	}
	return nil
}

// GetByProjectAndResource gets a project resource by project ID and human resource ID
func (r *ProjectResourceRepository) GetByProjectAndResource(ctx context.Context, projectID, humanResourceID uint) (*entities.ProjectResource, error) {
	var projectResource entities.ProjectResource
	err := r.db.WithContext(ctx).Model(&entities.ProjectResource{}).
		Where("project_id = ? AND human_resource_id = ?", projectID, humanResourceID).
		First(&projectResource).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "project_resource", "method", "GetByProjectAndResource", "error", err)
			return nil, entities.ErrRecordNotFound
		}
		internal.Logger.Error("failed to get project resource", "repository", "project_resource", "method", "GetByProjectAndResource", "error", err)
		return nil, err
	}
	return &projectResource, nil
}
