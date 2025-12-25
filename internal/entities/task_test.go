package entities

import "testing"

func TestTask_TableName(t *testing.T) {
	task := Task{}
	if task.TableName() != "tasks" {
		t.Errorf("Expected table name 'tasks', got '%s'", task.TableName())
	}
}

func TestTask_BeforeSave_ValidTask(t *testing.T) {
	task := Task{
		ProjectID:      1,
		Name:           "Implement feature",
		Status:         TaskStatusNotStarted,
		Priority:       PriorityMedium,
		Level:          1,
		Progress:       0,
		EstimatedHours: 40,
	}

	err := task.BeforeSave(nil)
	if err != nil {
		t.Errorf("Expected no error for valid task, got: %v", err)
	}
}

func TestTask_IsEpic(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		parentID *uint
		want     bool
	}{
		{"Epic (level 1, no parent)", 1, nil, true},
		{"Not epic (level 2)", 2, nil, false},
		{"Not epic (has parent)", 1, uintPtr(1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				Level:    tt.level,
				ParentID: tt.parentID,
			}
			if got := task.IsEpic(); got != tt.want {
				t.Errorf("IsEpic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_IsSubtask(t *testing.T) {
	tests := []struct {
		name     string
		parentID *uint
		want     bool
	}{
		{"Has parent", uintPtr(1), true},
		{"No parent", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{ParentID: tt.parentID}
			if got := task.IsSubtask(); got != tt.want {
				t.Errorf("IsSubtask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_EstimatedConversions(t *testing.T) {
	task := Task{
		EstimatedHours: 80,
	}

	t.Run("EstimatedDays", func(t *testing.T) {
		expected := 10.0 // 80 hours / 8 hours per day
		if got := task.EstimatedDays(); got != expected {
			t.Errorf("EstimatedDays() = %v, want %v", got, expected)
		}
	})

	t.Run("EstimatedMonths", func(t *testing.T) {
		expected := 0.5 // 80 hours / 160 hours per month
		if got := task.EstimatedMonths(); got != expected {
			t.Errorf("EstimatedMonths() = %v, want %v", got, expected)
		}
	})
}

func TestTask_ActualConversions(t *testing.T) {
	task := Task{
		ActualHours: 40,
	}

	t.Run("ActualDays", func(t *testing.T) {
		expected := 5.0 // 40 hours / 8 hours per day
		if got := task.ActualDays(); got != expected {
			t.Errorf("ActualDays() = %v, want %v", got, expected)
		}
	})

	t.Run("ActualMonths", func(t *testing.T) {
		expected := 0.25 // 40 hours / 160 hours per month
		if got := task.ActualMonths(); got != expected {
			t.Errorf("ActualMonths() = %v, want %v", got, expected)
		}
	})
}

func TestTaskDependency_TableName(t *testing.T) {
	td := TaskDependency{}
	if td.TableName() != "task_dependencies" {
		t.Errorf("Expected table name 'task_dependencies', got '%s'", td.TableName())
	}
}

func TestTaskDependency_BeforeSave_ValidDependency(t *testing.T) {
	td := TaskDependency{
		PredecessorTaskID: 1,
		DependentTaskID:   2,
		DependencyType:    DependencyFinishToStart,
		LagDays:           0,
		LeadDays:          0,
	}

	err := td.BeforeSave(nil)
	if err != nil {
		t.Errorf("Expected no error for valid dependency, got: %v", err)
	}
}

func TestTaskDependency_BeforeSave_SelfDependency(t *testing.T) {
	td := TaskDependency{
		PredecessorTaskID: 1,
		DependentTaskID:   1,
		DependencyType:    DependencyFinishToStart,
	}

	err := td.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for self-dependency, got nil")
	}
}

func TestTaskDependency_BeforeSave_BothLagAndLead(t *testing.T) {
	td := TaskDependency{
		PredecessorTaskID: 1,
		DependentTaskID:   2,
		DependencyType:    DependencyFinishToStart,
		LagDays:           5,
		LeadDays:          3,
	}

	err := td.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for both lag and lead time set, got nil")
	}
}

func TestTaskDependency_BeforeSave_InvalidDependencyType(t *testing.T) {
	td := TaskDependency{
		PredecessorTaskID: 1,
		DependentTaskID:   2,
		DependencyType:    DependencyType("invalid"),
	}

	err := td.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for invalid dependency type, got nil")
	}
}

// Helper function
func uintPtr(u uint) *uint {
	return &u
}
