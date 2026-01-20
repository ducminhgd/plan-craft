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

// TaskRepository is the repository for task entities
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create creates a new task and returns it with database-generated fields populated
func (r *TaskRepository) Create(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	err := r.db.WithContext(ctx).Create(task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "task", "method", "Create", "error", err)
			return nil, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "task", "method", "Create", "error", err)
			return nil, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			internal.Logger.Error("duplicated key", "repository", "task", "method", "Create", "error", err)
			return nil, entities.ErrDuplicatedKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "task", "method", "Create", "error", err)
			return nil, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "task", "method", "Create", "error", err)
			return nil, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to create task", "repository", "task", "method", "Create", "error", err)
		return nil, err
	}
	return task, nil
}

// GetOne gets a task by ID
func (r *TaskRepository) GetOne(ctx context.Context, id uint) (*entities.Task, error) {
	var task entities.Task
	err := r.db.WithContext(ctx).Model(&entities.Task{}).First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "task", "method", "GetOne", "error", err)
			return nil, entities.ErrRecordNotFound
		}
		internal.Logger.Error("failed to get task", "repository", "task", "method", "GetOne", "error", err)
		return nil, err
	}
	return &task, err
}

// GetMany gets multiple tasks by query parameters
func (r *TaskRepository) GetMany(ctx context.Context, qParams *entities.TaskQueryParams) ([]*entities.Task, int64, error) {
	var (
		tasks []*entities.Task
		count int64 = 0
	)
	q := r.db.WithContext(ctx).Model(&entities.Task{})

	if qParams == nil {
		qParams = &entities.TaskQueryParams{}
	}

	if len(qParams.ID_In) > 0 {
		q = q.Where("id IN @ID_In", sql.Named("ID_In", qParams.ID_In))
	}
	if qParams.Name != "" {
		q = q.Where("name = @Name", sql.Named("Name", qParams.Name))
	}
	if qParams.Level != 0 {
		q = q.Where("level = @Level", sql.Named("Level", qParams.Level))
	}
	if qParams.Level_Gte != nil {
		q = q.Where("level >= @Level_Gte", sql.Named("Level_Gte", *qParams.Level_Gte))
	}
	if qParams.Level_Lte != nil {
		q = q.Where("level <= @Level_Lte", sql.Named("Level_Lte", *qParams.Level_Lte))
	}
	if qParams.ProjectID != 0 {
		q = q.Where("project_id = @ProjectID", sql.Named("ProjectID", qParams.ProjectID))
	}
	if len(qParams.ProjectID_In) > 0 {
		q = q.Where("project_id IN ?", qParams.ProjectID_In)
	}
	if qParams.MilestoneID != nil {
		q = q.Where("milestone_id = @MilestoneID", sql.Named("MilestoneID", *qParams.MilestoneID))
	}
	if len(qParams.MilestoneID_In) > 0 {
		q = q.Where("milestone_id IN ?", qParams.MilestoneID_In)
	}
	if qParams.MilestoneID_IsNull != nil {
		if *qParams.MilestoneID_IsNull {
			q = q.Where("milestone_id IS NULL")
		} else {
			q = q.Where("milestone_id IS NOT NULL")
		}
	}
	if qParams.ParentID != nil {
		q = q.Where("parent_id = @ParentID", sql.Named("ParentID", *qParams.ParentID))
	}
	if len(qParams.ParentID_In) > 0 {
		q = q.Where("parent_id IN ?", qParams.ParentID_In)
	}
	if qParams.ParentID_IsNull != nil {
		if *qParams.ParentID_IsNull {
			q = q.Where("parent_id IS NULL")
		} else {
			q = q.Where("parent_id IS NOT NULL")
		}
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

	if qParams.Priority != entities.TaskPriorityUnknown {
		q = q.Where("priority = @Priority", sql.Named("Priority", qParams.Priority))
	}
	if len(qParams.Priority_In) > 0 {
		q = q.Where("priority IN ?", qParams.Priority_In)
	}
	if qParams.Status != entities.TaskWorkStatusUnknown {
		q = q.Where("status = @Status", sql.Named("Status", qParams.Status))
	}
	if len(qParams.Status_In) > 0 {
		q = q.Where("status IN ?", qParams.Status_In)
	}
	if qParams.EstimatedEffort_Gte != nil {
		q = q.Where("estimated_effort >= @EstimatedEffort_Gte", sql.Named("EstimatedEffort_Gte", *qParams.EstimatedEffort_Gte))
	}
	if qParams.EstimatedEffort_Lte != nil {
		q = q.Where("estimated_effort <= @EstimatedEffort_Lte", sql.Named("EstimatedEffort_Lte", *qParams.EstimatedEffort_Lte))
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
		internal.Logger.Error("failed to count tasks", "repository", "task", "method", "GetMany", "error", result.Error)
		return nil, 0, result.Error
	}

	// Apply sorting params
	if qParams.QueryParams != nil {
		if qParams.Sorts != nil {
			for _, sort := range qParams.Sorts {
				q = sort.Apply(q, entities.TaskAllowedSortField)
			}
		}
		if qParams.Pagination != nil {
			q = qParams.Pagination.Apply(q)
		}
	}

	// Execute query
	result = q.Find(&tasks)
	if result.Error != nil {
		internal.Logger.Error("failed to get tasks", "repository", "task", "method", "GetMany", "error", result.Error)
		return nil, count, result.Error
	}
	return tasks, count, nil
}

// Update updates a task and returns it with updated database fields
func (r *TaskRepository) Update(ctx context.Context, task *entities.Task) (int64, error) {
	result := r.db.WithContext(ctx).Model(task).Clauses(clause.Returning{}).Where("id = ?", task.ID).Select("*").Updates(&task)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "task", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "task", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "task", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "task", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to update task", "repository", "task", "method", "Update", "error", err)
		return result.RowsAffected, err
	}
	// Check if no rows were affected (record not found)
	if result.RowsAffected == 0 {
		return 0, entities.ErrRecordNotFound
	}
	return result.RowsAffected, nil
}

// Delete deletes a task by ID
func (r *TaskRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entities.Task{}, id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "task", "method", "Delete", "error", err)
			return entities.ErrForeignKeyViolated
		}
		internal.Logger.Error("failed to delete task", "repository", "task", "method", "Delete", "error", err)
		return err
	}
	// Check if no rows were affected (record not found)
	if result.RowsAffected == 0 {
		return entities.ErrRecordNotFound
	}
	return nil
}
