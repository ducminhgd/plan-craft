package entities

import (
	"testing"
	"time"
)

func TestWeekdayArrayValue(t *testing.T) {
	tests := []struct {
		name    string
		input   WeekdayArray
		want    string
		wantErr bool
	}{
		{
			name:    "Empty array",
			input:   WeekdayArray{},
			want:    "[]",
			wantErr: false,
		},
		{
			name:    "Nil array",
			input:   nil,
			want:    "[]",
			wantErr: false,
		},
		{
			name:    "Single day",
			input:   WeekdayArray{time.Monday},
			want:    "[1]",
			wantErr: false,
		},
		{
			name:    "Multiple days",
			input:   WeekdayArray{time.Monday, time.Tuesday, time.Wednesday},
			want:    "[1,2,3]",
			wantErr: false,
		},
		{
			name:    "Full work week",
			input:   WeekdayArray{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			want:    "[1,2,3,4,5]",
			wantErr: false,
		},
		{
			name:    "Weekend",
			input:   WeekdayArray{time.Saturday, time.Sunday},
			want:    "[6,0]",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("WeekdayArray.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WeekdayArray.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekdayArrayScan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    WeekdayArray
		wantErr bool
	}{
		{
			name:    "Nil value",
			input:   nil,
			want:    WeekdayArray{},
			wantErr: false,
		},
		{
			name:    "Empty string",
			input:   "",
			want:    WeekdayArray{},
			wantErr: false,
		},
		{
			name:    "Empty JSON array string",
			input:   "[]",
			want:    WeekdayArray{},
			wantErr: false,
		},
		{
			name:    "Empty JSON array bytes",
			input:   []byte("[]"),
			want:    WeekdayArray{},
			wantErr: false,
		},
		{
			name:    "Single day string",
			input:   "[1]",
			want:    WeekdayArray{time.Monday},
			wantErr: false,
		},
		{
			name:    "Multiple days string",
			input:   "[1,2,3]",
			want:    WeekdayArray{time.Monday, time.Tuesday, time.Wednesday},
			wantErr: false,
		},
		{
			name:    "Full work week bytes",
			input:   []byte("[1,2,3,4,5]"),
			want:    WeekdayArray{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			wantErr: false,
		},
		{
			name:    "Weekend",
			input:   "[6,0]",
			want:    WeekdayArray{time.Saturday, time.Sunday},
			wantErr: false,
		},
		{
			name:    "Invalid type",
			input:   123,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			input:   "invalid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w WeekdayArray
			err := w.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("WeekdayArray.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(w) != len(tt.want) {
					t.Errorf("WeekdayArray.Scan() length = %d, want %d", len(w), len(tt.want))
					return
				}
				for i, day := range w {
					if day != tt.want[i] {
						t.Errorf("WeekdayArray.Scan()[%d] = %v, want %v", i, day, tt.want[i])
					}
				}
			}
		})
	}
}

func TestWeekdayArrayRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input WeekdayArray
	}{
		{
			name:  "Empty array",
			input: WeekdayArray{},
		},
		{
			name:  "Single day",
			input: WeekdayArray{time.Monday},
		},
		{
			name:  "Full work week",
			input: WeekdayArray{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		},
		{
			name:  "All days",
			input: WeekdayArray{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to database format
			val, err := tt.input.Value()
			if err != nil {
				t.Errorf("WeekdayArray.Value() error = %v", err)
				return
			}

			// Convert back from database format
			var result WeekdayArray
			err = result.Scan(val)
			if err != nil {
				t.Errorf("WeekdayArray.Scan() error = %v", err)
				return
			}

			// Compare
			if len(result) != len(tt.input) {
				t.Errorf("Round trip length = %d, want %d", len(result), len(tt.input))
				return
			}
			for i, day := range result {
				if day != tt.input[i] {
					t.Errorf("Round trip[%d] = %v, want %v", i, day, tt.input[i])
				}
			}
		})
	}
}

func TestDefaultWorkingDays(t *testing.T) {
	days := DefaultWorkingDays()

	// Should have 5 days
	if len(days) != 5 {
		t.Errorf("DefaultWorkingDays() length = %d, want 5", len(days))
		return
	}

	// Should be Monday through Friday
	expected := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}
	for i, day := range days {
		if day != expected[i] {
			t.Errorf("DefaultWorkingDays()[%d] = %v, want %v", i, day, expected[i])
		}
	}
}

