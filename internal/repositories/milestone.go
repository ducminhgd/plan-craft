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

// MilestoneRepository is the repository for milestone entities
type MilestoneRepository struct {
	db *gorm.DB
}

// NewMilestoneRepository creates a new milestone repository
func NewMilestoneRepository(db *gorm.DB) *MilestoneRepository {
	return &MilestoneRepository{db: db}
}

// Create creates a new milestone and returns it with database-generated fields populated
func (r *MilestoneRepository) Create(ctx context.Context, milestone *entities.Milestone) (*entities.Milestone, error) {
	err := r.db.WithContext(ctx).Create(milestone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "milestone", "method", "Create", "error", err)
			return nil, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "milestone", "method", "Create", "error", err)
			return nil, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			internal.Logger.Error("duplicated key", "repository", "milestone", "method", "Create", "error", err)
			return nil, entities.ErrDuplicatedKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "milestone", "method", "Create", "error", err)
			return nil, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "milestone", "method", "Create", "error", err)
			return nil, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to create milestone", "repository", "milestone", "method", "Create", "error", err)
		return nil, err
	}
	return milestone, nil
}

// GetOne gets a milestone by ID
func (r *MilestoneRepository) GetOne(ctx context.Context, id uint) (*entities.Milestone, error) {
	var milestone entities.Milestone
	err := r.db.WithContext(ctx).Model(&entities.Milestone{}).First(&milestone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "milestone", "method", "GetOne", "error", err)
			return nil, entities.ErrRecordNotFound
		}
		internal.Logger.Error("failed to get milestone", "repository", "milestone", "method", "GetOne", "error", err)
		return nil, err
	}
	return &milestone, err
}

// GetMany gets multiple milestones by query parameters
func (r *MilestoneRepository) GetMany(ctx context.Context, qParams *entities.MilestoneQueryParams) ([]*entities.Milestone, int64, error) {
	var (
		milestones []*entities.Milestone
		count      int64 = 0
	)
	q := r.db.WithContext(ctx).Model(&entities.Milestone{})

	if qParams == nil {
		qParams = &entities.MilestoneQueryParams{}
	}

	if len(qParams.ID_In) > 0 {
		q = q.Where("id IN @ID_In", sql.Named("ID_In", qParams.ID_In))
	}
	if qParams.Name != "" {
		q = q.Where("name = @Name", sql.Named("Name", qParams.Name))
	}
	if qParams.ProjectID != 0 {
		q = q.Where("project_id = @ProjectID", sql.Named("ProjectID", qParams.ProjectID))
	}
	if len(qParams.ProjectID_In) > 0 {
		q = q.Where("project_id IN ?", qParams.ProjectID_In)
	}

	// Group LIKE conditions with OR for search functionality
	if qParams.Name_Like != "" || qParams.Description_Like != "" {
		orConditions := r.db.Where("1 = 0") // Start with false condition

		if qParams.Name_Like != "" {
			orConditions = orConditions.Or("name LIKE ?", "%"+qParams.Name_Like+"%")
		}
		if qParams.Description_Like != "" {
			orConditions = orConditions.Or("description LIKE ?", "%"+qParams.Description_Like+"%")
		}

		q = q.Where(orConditions)
	}

	if qParams.Status != entities.MilestoneStatusUnknown {
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
		internal.Logger.Error("failed to count milestones", "repository", "milestone", "method", "GetMany", "error", result.Error)
		return nil, 0, result.Error
	}

	// Apply sorting params
	if qParams.QueryParams != nil {
		if qParams.Sorts != nil {
			for _, sort := range qParams.Sorts {
				q = sort.Apply(q, entities.MilestoneAllowedSortField)
			}
		}
		if qParams.Pagination != nil {
			q = qParams.Pagination.Apply(q)
		}
	}

	// Execute query
	result = q.Find(&milestones)
	if result.Error != nil {
		internal.Logger.Error("failed to get milestones", "repository", "milestone", "method", "GetMany", "error", result.Error)
		return nil, count, result.Error
	}
	return milestones, count, nil
}

// Update updates a milestone and returns it with updated database fields
func (r *MilestoneRepository) Update(ctx context.Context, milestone *entities.Milestone) (int64, error) {
	result := r.db.WithContext(ctx).Model(milestone).Clauses(clause.Returning{}).Where("id = ?", milestone.ID).Select("*").Updates(&milestone)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "milestone", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "milestone", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "milestone", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "milestone", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "milestone", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to update milestone", "repository", "milestone", "method", "Update", "error", err)
		return result.RowsAffected, err
	}
	return result.RowsAffected, nil
}

// Delete deletes a milestone by ID
func (r *MilestoneRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entities.Milestone{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "milestone", "method", "Delete", "error", err)
			return entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "milestone", "method", "Delete", "error", err)
			return entities.ErrForeignKeyViolated
		}
		internal.Logger.Error("failed to delete milestone", "repository", "milestone", "method", "Delete", "error", err)
		return err
	}
	return nil
}
