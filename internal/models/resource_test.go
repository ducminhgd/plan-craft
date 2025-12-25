package models

import "testing"

func TestResource_TableName(t *testing.T) {
	resource := Resource{}
	if resource.TableName() != "resources" {
		t.Errorf("Expected table name 'resources', got '%s'", resource.TableName())
	}
}

func TestResource_BeforeSave_ValidResource(t *testing.T) {
	resource := Resource{
		Name:                 "John Doe",
		Role:                 "Backend Developer",
		DefaultHoursPerDay:   8,
		DefaultDaysPerWeek:   5,
		DefaultDaysPerMonth:  20,
		DefaultHourlyRate:    50,
		DefaultDailyRate:     400,
		DefaultMonthlyRate:   8000,
	}

	err := resource.BeforeSave(nil)
	if err != nil {
		t.Errorf("Expected no error for valid resource, got: %v", err)
	}
}

func TestResource_BeforeSave_InvalidCapacity(t *testing.T) {
	tests := []struct {
		name             string
		hoursPerDay      float64
		daysPerWeek      float64
		daysPerMonth     float64
		wantErr          bool
	}{
		{"Valid capacity", 8, 5, 20, false},
		{"Invalid hours per day - negative", -1, 5, 20, true},
		{"Invalid hours per day - over 24", 25, 5, 20, true},
		{"Invalid days per week - negative", 8, -1, 20, true},
		{"Invalid days per week - over 7", 8, 8, 20, true},
		{"Invalid days per month - negative", 8, 5, -1, true},
		{"Invalid days per month - over 31", 8, 5, 32, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := Resource{
				Name:                "Test Resource",
				Role:                "Developer",
				DefaultHoursPerDay:  tt.hoursPerDay,
				DefaultDaysPerWeek:  tt.daysPerWeek,
				DefaultDaysPerMonth: tt.daysPerMonth,
			}

			err := resource.BeforeSave(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeforeSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResource_BeforeSave_InvalidRates(t *testing.T) {
	tests := []struct {
		name         string
		hourlyRate   float64
		dailyRate    float64
		monthlyRate  float64
		wantErr      bool
	}{
		{"Valid rates", 50, 400, 8000, false},
		{"Negative hourly rate", -1, 400, 8000, true},
		{"Negative daily rate", 50, -1, 8000, true},
		{"Negative monthly rate", 50, 400, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := Resource{
				Name:               "Test Resource",
				Role:               "Developer",
				DefaultHourlyRate:  tt.hourlyRate,
				DefaultDailyRate:   tt.dailyRate,
				DefaultMonthlyRate: tt.monthlyRate,
			}

			err := resource.BeforeSave(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeforeSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResource_MonthlyCapacityHours(t *testing.T) {
	resource := Resource{
		DefaultHoursPerDay:  8,
		DefaultDaysPerMonth: 20,
	}

	expected := 160.0 // 8 hours * 20 days
	if got := resource.MonthlyCapacityHours(); got != expected {
		t.Errorf("MonthlyCapacityHours() = %v, want %v", got, expected)
	}
}

func TestProjectRole_TableName(t *testing.T) {
	pr := ProjectRole{}
	if pr.TableName() != "project_roles" {
		t.Errorf("Expected table name 'project_roles', got '%s'", pr.TableName())
	}
}

func TestProjectRole_GetEffectiveHoursPerDay(t *testing.T) {
	resource := &Resource{
		DefaultHoursPerDay: 8,
	}

	tests := []struct {
		name        string
		hoursPerDay *float64
		resource    *Resource
		want        float64
	}{
		{"Use project-specific value", floatPtr(6), resource, 6},
		{"Use resource default", nil, resource, 8},
		{"Use fallback", nil, nil, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := ProjectRole{
				HoursPerDay: tt.hoursPerDay,
			}
			if got := pr.GetEffectiveHoursPerDay(tt.resource); got != tt.want {
				t.Errorf("GetEffectiveHoursPerDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskAssignment_TableName(t *testing.T) {
	ta := TaskAssignment{}
	if ta.TableName() != "task_assignments" {
		t.Errorf("Expected table name 'task_assignments', got '%s'", ta.TableName())
	}
}

func TestTaskAssignment_EstimatedHours(t *testing.T) {
	ta := TaskAssignment{
		EstimatedManDays: 5,
	}

	expected := 40.0 // 5 days * 8 hours
	if got := ta.EstimatedHours(); got != expected {
		t.Errorf("EstimatedHours() = %v, want %v", got, expected)
	}
}

func TestTaskAssignment_ActualHours(t *testing.T) {
	ta := TaskAssignment{
		ActualManDays: 3,
	}

	expected := 24.0 // 3 days * 8 hours
	if got := ta.ActualHours(); got != expected {
		t.Errorf("ActualHours() = %v, want %v", got, expected)
	}
}

// Helper function
func floatPtr(f float64) *float64 {
	return &f
}