func TestIsValidProjectType(t *testing.T) {
	tests := []struct {
		name string
		pt   ProjectType
		want bool
	}{
		{"Valid: product", ProjectTypeProduct, true},
		{"Valid: service", ProjectTypeService, true},
		{"Valid: internal", ProjectTypeInternal, true},
		{"Valid: consulting", ProjectTypeConsulting, true},
		{"Valid: research", ProjectTypeResearch, true},
		{"Valid: maintenance", ProjectTypeMaintenance, true},
		{"Invalid: empty", ProjectType(""), false},
		{"Invalid: unknown", ProjectType("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidProjectType(tt.pt); got != tt.want {
				t.Errorf("IsValidProjectType(%v) = %v, want %v", tt.pt, got, tt.want)
			}
		})
	}
}

func TestIsValidTaskStatus(t *testing.T) {
	tests := []struct {
		name string
		ts   TaskStatus
		want bool
	}{
		{"Valid: not_started", TaskStatusNotStarted, true},
		{"Valid: in_progress", TaskStatusInProgress, true},
		{"Valid: on_hold", TaskStatusOnHold, true},
		{"Valid: completed", TaskStatusCompleted, true},
		{"Valid: cancelled", TaskStatusCancelled, true},
		{"Invalid: empty", TaskStatus(""), false},
		{"Invalid: unknown", TaskStatus("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidTaskStatus(tt.ts); got != tt.want {
				t.Errorf("IsValidTaskStatus(%v) = %v, want %v", tt.ts, got, tt.want)
			}
		})
	}
}

func TestIsValidDependencyType(t *testing.T) {
	tests := []struct {
		name string
		dt   DependencyType
		want bool
	}{
		{"Valid: finish_to_start", DependencyFinishToStart, true},
		{"Valid: start_to_start", DependencyStartToStart, true},
		{"Valid: finish_to_finish", DependencyFinishToFinish, true},
		{"Valid: start_to_finish", DependencyStartToFinish, true},
		{"Invalid: empty", DependencyType(""), false},
		{"Invalid: unknown", DependencyType("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDependencyType(tt.dt); got != tt.want {
				t.Errorf("IsValidDependencyType(%v) = %v, want %v", tt.dt, got, tt.want)
			}
		})
	}
}

func TestIsValidPriority(t *testing.T) {
	tests := []struct {
		name string
		p    Priority
		want bool
	}{
		{"Valid: low", PriorityLow, true},
		{"Valid: medium", PriorityMedium, true},
		{"Valid: high", PriorityHigh, true},
		{"Valid: critical", PriorityCritical, true},
		{"Invalid: empty", Priority(""), false},
		{"Invalid: unknown", Priority("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPriority(tt.p); got != tt.want {
				t.Errorf("IsValidPriority(%v) = %v, want %v", tt.p, got, tt.want)
			}
		})
	}
}

func TestIsValidCostType(t *testing.T) {
	tests := []struct {
		name string
		ct   CostType
		want bool
	}{
		{"Valid: labor", CostTypeLabor, true},
		{"Valid: material", CostTypeMaterial, true},
		{"Valid: equipment", CostTypeEquipment, true},
		{"Valid: overhead", CostTypeOverhead, true},
		{"Valid: infrastructure", CostTypeInfrastructure, true},
		{"Valid: service", CostTypeService, true},
		{"Valid: other", CostTypeOther, true},
		{"Invalid: empty", CostType(""), false},
		{"Invalid: unknown", CostType("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCostType(tt.ct); got != tt.want {
				t.Errorf("IsValidCostType(%v) = %v, want %v", tt.ct, got, tt.want)
			}
		})
	}
}

