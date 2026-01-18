package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProjectTableName(t *testing.T) {
	project := Project{}
	assert.Equal(t, "projects", project.TableName())
}

func TestProjectIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status uint
		want   bool
	}{
		{"Active project", ProjectStatusActive, true},
		{"Inactive project", ProjectStatusInactive, false},
		{"Invalid status treated as inactive", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := Project{Status: tt.status}
			assert.Equal(t, tt.want, project.IsActive())
		})
	}
}

func TestProjectStatusChange(t *testing.T) {
	project := Project{Status: ProjectStatusInactive}

	// Test changing to active
	project.Status = ProjectStatusActive
	assert.Equal(t, uint(ProjectStatusActive), project.Status)
	assert.True(t, project.IsActive())

	// Test changing to inactive
	project.Status = ProjectStatusInactive
	assert.Equal(t, uint(ProjectStatusInactive), project.Status)
	assert.False(t, project.IsActive())
}

func TestProjectValidateStatus(t *testing.T) {
	tests := []struct {
		name      string
		status    uint
		wantError error
	}{
		{"Valid: Active status", ProjectStatusActive, nil},
		{"Valid: Inactive status", ProjectStatusInactive, nil},
		{"Invalid: Status value 3", 3, ErrProjectInvalidStatus},
		{"Invalid: High value", 999, ErrProjectInvalidStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := Project{Status: tt.status}
			err := project.validateStatus()
			if tt.wantError != nil {
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectValidate(t *testing.T) {
	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)
	pastDate := now.AddDate(0, -1, 0)

	tests := []struct {
		name      string
		project   Project
		wantError error
	}{
		{
			name: "Valid project",
			project: Project{
				Name:     "Project Alpha",
				ClientID: 1,
				Status:   ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid with whitespace trimmed",
			project: Project{
				Name:        "  Project Alpha  ",
				Description: "  A test project  ",
				ClientID:    1,
				Status:      ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Empty name",
			project: Project{
				Name:     "",
				ClientID: 1,
				Status:   ProjectStatusActive,
			},
			wantError: ErrProjectNameRequired,
		},
		{
			name: "Whitespace-only name",
			project: Project{
				Name:     "   ",
				ClientID: 1,
				Status:   ProjectStatusActive,
			},
			wantError: ErrProjectNameRequired,
		},
		{
			name: "Missing client ID",
			project: Project{
				Name:     "Project Alpha",
				ClientID: 0,
				Status:   ProjectStatusActive,
			},
			wantError: ErrProjectInvalidClientID,
		},
		{
			name: "Invalid status",
			project: Project{
				Name:     "Project Alpha",
				ClientID: 1,
				Status:   999,
			},
			wantError: ErrProjectInvalidStatus,
		},
		{
			name: "Valid with dates",
			project: Project{
				Name:      "Project Alpha",
				ClientID:  1,
				StartDate: &now,
				EndDate:   &futureDate,
				Status:    ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Invalid dates - end before start",
			project: Project{
				Name:      "Project Alpha",
				ClientID:  1,
				StartDate: &now,
				EndDate:   &pastDate,
				Status:    ProjectStatusActive,
			},
			wantError: ErrProjectInvalidDates,
		},
		{
			name: "Valid with only start date",
			project: Project{
				Name:      "Project Alpha",
				ClientID:  1,
				StartDate: &now,
				Status:    ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid with only end date",
			project: Project{
				Name:     "Project Alpha",
				ClientID: 1,
				EndDate:  &futureDate,
				Status:   ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid with all optional fields",
			project: Project{
				Name:        "Project Alpha",
				Description: "A comprehensive project description",
				ClientID:    1,
				StartDate:   &now,
				EndDate:     &futureDate,
				Status:      ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid with configuration fields",
			project: Project{
				Name:               "Project Alpha",
				ClientID:           1,
				Status:             ProjectStatusActive,
				HoursPerDay:        6,
				DaysPerWeek:        4,
				WorkingDaysPerWeek: WeekdayArray{time.Monday, time.Tuesday, time.Wednesday, time.Thursday},
				Timezone:           "Asia/Ho_Chi_Minh",
				Currency:           "VND",
			},
			wantError: nil,
		},
		{
			name: "Invalid hours per day - too low",
			project: Project{
				Name:        "Project Alpha",
				ClientID:    1,
				Status:      ProjectStatusActive,
				HoursPerDay: 0, // 0 is allowed (uses default)
			},
			wantError: nil,
		},
		{
			name: "Invalid hours per day - negative",
			project: Project{
				Name:        "Project Alpha",
				ClientID:    1,
				Status:      ProjectStatusActive,
				HoursPerDay: -1,
			},
			wantError: ErrProjectInvalidHoursPerDay,
		},
		{
			name: "Invalid hours per day - too high",
			project: Project{
				Name:        "Project Alpha",
				ClientID:    1,
				Status:      ProjectStatusActive,
				HoursPerDay: 25,
			},
			wantError: ErrProjectInvalidHoursPerDay,
		},
		{
			name: "Invalid days per week - too high",
			project: Project{
				Name:        "Project Alpha",
				ClientID:    1,
				Status:      ProjectStatusActive,
				DaysPerWeek: 8,
			},
			wantError: ErrProjectInvalidDaysPerWeek,
		},
		{
			name: "Invalid days per week - negative",
			project: Project{
				Name:        "Project Alpha",
				ClientID:    1,
				Status:      ProjectStatusActive,
				DaysPerWeek: -1,
			},
			wantError: ErrProjectInvalidDaysPerWeek,
		},
		{
			name: "Valid with timezone and currency whitespace trimmed",
			project: Project{
				Name:     "Project Alpha",
				ClientID: 1,
				Status:   ProjectStatusActive,
				Timezone: "  UTC  ",
				Currency: "  USD  ",
			},
			wantError: nil,
		},
		{
			name: "Valid working days - all days",
			project: Project{
				Name:               "Project Alpha",
				ClientID:           1,
				Status:             ProjectStatusActive,
				WorkingDaysPerWeek: WeekdayArray{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday},
			},
			wantError: nil,
		},
		{
			name: "Invalid working days - duplicate",
			project: Project{
				Name:               "Project Alpha",
				ClientID:           1,
				Status:             ProjectStatusActive,
				WorkingDaysPerWeek: WeekdayArray{time.Monday, time.Monday, time.Tuesday},
			},
			wantError: ErrProjectDuplicateWorkingDays,
		},
		{
			name: "Invalid working days - value out of range",
			project: Project{
				Name:               "Project Alpha",
				ClientID:           1,
				Status:             ProjectStatusActive,
				WorkingDaysPerWeek: WeekdayArray{time.Weekday(7)},
			},
			wantError: ErrProjectInvalidWorkingDays,
		},
		{
			name: "Invalid working days - negative value",
			project: Project{
				Name:               "Project Alpha",
				ClientID:           1,
				Status:             ProjectStatusActive,
				WorkingDaysPerWeek: WeekdayArray{time.Weekday(-1)},
			},
			wantError: ErrProjectInvalidWorkingDays,
		},
		{
			name: "Invalid working days - exceeds 7 days",
			project: Project{
				Name:               "Project Alpha",
				ClientID:           1,
				Status:             ProjectStatusActive,
				WorkingDaysPerWeek: WeekdayArray{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
			},
			wantError: ErrProjectWorkingDaysExceedsWeek,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Validate()
			if tt.wantError != nil {
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
				// Verify whitespace was trimmed
				assert.NotContains(t, tt.project.Name, "  ")
			}
		})
	}
}

func TestProjectValidateTrimsWhitespace(t *testing.T) {
	project := Project{
		Name:        "  Project Alpha  ",
		Description: "  A test project  ",
		ClientID:    1,
		Status:      ProjectStatusActive,
		Timezone:    "  UTC  ",
		Currency:    "  USD  ",
	}

	err := project.Validate()
	assert.NoError(t, err)

	// Verify all fields were trimmed
	assert.Equal(t, "Project Alpha", project.Name)
	assert.Equal(t, "A test project", project.Description)
	assert.Equal(t, "UTC", project.Timezone)
	assert.Equal(t, "USD", project.Currency)
}

func TestProjectGetHoursPerDay(t *testing.T) {
	tests := []struct {
		name        string
		hoursPerDay int
		want        int
	}{
		{"Returns default when zero", 0, DefaultProjectHoursPerDay},
		{"Returns custom value when set", 6, 6},
		{"Returns custom value 10", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := Project{HoursPerDay: tt.hoursPerDay}
			assert.Equal(t, tt.want, project.GetHoursPerDay())
		})
	}
}

func TestProjectGetDaysPerWeek(t *testing.T) {
	tests := []struct {
		name        string
		daysPerWeek int
		want        int
	}{
		{"Returns default when zero", 0, DefaultProjectDaysPerWeek},
		{"Returns custom value when set", 4, 4},
		{"Returns custom value 7", 7, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := Project{DaysPerWeek: tt.daysPerWeek}
			assert.Equal(t, tt.want, project.GetDaysPerWeek())
		})
	}
}

func TestProjectGetWorkingDaysPerWeek(t *testing.T) {
	customDays := WeekdayArray{time.Monday, time.Tuesday, time.Wednesday}
	defaultDays := DefaultWorkingDays()

	tests := []struct {
		name               string
		workingDaysPerWeek WeekdayArray
		want               WeekdayArray
	}{
		{"Returns default when nil", nil, defaultDays},
		{"Returns default when empty", WeekdayArray{}, defaultDays},
		{"Returns custom value when set", customDays, customDays},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := Project{WorkingDaysPerWeek: tt.workingDaysPerWeek}
			assert.Equal(t, tt.want, project.GetWorkingDaysPerWeek())
		})
	}
}

func setupProjectTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the Client and Project tables
	err = db.AutoMigrate(&Client{}, &Project{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func createTestClient(t *testing.T, db *gorm.DB) *Client {
	client := &Client{
		Name:   "Test Client",
		Email:  "test@client.com",
		Status: ClientStatusActive,
	}
	result := db.Create(client)
	assert.NoError(t, result.Error)
	return client
}

func TestProjectBeforeCreate(t *testing.T) {
	db := setupProjectTestDB(t)
	client := createTestClient(t, db)

	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)
	pastDate := now.AddDate(0, -1, 0)

	tests := []struct {
		name      string
		project   Project
		wantError error
	}{
		{
			name: "Valid project creation",
			project: Project{
				Name:     "Project Alpha",
				ClientID: client.ID,
				Status:   ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Default status set to active when invalid",
			project: Project{
				Name:     "Project Beta",
				ClientID: client.ID,
				Status:   999, // Invalid status
			},
			wantError: nil, // Should succeed after setting default
		},
		{
			name: "Default status set to active when zero",
			project: Project{
				Name:     "Project Gamma",
				ClientID: client.ID,
				// Status not set (will be 0)
			},
			wantError: nil, // Should succeed after setting default
		},
		{
			name: "Validation fails - empty name",
			project: Project{
				Name:     "",
				ClientID: client.ID,
				Status:   ProjectStatusActive,
			},
			wantError: ErrProjectNameRequired,
		},
		{
			name: "Validation fails - missing client ID",
			project: Project{
				Name:     "Project Delta",
				ClientID: 0,
				Status:   ProjectStatusActive,
			},
			wantError: ErrProjectInvalidClientID,
		},
		{
			name: "Validation fails - invalid dates",
			project: Project{
				Name:      "Project Epsilon",
				ClientID:  client.ID,
				StartDate: &now,
				EndDate:   &pastDate,
				Status:    ProjectStatusActive,
			},
			wantError: ErrProjectInvalidDates,
		},
		{
			name: "Valid with dates",
			project: Project{
				Name:      "Project Zeta",
				ClientID:  client.ID,
				StartDate: &now,
				EndDate:   &futureDate,
				Status:    ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Whitespace trimmed during creation",
			project: Project{
				Name:        "  Project Eta  ",
				Description: "  Description  ",
				ClientID:    client.ID,
				Status:      ProjectStatusActive,
			},
			wantError: nil,
		},
		{
			name: "Valid with custom configuration",
			project: Project{
				Name:               "Project Theta",
				ClientID:           client.ID,
				Status:             ProjectStatusActive,
				HoursPerDay:        6,
				DaysPerWeek:        4,
				WorkingDaysPerWeek: WeekdayArray{time.Monday, time.Tuesday, time.Wednesday, time.Thursday},
				Timezone:           "Asia/Ho_Chi_Minh",
				Currency:           "VND",
			},
			wantError: nil,
		},
		{
			name: "Validation fails - invalid hours per day",
			project: Project{
				Name:        "Project Iota",
				ClientID:    client.ID,
				Status:      ProjectStatusActive,
				HoursPerDay: 25,
			},
			wantError: ErrProjectInvalidHoursPerDay,
		},
		{
			name: "Validation fails - invalid days per week",
			project: Project{
				Name:        "Project Kappa",
				ClientID:    client.ID,
				Status:      ProjectStatusActive,
				DaysPerWeek: 8,
			},
			wantError: ErrProjectInvalidDaysPerWeek,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.Create(&tt.project)

			if tt.wantError != nil {
				assert.Error(t, result.Error)
				assert.Contains(t, result.Error.Error(), tt.wantError.Error())
			} else {
				assert.NoError(t, result.Error)
				assert.NotZero(t, tt.project.ID)
				assert.NotZero(t, tt.project.CreatedAt)
				assert.NotZero(t, tt.project.UpdatedAt)

				// Verify status was set to active if it was invalid
				if tt.name == "Default status set to active when invalid" || tt.name == "Default status set to active when zero" {
					assert.Equal(t, uint(ProjectStatusActive), tt.project.Status)
				}

				// Verify whitespace was trimmed
				assert.NotContains(t, tt.project.Name, "  ")
				assert.NotContains(t, tt.project.Description, "  ")
			}
		})
	}
}

func TestProjectBeforeCreateDefaultConfiguration(t *testing.T) {
	db := setupProjectTestDB(t)
	client := createTestClient(t, db)

	// Create a project without specifying configuration values
	project := Project{
		Name:     "Project with defaults",
		ClientID: client.ID,
		Status:   ProjectStatusActive,
	}
	result := db.Create(&project)
	assert.NoError(t, result.Error)

	// Verify default working days were set by BeforeCreate hook
	assert.Equal(t, DefaultWorkingDays(), project.WorkingDaysPerWeek)

	// Verify from database (HoursPerDay and DaysPerWeek use GORM defaults)
	var retrieved Project
	db.First(&retrieved, project.ID)
	assert.Equal(t, DefaultProjectHoursPerDay, retrieved.HoursPerDay)
	assert.Equal(t, DefaultProjectDaysPerWeek, retrieved.DaysPerWeek)
	assert.Equal(t, DefaultWorkingDays(), retrieved.WorkingDaysPerWeek)
}

func TestProjectBeforeUpdate(t *testing.T) {
	db := setupProjectTestDB(t)
	client := createTestClient(t, db)

	now := time.Now()
	pastDate := now.AddDate(0, -1, 0)

	// Create a valid project first
	project := Project{
		Name:     "Project Alpha",
		ClientID: client.ID,
		Status:   ProjectStatusActive,
	}
	result := db.Create(&project)
	assert.NoError(t, result.Error)

	tests := []struct {
		name       string
		updateFunc func(*Project)
		wantError  error
	}{
		{
			name: "Valid update",
			updateFunc: func(p *Project) {
				p.Name = "Project Alpha Updated"
				p.Description = "Updated description"
			},
			wantError: nil,
		},
		{
			name: "Update with whitespace trimming",
			updateFunc: func(p *Project) {
				p.Name = "  Updated Name  "
				p.Description = "  Updated description  "
			},
			wantError: nil,
		},
		{
			name: "Update fails - empty name",
			updateFunc: func(p *Project) {
				p.Name = ""
			},
			wantError: ErrProjectNameRequired,
		},
		{
			name: "Update fails - invalid status",
			updateFunc: func(p *Project) {
				p.Status = 999
			},
			wantError: ErrProjectInvalidStatus,
		},
		{
			name: "Update fails - invalid dates",
			updateFunc: func(p *Project) {
				p.StartDate = &now
				p.EndDate = &pastDate
			},
			wantError: ErrProjectInvalidDates,
		},
		{
			name: "Update status to inactive",
			updateFunc: func(p *Project) {
				p.Status = ProjectStatusInactive
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reload project from database
			var testProject Project
			db.First(&testProject, project.ID)

			// Apply update
			tt.updateFunc(&testProject)

			// Attempt to save
			result := db.Save(&testProject)

			if tt.wantError != nil {
				assert.Error(t, result.Error)
				assert.Contains(t, result.Error.Error(), tt.wantError.Error())
			} else {
				assert.NoError(t, result.Error)

				// Verify whitespace was trimmed
				assert.NotContains(t, testProject.Name, "  ")
				assert.NotContains(t, testProject.Description, "  ")
			}
		})
	}
}

func TestProjectCRUDOperations(t *testing.T) {
	db := setupProjectTestDB(t)
	client := createTestClient(t, db)

	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)

	// Create
	project := Project{
		Name:        "Project Alpha",
		Description: "A comprehensive project",
		ClientID:    client.ID,
		StartDate:   &now,
		EndDate:     &futureDate,
		Status:      ProjectStatusActive,
	}
	result := db.Create(&project)
	assert.NoError(t, result.Error)
	assert.NotZero(t, project.ID)

	// Read
	var retrieved Project
	result = db.First(&retrieved, project.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, project.Name, retrieved.Name)
	assert.Equal(t, project.Description, retrieved.Description)
	assert.Equal(t, project.ClientID, retrieved.ClientID)
	assert.Equal(t, project.Status, retrieved.Status)

	// Update
	retrieved.Name = "Project Alpha Updated"
	retrieved.Description = "Updated description"
	result = db.Save(&retrieved)
	assert.NoError(t, result.Error)

	// Verify update
	var updated Project
	result = db.First(&updated, project.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Project Alpha Updated", updated.Name)
	assert.Equal(t, "Updated description", updated.Description)

	// Delete
	result = db.Delete(&updated)
	assert.NoError(t, result.Error)

	// Verify deletion
	var deleted Project
	result = db.First(&deleted, project.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestProjectQueryOperations(t *testing.T) {
	db := setupProjectTestDB(t)
	client := createTestClient(t, db)

	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)

	// Create multiple projects
	projects := []Project{
		{Name: "Project Alpha", ClientID: client.ID, StartDate: &now, Status: ProjectStatusActive},
		{Name: "Project Beta", ClientID: client.ID, EndDate: &futureDate, Status: ProjectStatusActive},
		{Name: "Project Gamma", ClientID: client.ID, Status: ProjectStatusInactive},
	}

	for i := range projects {
		db.Create(&projects[i])
	}

	// Query all projects
	var allProjects []Project
	result := db.Find(&allProjects)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(allProjects), 3)

	// Query active projects only
	var activeProjects []Project
	result = db.Where("status = ?", ProjectStatusActive).Find(&activeProjects)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(activeProjects), 2)

	// Query by client ID
	var clientProjects []Project
	result = db.Where("client_id = ?", client.ID).Find(&clientProjects)
	assert.NoError(t, result.Error)
	assert.GreaterOrEqual(t, len(clientProjects), 3)

	// Query by name
	var project Project
	result = db.Where("name = ?", "Project Alpha").First(&project)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Project Alpha", project.Name)
}

func TestProjectClientRelationship(t *testing.T) {
	db := setupProjectTestDB(t)
	client := createTestClient(t, db)

	// Create project
	project := Project{
		Name:     "Project Alpha",
		ClientID: client.ID,
		Status:   ProjectStatusActive,
	}
	result := db.Create(&project)
	assert.NoError(t, result.Error)

	// Load project with client relationship
	var loadedProject Project
	result = db.Preload("Client").First(&loadedProject, project.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, loadedProject.Client)
	assert.Equal(t, client.ID, loadedProject.Client.ID)
	assert.Equal(t, client.Name, loadedProject.Client.Name)
}

func TestProjectAllowedSortFields(t *testing.T) {
	expectedFields := map[string]string{
		"id":          "id",
		"name":        "name",
		"description": "description",
		"client_id":   "client_id",
		"start_date":  "start_date",
		"end_date":    "end_date",
		"status":      "status",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
	}

	assert.Equal(t, expectedFields, ProjectAllowedSortField)
}

func BenchmarkProjectValidate(b *testing.B) {
	project := Project{
		Name:     "Project Alpha",
		ClientID: 1,
		Status:   ProjectStatusActive,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = project.Validate()
	}
}

func BenchmarkProjectCreate(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Client{}, &Project{}) //nolint:errcheck

	client := &Client{Name: "Test", Email: "test@test.com", Status: ClientStatusActive}
	db.Create(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		project := Project{
			Name:     "Project Alpha",
			ClientID: client.ID,
			Status:   ProjectStatusActive,
		}
		db.Create(&project)
	}
}
