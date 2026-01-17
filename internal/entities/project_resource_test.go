package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProjectResourceTableName(t *testing.T) {
	pr := ProjectResource{}
	assert.Equal(t, "project_resources", pr.TableName())
}

func TestProjectResourceIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status uint
		want   bool
	}{
		{"Active project resource", ProjectResourceStatusActive, true},
		{"Inactive project resource", ProjectResourceStatusInactive, false},
		{"Invalid status treated as inactive", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := ProjectResource{Status: tt.status}
			assert.Equal(t, tt.want, pr.IsActive())
		})
	}
}

func TestProjectResourceStatusChange(t *testing.T) {
	pr := ProjectResource{Status: ProjectResourceStatusInactive}

	// Test changing to active
	pr.Status = ProjectResourceStatusActive
	assert.Equal(t, uint(ProjectResourceStatusActive), pr.Status)
	assert.True(t, pr.IsActive())

	// Test changing to inactive
	pr.Status = ProjectResourceStatusInactive
	assert.Equal(t, uint(ProjectResourceStatusInactive), pr.Status)
	assert.False(t, pr.IsActive())
}

func TestProjectResourceValidateStatus(t *testing.T) {
	tests := []struct {
		name      string
		status    uint
		wantError error
	}{
		{"Valid: Active status", ProjectResourceStatusActive, nil},
		{"Valid: Inactive status", ProjectResourceStatusInactive, nil},
		{"Invalid: Status value 3", 3, ErrProjectResourceInvalidStatus},
		{"Invalid: High value", 999, ErrProjectResourceInvalidStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := ProjectResource{Status: tt.status}
			err := pr.validateStatus()
			if tt.wantError != nil {
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectResourceValidate(t *testing.T) {
	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)
	pastDate := now.AddDate(0, -1, 0)

	tests := []struct {
		name      string
		pr        ProjectResource
		wantError error
	}{
		{
			name: "Valid project resource",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      100,
				Status:          ProjectResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid with role and notes",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Role:            "Developer",
				Allocation:      50,
				Notes:           "Part-time allocation",
				Status:          ProjectResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Missing project ID",
			pr: ProjectResource{
				ProjectID:       0,
				HumanResourceID: 1,
				Allocation:      100,
				Status:          ProjectResourceStatusActive,
			},
			wantError: ErrProjectResourceInvalidProjectID,
		},
		{
			name: "Missing human resource ID",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 0,
				Allocation:      100,
				Status:          ProjectResourceStatusActive,
			},
			wantError: ErrProjectResourceInvalidHumanResourceID,
		},
		{
			name: "Invalid allocation - negative",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      -10,
				Status:          ProjectResourceStatusActive,
			},
			wantError: ErrProjectResourceInvalidAllocation,
		},
		{
			name: "Invalid allocation - over 100",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      150,
				Status:          ProjectResourceStatusActive,
			},
			wantError: ErrProjectResourceInvalidAllocation,
		},
		{
			name: "Valid allocation - zero",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      0,
				Status:          ProjectResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid allocation - exactly 100",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      100,
				Status:          ProjectResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid allocation - fractional",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      33.33,
				Status:          ProjectResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Invalid status",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      100,
				Status:          999,
			},
			wantError: ErrProjectResourceInvalidStatus,
		},
		{
			name: "Valid with dates",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      100,
				StartDate:       &now,
				EndDate:         &futureDate,
				Status:          ProjectResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Invalid dates - end before start",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      100,
				StartDate:       &now,
				EndDate:         &pastDate,
				Status:          ProjectResourceStatusActive,
			},
			wantError: ErrProjectResourceInvalidDates,
		},
		{
			name: "Valid with only start date",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      100,
				StartDate:       &now,
				Status:          ProjectResourceStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid with only end date",
			pr: ProjectResource{
				ProjectID:       1,
				HumanResourceID: 1,
				Allocation:      100,
				EndDate:         &futureDate,
				Status:          ProjectResourceStatusActive,
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pr.Validate()
			if tt.wantError != nil {
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func setupProjectResourceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all required tables
	err = db.AutoMigrate(&Client{}, &HumanResource{}, &Project{}, &ProjectResource{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func createTestClientForPR(t *testing.T, db *gorm.DB) *Client {
	client := &Client{
		Name:   "Test Client",
		Email:  "test@client.com",
		Status: ClientStatusActive,
	}
	result := db.Create(client)
	assert.NoError(t, result.Error)
	return client
}

func createTestProjectForPR(t *testing.T, db *gorm.DB, clientID uint) *Project {
	project := &Project{
		Name:     "Test Project",
		ClientID: clientID,
		Status:   ProjectStatusActive,
	}
	result := db.Create(project)
	assert.NoError(t, result.Error)
	return project
}

func createTestHumanResource(t *testing.T, db *gorm.DB) *HumanResource {
	hr := &HumanResource{
		Name:   "John Doe",
		Title:  "Software Engineer",
		Level:  "Senior",
		Status: HumanResourceStatusActive,
	}
	result := db.Create(hr)
	assert.NoError(t, result.Error)
	return hr
}

func TestProjectResourceBeforeCreate(t *testing.T) {
	db := setupProjectResourceTestDB(t)
	client := createTestClientForPR(t, db)
	project := createTestProjectForPR(t, db, client.ID)

	// Create multiple human resources to avoid unique constraint violations
	hrCounter := 0
	createUniqueHR := func() *HumanResource {
		hrCounter++
		hr := &HumanResource{
			Name:   "Test HR " + string(rune('A'+hrCounter)),
			Title:  "Engineer",
			Level:  "Senior",
			Status: HumanResourceStatusActive,
		}
		db.Create(hr)
		return hr
	}

	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)
	pastDate := now.AddDate(0, -1, 0)

	t.Run("Valid project resource creation", func(t *testing.T) {
		hr := createUniqueHR()
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: hr.ID,
			Role:            "Developer",
			Allocation:      100,
			Status:          ProjectResourceStatusActive,
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.NotZero(t, pr.ID)
		assert.NotZero(t, pr.CreatedAt)
		assert.NotZero(t, pr.UpdatedAt)
	})

	t.Run("Default status set to active when invalid", func(t *testing.T) {
		hr := createUniqueHR()
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: hr.ID,
			Allocation:      50,
			Status:          999, // Invalid status
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.Equal(t, uint(ProjectResourceStatusActive), pr.Status)
	})

	t.Run("Default status set to active when zero", func(t *testing.T) {
		hr := createUniqueHR()
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: hr.ID,
			Allocation:      75,
			// Status not set (will be 0)
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.Equal(t, uint(ProjectResourceStatusActive), pr.Status)
	})

	t.Run("Default allocation set to 100 when zero", func(t *testing.T) {
		hr := createUniqueHR()
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: hr.ID,
			// Allocation not set (will be 0)
			Status: ProjectResourceStatusActive,
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.Equal(t, float64(100), pr.Allocation)
	})

	t.Run("Validation fails - missing project ID", func(t *testing.T) {
		hr := createUniqueHR()
		pr := ProjectResource{
			ProjectID:       0,
			HumanResourceID: hr.ID,
			Allocation:      100,
			Status:          ProjectResourceStatusActive,
		}
		result := db.Create(&pr)
		assert.Error(t, result.Error)
		assert.Contains(t, result.Error.Error(), ErrProjectResourceInvalidProjectID.Error())
	})

	t.Run("Validation fails - missing human resource ID", func(t *testing.T) {
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: 0,
			Allocation:      100,
			Status:          ProjectResourceStatusActive,
		}
		result := db.Create(&pr)
		assert.Error(t, result.Error)
		assert.Contains(t, result.Error.Error(), ErrProjectResourceInvalidHumanResourceID.Error())
	})

	t.Run("Validation fails - invalid allocation", func(t *testing.T) {
		hr := createUniqueHR()
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: hr.ID,
			Allocation:      150,
			Status:          ProjectResourceStatusActive,
		}
		result := db.Create(&pr)
		assert.Error(t, result.Error)
		assert.Contains(t, result.Error.Error(), ErrProjectResourceInvalidAllocation.Error())
	})

	t.Run("Validation fails - invalid dates", func(t *testing.T) {
		hr := createUniqueHR()
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: hr.ID,
			Allocation:      100,
			StartDate:       &now,
			EndDate:         &pastDate,
			Status:          ProjectResourceStatusActive,
		}
		result := db.Create(&pr)
		assert.Error(t, result.Error)
		assert.Contains(t, result.Error.Error(), ErrProjectResourceInvalidDates.Error())
	})

	t.Run("Valid with dates", func(t *testing.T) {
		hr := createUniqueHR()
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: hr.ID,
			Allocation:      100,
			StartDate:       &now,
			EndDate:         &futureDate,
			Status:          ProjectResourceStatusActive,
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.NotZero(t, pr.ID)
	})
}

func TestProjectResourceBeforeUpdate(t *testing.T) {
	db := setupProjectResourceTestDB(t)
	client := createTestClientForPR(t, db)
	project := createTestProjectForPR(t, db, client.ID)
	hr := createTestHumanResource(t, db)

	now := time.Now()
	pastDate := now.AddDate(0, -1, 0)

	// Create a valid project resource first
	pr := ProjectResource{
		ProjectID:       project.ID,
		HumanResourceID: hr.ID,
		Role:            "Developer",
		Allocation:      100,
		Status:          ProjectResourceStatusActive,
	}
	result := db.Create(&pr)
	assert.NoError(t, result.Error)

	tests := []struct {
		name       string
		updateFunc func(*ProjectResource)
		wantError  error
	}{
		{
			name: "Valid update",
			updateFunc: func(p *ProjectResource) {
				p.Role = "Tech Lead"
				p.Allocation = 80
			},
			wantError: nil,
		},
		{
			name: "Update allocation to 50%",
			updateFunc: func(p *ProjectResource) {
				p.Allocation = 50
			},
			wantError: nil,
		},
		{
			name: "Update fails - invalid allocation over 100",
			updateFunc: func(p *ProjectResource) {
				p.Allocation = 150
			},
			wantError: ErrProjectResourceInvalidAllocation,
		},
		{
			name: "Update fails - negative allocation",
			updateFunc: func(p *ProjectResource) {
				p.Allocation = -10
			},
			wantError: ErrProjectResourceInvalidAllocation,
		},
		{
			name: "Update fails - invalid status",
			updateFunc: func(p *ProjectResource) {
				p.Status = 999
			},
			wantError: ErrProjectResourceInvalidStatus,
		},
		{
			name: "Update fails - invalid dates",
			updateFunc: func(p *ProjectResource) {
				p.StartDate = &now
				p.EndDate = &pastDate
			},
			wantError: ErrProjectResourceInvalidDates,
		},
		{
			name: "Update status to inactive",
			updateFunc: func(p *ProjectResource) {
				p.Status = ProjectResourceStatusInactive
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reload project resource from database
			var testPR ProjectResource
			db.First(&testPR, pr.ID)

			// Apply update
			tt.updateFunc(&testPR)

			// Attempt to save
			result := db.Save(&testPR)

			if tt.wantError != nil {
				assert.Error(t, result.Error)
				assert.Contains(t, result.Error.Error(), tt.wantError.Error())
			} else {
				assert.NoError(t, result.Error)
			}
		})
	}
}

func TestProjectResourceCRUDOperations(t *testing.T) {
	db := setupProjectResourceTestDB(t)
	client := createTestClientForPR(t, db)
	project := createTestProjectForPR(t, db, client.ID)
	hr := createTestHumanResource(t, db)

	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)

	// Create
	pr := ProjectResource{
		ProjectID:       project.ID,
		HumanResourceID: hr.ID,
		Role:            "Developer",
		Allocation:      100,
		StartDate:       &now,
		EndDate:         &futureDate,
		Notes:           "Primary developer for this project",
		Status:          ProjectResourceStatusActive,
	}
	result := db.Create(&pr)
	assert.NoError(t, result.Error)
	assert.NotZero(t, pr.ID)

	// Read
	var retrieved ProjectResource
	result = db.First(&retrieved, pr.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, pr.ProjectID, retrieved.ProjectID)
	assert.Equal(t, pr.HumanResourceID, retrieved.HumanResourceID)
	assert.Equal(t, pr.Role, retrieved.Role)
	assert.Equal(t, pr.Allocation, retrieved.Allocation)
	assert.Equal(t, pr.Notes, retrieved.Notes)
	assert.Equal(t, pr.Status, retrieved.Status)

	// Update
	retrieved.Role = "Tech Lead"
	retrieved.Allocation = 80
	result = db.Save(&retrieved)
	assert.NoError(t, result.Error)

	// Verify update
	var updated ProjectResource
	result = db.First(&updated, pr.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Tech Lead", updated.Role)
	assert.Equal(t, float64(80), updated.Allocation)

	// Delete
	result = db.Delete(&updated)
	assert.NoError(t, result.Error)

	// Verify deletion
	var deleted ProjectResource
	result = db.First(&deleted, pr.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestProjectResourceQueryOperations(t *testing.T) {
	db := setupProjectResourceTestDB(t)
	client := createTestClientForPR(t, db)
	project1 := createTestProjectForPR(t, db, client.ID)
	project2 := &Project{Name: "Project 2", ClientID: client.ID, Status: ProjectStatusActive}
	db.Create(project2)

	hr1 := createTestHumanResource(t, db)
	hr2 := &HumanResource{Name: "Jane Doe", Title: "QA Engineer", Level: "Mid", Status: HumanResourceStatusActive}
	db.Create(hr2)

	// Create multiple project resources
	prs := []ProjectResource{
		{ProjectID: project1.ID, HumanResourceID: hr1.ID, Role: "Developer", Allocation: 100, Status: ProjectResourceStatusActive},
		{ProjectID: project1.ID, HumanResourceID: hr2.ID, Role: "QA", Allocation: 50, Status: ProjectResourceStatusActive},
		{ProjectID: project2.ID, HumanResourceID: hr1.ID, Role: "Tech Lead", Allocation: 50, Status: ProjectResourceStatusInactive},
	}

	for i := range prs {
		db.Create(&prs[i])
	}

	// Query all project resources
	var allPRs []ProjectResource
	result := db.Find(&allPRs)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(allPRs), 3)

	// Query active project resources only
	var activePRs []ProjectResource
	result = db.Where("status = ?", ProjectResourceStatusActive).Find(&activePRs)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(activePRs), 2)

	// Query by project ID
	var projectPRs []ProjectResource
	result = db.Where("project_id = ?", project1.ID).Find(&projectPRs)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(projectPRs), 2)

	// Query by human resource ID
	var hrPRs []ProjectResource
	result = db.Where("human_resource_id = ?", hr1.ID).Find(&hrPRs)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(hrPRs), 2)

	// Query by allocation
	var fullAllocationPRs []ProjectResource
	result = db.Where("allocation = ?", 100).Find(&fullAllocationPRs)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(fullAllocationPRs), 1)
}

