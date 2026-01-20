package handlers

import (
	"context"
	"fmt"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// TaskHandler handles task-related operations for Wails bindings
type TaskHandler struct {
	ctx     context.Context
	service *services.TaskService
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(ctx context.Context, service *services.TaskService) *TaskHandler {
	return &TaskHandler{
		ctx:     ctx,
		service: service,
	}
}

// GetTasks retrieves multiple tasks with optional query parameters
func (h *TaskHandler) GetTasks(params *entities.TaskQueryParams) (*entities.TaskListResponse, error) {
	if h.service == nil {
		return nil, fmt.Errorf("task service not initialized")
	}
	return h.service.GetTasks(h.ctx, params)
}

// GetTask retrieves a single task by ID
func (h *TaskHandler) GetTask(id uint) (*entities.Task, error) {
	if h.service == nil {
		return nil, fmt.Errorf("task service not initialized")
	}
	return h.service.GetTask(h.ctx, id)
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(task *entities.Task) (*entities.Task, error) {
	if h.service == nil {
		return nil, fmt.Errorf("task service not initialized")
	}
	return h.service.CreateTask(h.ctx, task)
}

// UpdateTask updates an existing task
func (h *TaskHandler) UpdateTask(task *entities.Task) (int64, error) {
	if h.service == nil {
		return 0, fmt.Errorf("task service not initialized")
	}
	return h.service.UpdateTask(h.ctx, task)
}

// DeleteTask deletes a task by ID
func (h *TaskHandler) DeleteTask(id uint) error {
	if h.service == nil {
		return fmt.Errorf("task service not initialized")
	}
	return h.service.DeleteTask(h.ctx, id)
}