func TestIsValidRateType(t *testing.T) {
	tests := []struct {
		name string
		rt   RateType
		want bool
	}{
		{"Valid: hourly", RateTypeHourly, true},
		{"Valid: daily", RateTypeDaily, true},
		{"Valid: monthly", RateTypeMonthly, true},
		{"Valid: fixed", RateTypeFixed, true},
		{"Invalid: empty", RateType(""), false},
		{"Invalid: unknown", RateType("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRateType(tt.rt); got != tt.want {
				t.Errorf("IsValidRateType(%v) = %v, want %v", tt.rt, got, tt.want)
			}
		})
	}
}

func TestUnitConversions(t *testing.T) {
	t.Run("HoursToDays", func(t *testing.T) {
		tests := []struct {
			hours float64
			want  float64
		}{
			{0, 0},
			{8, 1},
			{16, 2},
			{40, 5},
			{160, 20},
		}
		for _, tt := range tests {
			if got := HoursToDays(tt.hours); got != tt.want {
				t.Errorf("HoursToDays(%v) = %v, want %v", tt.hours, got, tt.want)
			}
		}
	})

	t.Run("DaysToHours", func(t *testing.T) {
		tests := []struct {
			days float64
			want float64
		}{
			{0, 0},
			{1, 8},
			{2, 16},
			{5, 40},
			{20, 160},
		}
		for _, tt := range tests {
			if got := DaysToHours(tt.days); got != tt.want {
				t.Errorf("DaysToHours(%v) = %v, want %v", tt.days, got, tt.want)
			}
		}
	})

	t.Run("HoursToMonths", func(t *testing.T) {
		tests := []struct {
			hours float64
			want  float64
		}{
			{0, 0},
			{160, 1},
			{320, 2},
			{80, 0.5},
		}
		for _, tt := range tests {
			if got := HoursToMonths(tt.hours); got != tt.want {
				t.Errorf("HoursToMonths(%v) = %v, want %v", tt.hours, got, tt.want)
			}
		}
	})

	t.Run("MonthsToHours", func(t *testing.T) {
		tests := []struct {
			months float64
			want   float64
		}{
			{0, 0},
			{1, 160},
			{2, 320},
			{0.5, 80},
		}
		for _, tt := range tests {
			if got := MonthsToHours(tt.months); got != tt.want {
				t.Errorf("MonthsToHours(%v) = %v, want %v", tt.months, got, tt.want)
			}
		}
	})

	t.Run("DaysToMonths", func(t *testing.T) {
		tests := []struct {
			days float64
			want float64
		}{
			{0, 0},
			{20, 1},
			{40, 2},
			{10, 0.5},
		}
		for _, tt := range tests {
			if got := DaysToMonths(tt.days); got != tt.want {
				t.Errorf("DaysToMonths(%v) = %v, want %v", tt.days, got, tt.want)
			}
		}
	})

	t.Run("MonthsToDays", func(t *testing.T) {
		tests := []struct {
			months float64
			want   float64
		}{
			{0, 0},
			{1, 20},
			{2, 40},
			{0.5, 10},
		}
		for _, tt := range tests {
			if got := MonthsToDays(tt.months); got != tt.want {
				t.Errorf("MonthsToDays(%v) = %v, want %v", tt.months, got, tt.want)
			}
		}
	})
}