func TestProjectResourceRelationships(t *testing.T) {
	db := setupProjectResourceTestDB(t)
	client := createTestClientForPR(t, db)
	project := createTestProjectForPR(t, db, client.ID)
	hr := createTestHumanResource(t, db)

	// Create project resource
	pr := ProjectResource{
		ProjectID:       project.ID,
		HumanResourceID: hr.ID,
		Role:            "Developer",
		Allocation:      100,
		Status:          ProjectResourceStatusActive,
	}
	result := db.Create(&pr)
	assert.NoError(t, result.Error)

	// Load project resource with project relationship
	var loadedPR ProjectResource
	result = db.Preload("Project").First(&loadedPR, pr.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, loadedPR.Project)
	assert.Equal(t, project.ID, loadedPR.Project.ID)
	assert.Equal(t, project.Name, loadedPR.Project.Name)

	// Load project resource with human resource relationship
	var loadedPR2 ProjectResource
	result = db.Preload("HumanResource").First(&loadedPR2, pr.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, loadedPR2.HumanResource)
	assert.Equal(t, hr.ID, loadedPR2.HumanResource.ID)
	assert.Equal(t, hr.Name, loadedPR2.HumanResource.Name)

	// Load project resource with both relationships
	var loadedPR3 ProjectResource
	result = db.Preload("Project").Preload("HumanResource").First(&loadedPR3, pr.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, loadedPR3.Project)
	assert.NotNil(t, loadedPR3.HumanResource)
}

