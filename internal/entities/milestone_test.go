package entities

import (
	"testing"
	"time"
)

func TestMilestone_TableName(t *testing.T) {
	milestone := Milestone{}
	if milestone.TableName() != "milestones" {
		t.Errorf("Expected table name 'milestones', got '%s'", milestone.TableName())
	}
}

func TestMilestone_BeforeSave_ValidMilestone(t *testing.T) {
	milestone := Milestone{
		ProjectID: 1,
		Name:      "Phase 1",
		Status:    TaskStatusNotStarted,
		Progress:  0,
	}

	err := milestone.BeforeSave(nil)
	if err != nil {
		t.Errorf("Expected no error for valid milestone, got: %v", err)
	}
}

func TestMilestone_BeforeSave_InvalidStatus(t *testing.T) {
	milestone := Milestone{
		ProjectID: 1,
		Name:      "Phase 1",
		Status:    TaskStatus("invalid"),
	}

	err := milestone.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for invalid status, got nil")
	}
}

func TestMilestone_BeforeSave_InvalidProgress(t *testing.T) {
	tests := []struct {
		name     string
		progress float64
		wantErr  bool
	}{
		{"Valid progress 0", 0, false},
		{"Valid progress 50", 50, false},
		{"Valid progress 100", 100, false},
		{"Invalid progress -1", -1, true},
		{"Invalid progress 101", 101, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			milestone := Milestone{
				ProjectID: 1,
				Name:      "Test Milestone",
				Status:    TaskStatusNotStarted,
				Progress:  tt.progress,
			}

			err := milestone.BeforeSave(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeforeSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMilestone_BeforeSave_InvalidDates(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.Add(-24 * time.Hour) // End date before start date

	milestone := Milestone{
		ProjectID:        1,
		Name:             "Test Milestone",
		Status:           TaskStatusNotStarted,
		PlannedStartDate: &startDate,
		PlannedEndDate:   &endDate,
	}

	err := milestone.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for end date before start date, got nil")
	}
}

func TestMilestone_IsCompleted(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{"Completed", TaskStatusCompleted, true},
		{"In progress", TaskStatusInProgress, false},
		{"Not started", TaskStatusNotStarted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			milestone := Milestone{Status: tt.status}
			if got := milestone.IsCompleted(); got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMilestone_PlannedDuration(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		startDate *time.Time
		endDate   *time.Time
		wantNil   bool
	}{
		{"Both dates set", &startDate, &endDate, false},
		{"No start date", nil, &endDate, true},
		{"No end date", &startDate, nil, true},
		{"No dates", nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			milestone := Milestone{
				PlannedStartDate: tt.startDate,
				PlannedEndDate:   tt.endDate,
			}
			duration := milestone.PlannedDuration()
			if (duration == nil) != tt.wantNil {
				t.Errorf("PlannedDuration() nil = %v, want nil %v", duration == nil, tt.wantNil)
			}
			if !tt.wantNil && duration != nil {
				expected := endDate.Sub(startDate)
				if *duration != expected {
					t.Errorf("PlannedDuration() = %v, want %v", *duration, expected)
				}
			}
		})
	}
}

func TestMilestone_EstimatedConversions(t *testing.T) {
	milestone := Milestone{
		EstimatedEffort: 160, // 160 hours
	}

	t.Run("EstimatedDays", func(t *testing.T) {
		expected := 20.0 // 160 hours / 8 hours per day
		if got := milestone.EstimatedDays(); got != expected {
			t.Errorf("EstimatedDays() = %v, want %v", got, expected)
		}
	})

	t.Run("EstimatedMonths", func(t *testing.T) {
		expected := 1.0 // 160 hours / 160 hours per month
		if got := milestone.EstimatedMonths(); got != expected {
			t.Errorf("EstimatedMonths() = %v, want %v", got, expected)
		}
	})
}

func TestMilestone_ActualConversions(t *testing.T) {
	milestone := Milestone{
		ActualEffort: 80, // 80 hours
	}

	t.Run("ActualDays", func(t *testing.T) {
		expected := 10.0 // 80 hours / 8 hours per day
		if got := milestone.ActualDays(); got != expected {
			t.Errorf("ActualDays() = %v, want %v", got, expected)
		}
	})

	t.Run("ActualMonths", func(t *testing.T) {
		expected := 0.5 // 80 hours / 160 hours per month
		if got := milestone.ActualMonths(); got != expected {
			t.Errorf("ActualMonths() = %v, want %v", got, expected)
		}
	})
}
