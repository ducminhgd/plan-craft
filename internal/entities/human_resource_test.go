package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHumanResourceTableName(t *testing.T) {
	humanResource := HumanResource{}
	assert.Equal(t, "human_resources", humanResource.TableName())
}

func TestHumanResourceIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status uint
		want   bool
	}{
		{"Active HumanResource", HumanResourceStatusActive, true},
		{"Inactive HumanResource", HumanResourceStatusInactive, false},
		{"Invalid status treated as inactive", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			humanResource := HumanResource{Status: tt.status}
			assert.Equal(t, tt.want, humanResource.IsActive())
		})
	}
}

func TestHumanResourceStatusChange(t *testing.T) {
	humanResource := HumanResource{Status: HumanResourceStatusInactive}

	// Test changing to active
	humanResource.Status = HumanResourceStatusActive
	assert.Equal(t, uint(HumanResourceStatusActive), humanResource.Status)
	assert.True(t, humanResource.IsActive())

	// Test changing to inactive
	humanResource.Status = HumanResourceStatusInactive
	assert.Equal(t, uint(HumanResourceStatusInactive), humanResource.Status)
	assert.False(t, humanResource.IsActive())
}

func TestHumanResourceValidateStatus(t *testing.T) {
	tests := []struct {
		name      string
		status    uint
		wantError error
	}{
		{"Valid: Active status", HumanResourceStatusActive, nil},
		{"Valid: Inactive status", HumanResourceStatusInactive, nil},
		{"Invalid: Status value 3", 3, ErrHumanResourceInvalidStatus},
		{"Invalid: High value", 999, ErrHumanResourceInvalidStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			humanResource := HumanResource{Status: tt.status}
			err := humanResource.validateStatus()
			if tt.wantError != nil {
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHumanResourceValidate(t *testing.T) {
	tests := []struct {
		name          string
		humanResource HumanResource
		wantError     error
	}{
		{
			name: "Valid HumanResource",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: HumanResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid with whitespace trimmed",
			humanResource: HumanResource{
				Name:   "  John Doe  ",
				Title:  "  Senior Developer  ",
				Level:  "  Senior  ",
				Status: HumanResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Empty name",
			humanResource: HumanResource{
				Name:   "",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceNameRequired,
		},
		{
			name: "Whitespace-only name",
			humanResource: HumanResource{
				Name:   "   ",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceNameRequired,
		},
		{
			name: "Empty title",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "",
				Level:  "Senior",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceTitleRequired,
		},
		{
			name: "Whitespace-only title",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "   ",
				Level:  "Senior",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceTitleRequired,
		},
		{
			name: "Empty level",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceLevelRequired,
		},
		{
			name: "Whitespace-only level",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "   ",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceLevelRequired,
		},
		{
			name: "Invalid status",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: 999,
			},
			wantError: ErrHumanResourceInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.humanResource.Validate()
			if tt.wantError != nil {
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
				// Verify whitespace was trimmed
				assert.NotContains(t, tt.humanResource.Name, "  ")
				assert.NotContains(t, tt.humanResource.Title, "  ")
				assert.NotContains(t, tt.humanResource.Level, "  ")
			}
		})
	}
}

func TestHumanResourceValidateTrimsWhitespace(t *testing.T) {
	humanResource := HumanResource{
		Name:   "  John Doe  ",
		Title:  "  Senior Developer  ",
		Level:  "  Senior  ",
		Status: HumanResourceStatusActive,
	}

	err := humanResource.Validate()
	assert.NoError(t, err)

	// Verify all fields were trimmed
	assert.Equal(t, "John Doe", humanResource.Name)
	assert.Equal(t, "Senior Developer", humanResource.Title)
	assert.Equal(t, "Senior", humanResource.Level)
}

func setupHumanResourceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the HumanResource table
	err = db.AutoMigrate(&HumanResource{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestHumanResourceBeforeCreate(t *testing.T) {
	db := setupHumanResourceTestDB(t)

	tests := []struct {
		name          string
		humanResource HumanResource
		wantError     error
	}{
		{
			name: "Valid HumanResource creation",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: HumanResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Default status set to active when invalid",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: 999, // Invalid status
			},
			wantError: nil, // Should succeed after setting default
		},
		{
			name: "Default status set to active when zero",
			humanResource: HumanResource{
				Name:  "John Doe",
				Title: "Senior Developer",
				Level: "Senior",
				// Status not set (will be 0)
			},
			wantError: nil, // Should succeed after setting default
		},
		{
			name: "Validation fails - empty name",
			humanResource: HumanResource{
				Name:   "",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceNameRequired,
		},
		{
			name: "Validation fails - empty title",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "",
				Level:  "Senior",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceTitleRequired,
		},
		{
			name: "Validation fails - empty level",
			humanResource: HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "",
				Status: HumanResourceStatusActive,
			},
			wantError: ErrHumanResourceLevelRequired,
		},
		{
			name: "Whitespace trimmed during creation",
			humanResource: HumanResource{
				Name:   "  John Doe  ",
				Title:  "  Senior Developer  ",
				Level:  "  Senior  ",
				Status: HumanResourceStatusActive,
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.Create(&tt.humanResource)

			if tt.wantError != nil {
				assert.Error(t, result.Error)
				assert.Contains(t, result.Error.Error(), tt.wantError.Error())
			} else {
				assert.NoError(t, result.Error)
				assert.NotZero(t, tt.humanResource.ID)
				assert.NotZero(t, tt.humanResource.CreatedAt)
				assert.NotZero(t, tt.humanResource.UpdatedAt)

				// Verify status was set to active if it was invalid
				if tt.name == "Default status set to active when invalid" || tt.name == "Default status set to active when zero" {
					assert.Equal(t, uint(HumanResourceStatusActive), tt.humanResource.Status)
				}

				// Verify whitespace was trimmed
				assert.NotContains(t, tt.humanResource.Name, "  ")
				assert.NotContains(t, tt.humanResource.Title, "  ")
				assert.NotContains(t, tt.humanResource.Level, "  ")
			}
		})
	}
}

func TestHumanResourceBeforeUpdate(t *testing.T) {
	db := setupHumanResourceTestDB(t)

	// Create a valid HumanResource first
	humanResource := HumanResource{
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: HumanResourceStatusActive,
	}
	result := db.Create(&humanResource)
	assert.NoError(t, result.Error)

	tests := []struct {
		name       string
		updateFunc func(*HumanResource)
		wantError  error
	}{
		{
			name: "Valid update",
			updateFunc: func(h *HumanResource) {
				h.Name = "Jane Smith"
				h.Title = "Lead Developer"
			},
			wantError: nil,
		},
		{
			name: "Update with whitespace trimming",
			updateFunc: func(h *HumanResource) {
				h.Name = "  Updated Name  "
				h.Title = "  Updated Title  "
				h.Level = "  Updated Level  "
			},
			wantError: nil,
		},
		{
			name: "Update fails - empty name",
			updateFunc: func(h *HumanResource) {
				h.Name = ""
			},
			wantError: ErrHumanResourceNameRequired,
		},
		{
			name: "Update fails - empty title",
			updateFunc: func(h *HumanResource) {
				h.Title = ""
			},
			wantError: ErrHumanResourceTitleRequired,
		},
		{
			name: "Update fails - empty level",
			updateFunc: func(h *HumanResource) {
				h.Level = ""
			},
			wantError: ErrHumanResourceLevelRequired,
		},
		{
			name: "Update fails - invalid status",
			updateFunc: func(h *HumanResource) {
				h.Status = 999
			},
			wantError: ErrHumanResourceInvalidStatus,
		},
		{
			name: "Update status to inactive",
			updateFunc: func(h *HumanResource) {
				h.Status = HumanResourceStatusInactive
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reload HumanResource from database
			var testHR HumanResource
			db.First(&testHR, humanResource.ID)

			// Apply update
			tt.updateFunc(&testHR)

			// Attempt to save
			result := db.Save(&testHR)

			if tt.wantError != nil {
				assert.Error(t, result.Error)
				assert.Contains(t, result.Error.Error(), tt.wantError.Error())
			} else {
				assert.NoError(t, result.Error)

				// Verify whitespace was trimmed
				assert.NotContains(t, testHR.Name, "  ")
				assert.NotContains(t, testHR.Title, "  ")
				assert.NotContains(t, testHR.Level, "  ")
			}
		})
	}
}

func TestHumanResourceCRUDOperations(t *testing.T) {
	db := setupHumanResourceTestDB(t)

	// Create
	humanResource := HumanResource{
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: HumanResourceStatusActive,
	}
	result := db.Create(&humanResource)
	assert.NoError(t, result.Error)
	assert.NotZero(t, humanResource.ID)

	// Read
	var retrieved HumanResource
	result = db.First(&retrieved, humanResource.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, humanResource.Name, retrieved.Name)
	assert.Equal(t, humanResource.Title, retrieved.Title)
	assert.Equal(t, humanResource.Level, retrieved.Level)
	assert.Equal(t, humanResource.Status, retrieved.Status)

	// Update
	retrieved.Name = "Jane Smith"
	retrieved.Title = "Lead Developer"
	result = db.Save(&retrieved)
	assert.NoError(t, result.Error)

	// Verify update
	var updated HumanResource
	result = db.First(&updated, humanResource.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Jane Smith", updated.Name)
	assert.Equal(t, "Lead Developer", updated.Title)

	// Delete
	result = db.Delete(&updated)
	assert.NoError(t, result.Error)

	// Verify deletion
	var deleted HumanResource
	result = db.First(&deleted, humanResource.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestHumanResourceQueryOperations(t *testing.T) {
	db := setupHumanResourceTestDB(t)

	// Create multiple HRs
	hrs := []HumanResource{
		{Name: "John Doe", Title: "Senior Developer", Level: "Senior", Status: HumanResourceStatusActive},
		{Name: "Jane Smith", Title: "Lead Developer", Level: "Lead", Status: HumanResourceStatusActive},
		{Name: "Bob Johnson", Title: "Junior Developer", Level: "Junior", Status: HumanResourceStatusInactive},
	}

	for _, h := range hrs {
		db.Create(&h)
	}

	// Query all HRs
	var allHRs []HumanResource
	result := db.Find(&allHRs)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(allHRs), 3)

	// Query active HRs only
	var activeHRs []HumanResource
	result = db.Where("status = ?", HumanResourceStatusActive).Find(&activeHRs)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(activeHRs), 2)

	// Query by level
	var humanResource HumanResource
	result = db.Where("level = ?", "Senior").First(&humanResource)
	assert.NoError(t, result.Error)
	assert.Equal(t, "John Doe", humanResource.Name)
}

func BenchmarkHRValidate(b *testing.B) {
	humanResource := HumanResource{
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: HumanResourceStatusActive,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = humanResource.Validate()
	}
}

func BenchmarkHRCreate(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&HumanResource{}) //nolint:errcheck

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		humanResource := HumanResource{
			Name:   "John Doe",
			Title:  "Senior Developer",
			Level:  "Senior",
			Status: HumanResourceStatusActive,
		}
		db.Create(&humanResource)
	}
}