func TestProjectResourceAllowedSortFields(t *testing.T) {
	expectedFields := map[string]string{
		"id":                "id",
		"project_id":        "project_id",
		"human_resource_id": "human_resource_id",
		"role":              "role",
		"allocation":        "allocation",
		"start_date":        "start_date",
		"end_date":          "end_date",
		"status":            "status",
		"created_at":        "created_at",
		"updated_at":        "updated_at",
	}

	assert.Equal(t, expectedFields, ProjectResourceAllowedSortField)
}

func TestProjectWithProjectResources(t *testing.T) {
	db := setupProjectResourceTestDB(t)
	client := createTestClientForPR(t, db)
	project := createTestProjectForPR(t, db, client.ID)
	hr1 := createTestHumanResource(t, db)
	hr2 := &HumanResource{Name: "Jane Doe", Title: "QA Engineer", Level: "Mid", Status: HumanResourceStatusActive}
	db.Create(hr2)

	// Create multiple project resources for the project
	pr1 := ProjectResource{ProjectID: project.ID, HumanResourceID: hr1.ID, Role: "Developer", Allocation: 100, Status: ProjectResourceStatusActive}
	pr2 := ProjectResource{ProjectID: project.ID, HumanResourceID: hr2.ID, Role: "QA", Allocation: 50, Status: ProjectResourceStatusActive}
	db.Create(&pr1)
	db.Create(&pr2)

	// Load project with project resources
	var loadedProject Project
	result := db.Preload("ProjectResources").First(&loadedProject, project.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, loadedProject.ProjectResources)
	assert.Len(t, loadedProject.ProjectResources, 2)

	// Verify the project resources
	roles := make([]string, 0, 2)
	for _, pr := range loadedProject.ProjectResources {
		roles = append(roles, pr.Role)
	}
	assert.Contains(t, roles, "Developer")
	assert.Contains(t, roles, "QA")
}

func BenchmarkProjectResourceValidate(b *testing.B) {
	pr := ProjectResource{
		ProjectID:       1,
		HumanResourceID: 1,
		Allocation:      100,
		Status:          ProjectResourceStatusActive,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pr.Validate()
	}
}

func BenchmarkProjectResourceCreate(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Client{}, &HumanResource{}, &Project{}, &ProjectResource{})

	client := &Client{Name: "Test", Email: "test@test.com", Status: ClientStatusActive}
	db.Create(client)
	project := &Project{Name: "Test", ClientID: client.ID, Status: ProjectStatusActive}
	db.Create(project)
	hr := &HumanResource{Name: "Test", Title: "Dev", Level: "Sr", Status: HumanResourceStatusActive}
	db.Create(hr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pr := ProjectResource{
			ProjectID:       project.ID,
			HumanResourceID: hr.ID,
			Allocation:      100,
			Status:          ProjectResourceStatusActive,
		}
		db.Create(&pr)
	}
}
