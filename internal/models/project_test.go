package models

import (
	"testing"
	"time"
)

func TestProject_TableName(t *testing.T) {
	project := Project{}
	if project.TableName() != "projects" {
		t.Errorf("Expected table name 'projects', got '%s'", project.TableName())
	}
}

func TestProject_BeforeSave_ValidProject(t *testing.T) {
	project := Project{
		Name:     "Test Project",
		Type:     ProjectTypeProduct,
		Status:   TaskStatusNotStarted,
		Progress: 0,
	}

	err := project.BeforeSave(nil)
	if err != nil {
		t.Errorf("Expected no error for valid project, got: %v", err)
	}
}

func TestProject_BeforeSave_InvalidType(t *testing.T) {
	project := Project{
		Name:   "Test Project",
		Type:   ProjectType("invalid"),
		Status: TaskStatusNotStarted,
	}

	err := project.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for invalid project type, got nil")
	}
}

func TestProject_BeforeSave_InvalidStatus(t *testing.T) {
	project := Project{
		Name:   "Test Project",
		Type:   ProjectTypeProduct,
		Status: TaskStatus("invalid"),
	}

	err := project.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for invalid status, got nil")
	}
}

func TestProject_BeforeSave_InvalidProgress(t *testing.T) {
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
			project := Project{
				Name:     "Test Project",
				Type:     ProjectTypeProduct,
				Status:   TaskStatusNotStarted,
				Progress: tt.progress,
			}

			err := project.BeforeSave(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeforeSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProject_BeforeSave_InvalidDates(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.Add(-24 * time.Hour) // End date before start date

	project := Project{
		Name:          "Test Project",
		Type:          ProjectTypeProduct,
		Status:        TaskStatusNotStarted,
		StartDate:     &startDate,
		TargetEndDate: &endDate,
	}

	err := project.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for end date before start date, got nil")
	}
}

func TestProject_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{"In progress", TaskStatusInProgress, true},
		{"Not started", TaskStatusNotStarted, false},
		{"Completed", TaskStatusCompleted, false},
		{"On hold", TaskStatusOnHold, false},
		{"Cancelled", TaskStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := Project{Status: tt.status}
			if got := project.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_IsCompleted(t *testing.T) {
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
			project := Project{Status: tt.status}
			if got := project.IsCompleted(); got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_Duration(t *testing.T) {
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
			project := Project{
				StartDate:     tt.startDate,
				TargetEndDate: tt.endDate,
			}
			duration := project.Duration()
			if (duration == nil) != tt.wantNil {
				t.Errorf("Duration() nil = %v, want nil %v", duration == nil, tt.wantNil)
			}
			if !tt.wantNil && duration != nil {
				expected := endDate.Sub(startDate)
				if *duration != expected {
					t.Errorf("Duration() = %v, want %v", *duration, expected)
				}
			}
		})
	}
}

func TestProject_ActualDuration(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	actualEndDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		startDate     *time.Time
		actualEndDate *time.Time
		wantNil       bool
	}{
		{"Both dates set", &startDate, &actualEndDate, false},
		{"No start date", nil, &actualEndDate, true},
		{"No actual end date", &startDate, nil, true},
		{"No dates", nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := Project{
				StartDate:     tt.startDate,
				ActualEndDate: tt.actualEndDate,
			}
			duration := project.ActualDuration()
			if (duration == nil) != tt.wantNil {
				t.Errorf("ActualDuration() nil = %v, want nil %v", duration == nil, tt.wantNil)
			}
			if !tt.wantNil && duration != nil {
				expected := actualEndDate.Sub(startDate)
				if *duration != expected {
					t.Errorf("ActualDuration() = %v, want %v", *duration, expected)
				}
			}
		})
	}
}

func TestStringArray_ScanAndValue(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		wantErr bool
	}{
		{"Empty array", []string{}, false},
		{"Single item", []string{"item1"}, false},
		{"Multiple items", []string{"item1", "item2", "item3"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sa StringArray = tt.input

			// Test Value
			value, err := sa.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Test Scan
			var scanned StringArray
			err = scanned.Scan(value)
			if err != nil {
				t.Errorf("Scan() error = %v", err)
				return
			}

			// Compare
			if len(tt.input) > 0 {
				if len(scanned) != len(tt.input) {
					t.Errorf("Scanned length = %d, want %d", len(scanned), len(tt.input))
				}
				for i, v := range tt.input {
					if scanned[i] != v {
						t.Errorf("Scanned[%d] = %s, want %s", i, scanned[i], v)
					}
				}
			}
		})
	}
}

func TestStringArray_ScanNil(t *testing.T) {
	var sa StringArray
	err := sa.Scan(nil)
	if err != nil {
		t.Errorf("Scan(nil) error = %v, want nil", err)
	}
	if len(sa) != 0 {
		t.Errorf("Scan(nil) length = %d, want 0", len(sa))
	}
}

func TestJSONB_ScanAndValue(t *testing.T) {
	tests := []struct {
		name    string
		input   JSONB
		wantErr bool
	}{
		{"Empty map", JSONB{}, false},
		{"Single field", JSONB{"key": "value"}, false},
		{"Multiple fields", JSONB{"key1": "value1", "key2": 123, "key3": true}, false},
		{"Nested", JSONB{"outer": map[string]interface{}{"inner": "value"}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Value
			value, err := tt.input.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Test Scan
			var scanned JSONB
			err = scanned.Scan(value)
			if err != nil {
				t.Errorf("Scan() error = %v", err)
				return
			}

			// Compare length
			if len(tt.input) > 0 && len(scanned) != len(tt.input) {
				t.Errorf("Scanned length = %d, want %d", len(scanned), len(tt.input))
			}
		})
	}
}