func TestCustomUnitConversions(t *testing.T) {
	t.Run("HoursToDaysCustom", func(t *testing.T) {
		tests := []struct {
			hours       float64
			hoursPerDay float64
			want        float64
		}{
			{48, 6, 8},      // 48 hours with 6 hours/day = 8 days
			{40, 10, 4},     // 40 hours with 10 hours/day = 4 days
			{32, 8, 4},      // 32 hours with 8 hours/day = 4 days
			{24, 0, 3},      // 24 hours with 0 hours/day (use default 8) = 3 days
			{24, -1, 3},     // 24 hours with negative hours/day (use default 8) = 3 days
		}
		for _, tt := range tests {
			if got := HoursToDaysCustom(tt.hours, tt.hoursPerDay); got != tt.want {
				t.Errorf("HoursToDaysCustom(%v, %v) = %v, want %v", tt.hours, tt.hoursPerDay, got, tt.want)
			}
		}
	})

	t.Run("DaysToHoursCustom", func(t *testing.T) {
		tests := []struct {
			days        float64
			hoursPerDay float64
			want        float64
		}{
			{5, 6, 30},   // 5 days with 6 hours/day = 30 hours
			{4, 10, 40},  // 4 days with 10 hours/day = 40 hours
			{10, 0, 80},  // 10 days with 0 hours/day (use default 8) = 80 hours
		}
		for _, tt := range tests {
			if got := DaysToHoursCustom(tt.days, tt.hoursPerDay); got != tt.want {
				t.Errorf("DaysToHoursCustom(%v, %v) = %v, want %v", tt.days, tt.hoursPerDay, got, tt.want)
			}
		}
	})

	t.Run("HoursToMonthsCustom", func(t *testing.T) {
		tests := []struct {
			hours        float64
			hoursPerDay  float64
			daysPerMonth float64
			want         float64
		}{
			{96, 6, 16, 1},    // 96 hours with 6 hours/day and 16 days/month = 1 month
			{120, 10, 12, 1},  // 120 hours with 10 hours/day and 12 days/month = 1 month
			{160, 8, 20, 1},   // 160 hours with 8 hours/day and 20 days/month = 1 month
			{80, 0, 0, 0.5},   // 80 hours with defaults = 0.5 months
		}
		for _, tt := range tests {
			if got := HoursToMonthsCustom(tt.hours, tt.hoursPerDay, tt.daysPerMonth); got != tt.want {
				t.Errorf("HoursToMonthsCustom(%v, %v, %v) = %v, want %v",
					tt.hours, tt.hoursPerDay, tt.daysPerMonth, got, tt.want)
			}
		}
	})

	t.Run("MonthsToHoursCustom", func(t *testing.T) {
		tests := []struct {
			months       float64
			hoursPerDay  float64
			daysPerMonth float64
			want         float64
		}{
			{1, 6, 16, 96},   // 1 month with 6 hours/day and 16 days/month = 96 hours
			{2, 10, 12, 240}, // 2 months with 10 hours/day and 12 days/month = 240 hours
			{0.5, 8, 20, 80}, // 0.5 months with 8 hours/day and 20 days/month = 80 hours
		}
		for _, tt := range tests {
			if got := MonthsToHoursCustom(tt.months, tt.hoursPerDay, tt.daysPerMonth); got != tt.want {
				t.Errorf("MonthsToHoursCustom(%v, %v, %v) = %v, want %v",
					tt.months, tt.hoursPerDay, tt.daysPerMonth, got, tt.want)
			}
		}
	})

	t.Run("DaysToMonthsCustom", func(t *testing.T) {
		tests := []struct {
			days         float64
			daysPerMonth float64
			want         float64
		}{
			{16, 16, 1},   // 16 days with 16 days/month = 1 month
			{12, 12, 1},   // 12 days with 12 days/month = 1 month
			{10, 20, 0.5}, // 10 days with 20 days/month = 0.5 months
			{20, 0, 1},    // 20 days with 0 days/month (use default 20) = 1 month
		}
		for _, tt := range tests {
			if got := DaysToMonthsCustom(tt.days, tt.daysPerMonth); got != tt.want {
				t.Errorf("DaysToMonthsCustom(%v, %v) = %v, want %v", tt.days, tt.daysPerMonth, got, tt.want)
			}
		}
	})

	t.Run("MonthsToDaysCustom", func(t *testing.T) {
		tests := []struct {
			months       float64
			daysPerMonth float64
			want         float64
		}{
			{1, 16, 16},   // 1 month with 16 days/month = 16 days
			{2, 12, 24},   // 2 months with 12 days/month = 24 days
			{0.5, 20, 10}, // 0.5 months with 20 days/month = 10 days
			{1, 0, 20},    // 1 month with 0 days/month (use default 20) = 20 days
		}
		for _, tt := range tests {
			if got := MonthsToDaysCustom(tt.months, tt.daysPerMonth); got != tt.want {
				t.Errorf("MonthsToDaysCustom(%v, %v) = %v, want %v", tt.months, tt.daysPerMonth, got, tt.want)
			}
		}
	})
}

