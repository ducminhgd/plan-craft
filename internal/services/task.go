package services

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
)

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	Create(ctx context.Context, task *entities.Task) (*entities.Task, error)
	GetOne(ctx context.Context, id uint) (*entities.Task, error)
	GetMany(ctx context.Context, qParams *entities.TaskQueryParams) ([]*entities.Task, int64, error)
	Update(ctx context.Context, task *entities.Task) (int64, error)
	Delete(ctx context.Context, id uint) error
}

// TaskService handles task business logic
type TaskService struct {
	repo TaskRepository
}

// NewTaskService creates a new task service
func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	return s.repo.Create(ctx, task)
}

// GetTask retrieves a single task by ID
func (s *TaskService) GetTask(ctx context.Context, id uint) (*entities.Task, error) {
	return s.repo.GetOne(ctx, id)
}

// GetTasks retrieves multiple tasks with optional query parameters
func (s *TaskService) GetTasks(ctx context.Context, params *entities.TaskQueryParams) (*entities.TaskListResponse, error) {
	data, total, err := s.repo.GetMany(ctx, params)
	if err != nil {
		return nil, err
	}
	return &entities.TaskListResponse{
		Data:  data,
		Total: total,
	}, nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(ctx context.Context, task *entities.Task) (int64, error) {
	return s.repo.Update(ctx, task)
}

// DeleteTask deletes a task by ID
func (s *TaskService) DeleteTask(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