func TestJSONB_ScanNil(t *testing.T) {
	var j JSONB
	err := j.Scan(nil)
	if err != nil {
		t.Errorf("Scan(nil) error = %v, want nil", err)
	}
	if len(j) != 0 {
		t.Errorf("Scan(nil) length = %d, want 0", len(j))
	}
}

func TestProject_WorkTimeConfiguration(t *testing.T) {
	t.Run("Default values when not set", func(t *testing.T) {
		project := Project{
			Name:   "Test Project",
			Type:   ProjectTypeProduct,
			Status: TaskStatusNotStarted,
		}

		if got := project.GetHoursPerDay(); got != DefaultHoursPerDay {
			t.Errorf("GetHoursPerDay() = %v, want %v", got, DefaultHoursPerDay)
		}
		if got := project.GetDaysPerWeek(); got != DefaultDaysPerWeek {
			t.Errorf("GetDaysPerWeek() = %v, want %v", got, DefaultDaysPerWeek)
		}
		if got := project.GetDaysPerMonth(); got != DefaultDaysPerMonth {
			t.Errorf("GetDaysPerMonth() = %v, want %v", got, DefaultDaysPerMonth)
		}
		if got := project.GetHoursPerMonth(); got != DefaultHoursPerMonth {
			t.Errorf("GetHoursPerMonth() = %v, want %v", got, DefaultHoursPerMonth)
		}
	})

	t.Run("Custom values when set", func(t *testing.T) {
		hoursPerDay := 6.0
		daysPerWeek := 4.0
		daysPerMonth := 16.0

		project := Project{
			Name:         "Test Project",
			Type:         ProjectTypeProduct,
			Status:       TaskStatusNotStarted,
			HoursPerDay:  &hoursPerDay,
			DaysPerWeek:  &daysPerWeek,
			DaysPerMonth: &daysPerMonth,
		}

		if got := project.GetHoursPerDay(); got != 6.0 {
			t.Errorf("GetHoursPerDay() = %v, want 6.0", got)
		}
		if got := project.GetDaysPerWeek(); got != 4.0 {
			t.Errorf("GetDaysPerWeek() = %v, want 4.0", got)
		}
		if got := project.GetDaysPerMonth(); got != 16.0 {
			t.Errorf("GetDaysPerMonth() = %v, want 16.0", got)
		}
		if got := project.GetHoursPerMonth(); got != 96.0 { // 6 * 16
			t.Errorf("GetHoursPerMonth() = %v, want 96.0", got)
		}
	})

	t.Run("Validation - negative values", func(t *testing.T) {
		negativeHours := -8.0
		project := Project{
			Name:        "Test Project",
			Type:        ProjectTypeProduct,
			Status:      TaskStatusNotStarted,
			HoursPerDay: &negativeHours,
		}

		err := project.BeforeSave(nil)
		if err == nil {
			t.Error("Expected error for negative hours per day, got nil")
		}
	})

	t.Run("Validation - zero values", func(t *testing.T) {
		zeroDays := 0.0
		project := Project{
			Name:         "Test Project",
			Type:         ProjectTypeProduct,
			Status:       TaskStatusNotStarted,
			DaysPerMonth: &zeroDays,
		}

		err := project.BeforeSave(nil)
		if err == nil {
			t.Error("Expected error for zero days per month, got nil")
		}
	})
}

func TestProject_ConversionsWithCustomSettings(t *testing.T) {
	t.Run("Default settings", func(t *testing.T) {
		project := Project{
			Name:            "Test Project",
			Type:            ProjectTypeProduct,
			Status:          TaskStatusNotStarted,
			EstimatedEffort: 160, // 160 hours
			ActualEffort:    80,  // 80 hours
		}

		// With default 8 hours/day, 20 days/month
		if got := project.EstimatedDays(); got != 20.0 {
			t.Errorf("EstimatedDays() = %v, want 20.0", got)
		}
		if got := project.EstimatedMonths(); got != 1.0 {
			t.Errorf("EstimatedMonths() = %v, want 1.0", got)
		}
		if got := project.ActualDays(); got != 10.0 {
			t.Errorf("ActualDays() = %v, want 10.0", got)
		}
		if got := project.ActualMonths(); got != 0.5 {
			t.Errorf("ActualMonths() = %v, want 0.5", got)
		}
	})

	t.Run("Custom settings - 6 hours/day, 16 days/month", func(t *testing.T) {
		hoursPerDay := 6.0
		daysPerMonth := 16.0

		project := Project{
			Name:            "Test Project",
			Type:            ProjectTypeProduct,
			Status:          TaskStatusNotStarted,
			EstimatedEffort: 96, // 96 hours
			ActualEffort:    48, // 48 hours
			HoursPerDay:     &hoursPerDay,
			DaysPerMonth:    &daysPerMonth,
		}

		// With 6 hours/day: 96 hours = 16 days
		if got := project.EstimatedDays(); got != 16.0 {
			t.Errorf("EstimatedDays() = %v, want 16.0", got)
		}
		// With 96 hours/month (6 * 16): 96 hours = 1 month
		if got := project.EstimatedMonths(); got != 1.0 {
			t.Errorf("EstimatedMonths() = %v, want 1.0", got)
		}
		// With 6 hours/day: 48 hours = 8 days
		if got := project.ActualDays(); got != 8.0 {
			t.Errorf("ActualDays() = %v, want 8.0", got)
		}
		// With 96 hours/month: 48 hours = 0.5 months
		if got := project.ActualMonths(); got != 0.5 {
			t.Errorf("ActualMonths() = %v, want 0.5", got)
		}
	})
}

