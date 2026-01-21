package handlers

// Handlers holds all handler dependencies
type Handlers struct {
	*ClientHandler
	*HumanResourceHandler
	*ProjectHandler
	*ProjectResourceHandler
	*ProjectRoleHandler
	*MilestoneHandler
	*TaskHandler
}

// NewHandlers creates a new Handlers instance with all handler dependencies
func NewHandlers(clientHandler *ClientHandler, hrHandler *HumanResourceHandler, projectHandler *ProjectHandler, projectResourceHandler *ProjectResourceHandler, projectRoleHandler *ProjectRoleHandler, milestoneHandler *MilestoneHandler, taskHandler *TaskHandler) *Handlers {
	return &Handlers{
		ClientHandler:          clientHandler,
		HumanResourceHandler:   hrHandler,
		ProjectHandler:         projectHandler,
		ProjectResourceHandler: projectResourceHandler,
		ProjectRoleHandler:     projectRoleHandler,
		MilestoneHandler:       milestoneHandler,
		TaskHandler:            taskHandler,
	}
}
