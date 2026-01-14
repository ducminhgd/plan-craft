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

// HRRepository is the repository for human resource entities
type HRRepository struct {
	db *gorm.DB
}

// NewHRRepository creates a new human resource repository
func NewHRRepository(db *gorm.DB) *HRRepository {
	return &HRRepository{db: db}
}

// Create creates a new human resource and returns it with database-generated fields populated
func (r *HRRepository) Create(ctx context.Context, humanResource *entities.HumanResource) (*entities.HumanResource, error) {
	err := r.db.WithContext(ctx).Create(humanResource).Error
	if err != nil {
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "humanResource", "method", "Create", "error", err)
			return nil, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "humanResource", "method", "Create", "error", err)
			return nil, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			internal.Logger.Error("duplicated key", "repository", "humanResource", "method", "Create", "error", err)
			return nil, entities.ErrDuplicatedKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "humanResource", "method", "Create", "error", err)
			return nil, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "humanResource", "method", "Create", "error", err)
			return nil, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to create human resource", "repository", "humanResource", "method", "Create", "error", err)
		return nil, err
	}
	return humanResource, nil
}

// GetOne gets a human resource by ID
func (r *HRRepository) GetOne(ctx context.Context, id uint) (*entities.HumanResource, error) {
	var humanResource entities.HumanResource
	err := r.db.WithContext(ctx).Model(&entities.HumanResource{}).First(&humanResource, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "humanResource", "method", "GetOne", "error", err)
			return nil, entities.ErrRecordNotFound
		}
		internal.Logger.Error("failed to get human resource", "repository", "humanResource", "method", "GetOne", "error", err)
		return nil, err
	}
	return &humanResource, err
}

// GetMany gets multiple human resources by query parameters
func (r *HRRepository) GetMany(ctx context.Context, qParams *entities.HumanResourceQueryParams) ([]*entities.HumanResource, int64, error) {
	var (
		humanResources   []*entities.HumanResource
		count int64 = 0
	)
	q := r.db.WithContext(ctx).Model(&entities.HumanResource{})

	if qParams == nil {
		qParams = &entities.HumanResourceQueryParams{}
	}

	if len(qParams.ID_In) > 0 {
		q = q.Where("id IN @ID_In", sql.Named("ID_In", qParams.ID_In))
	}
	if qParams.Name != "" {
		q = q.Where("name = @Name", sql.Named("Name", qParams.Name))
	}
	if qParams.Title != "" {
		q = q.Where("title = @Title", sql.Named("Title", qParams.Title))
	}
	if qParams.Level != "" {
		q = q.Where("level = @Level", sql.Named("Level", qParams.Level))
	}

	// Group LIKE conditions with OR for search functionality
	// These fields are used by the frontend search box and should match ANY of them
	if qParams.Name_Like != "" || qParams.Title_Like != "" || qParams.Level_Like != "" {

		orConditions := r.db.Where("1 = 0") // Start with false condition

		if qParams.Name_Like != "" {
			orConditions = orConditions.Or("name LIKE ?", "%"+qParams.Name_Like+"%")
		}
		if qParams.Title_Like != "" {
			orConditions = orConditions.Or("title LIKE ?", "%"+qParams.Title_Like+"%")
		}
		if qParams.Level_Like != "" {
			orConditions = orConditions.Or("level LIKE ?", "%"+qParams.Level_Like+"%")
		}

		q = q.Where(orConditions)
	}
	if qParams.Status != entities.HumanResourceStatusUnknown {
		q = q.Where("status = @Status", sql.Named("Status", qParams.Status))
	}
	if len(qParams.Status_In) > 0 {
		q = q.Where("status IN ?", qParams.Status_In)
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
		internal.Logger.Error("failed to count human resources", "repository", "humanResource", "method", "GetMany", "error", result.Error)
		return nil, 0, result.Error
	}

	// Apply sorting params
	if qParams.QueryParams != nil {
		if qParams.Sorts != nil {
			for _, sort := range qParams.Sorts {
				q = sort.Apply(q, entities.HumanResourceAllowedSortField)
			}
		}
		if qParams.Pagination != nil {
			q = qParams.Pagination.Apply(q)
		}
	}

	// Execute query
	result = q.Find(&humanResources)
	if result.Error != nil {
		internal.Logger.Error("failed to get human resources", "repository", "humanResource", "method", "GetMany", "error", result.Error)
		return nil, count, result.Error
	}
	return humanResources, count, nil
}

// Update updates a human resource and returns it with updated database fields
func (r *HRRepository) Update(ctx context.Context, humanResource *entities.HumanResource) (int64, error) {
	result := r.db.WithContext(ctx).Model(humanResource).Clauses(clause.Returning{}).Where("id = ?", humanResource.ID).Select("*").Updates(&humanResource)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "humanResource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "humanResource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "humanResource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "humanResource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "humanResource", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to update human resource", "repository", "humanResource", "method", "Update", "error", err)
		return result.RowsAffected, err
	}
	return result.RowsAffected, nil
}

// Delete deletes a human resource by ID
func (r *HRRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entities.HumanResource{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "humanResource", "method", "Delete", "error", err)
			return entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "humanResource", "method", "Delete", "error", err)
			return entities.ErrForeignKeyViolated
		}
		internal.Logger.Error("failed to delete human resource", "repository", "humanResource", "method", "Delete", "error", err)
		return err
	}
	return nil
}
