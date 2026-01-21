package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProjectRoleTableName(t *testing.T) {
	pr := ProjectRole{}
	assert.Equal(t, "project_roles", pr.TableName())
}

func TestRoleLevelName(t *testing.T) {
	tests := []struct {
		name  string
		level uint
		want  string
	}{
		{"Junior level", RoleLevelJunior, "Junior"},
		{"Mid level", RoleLevelMid, "Mid"},
		{"Senior level", RoleLevelSenior, "Senior"},
		{"Lead level", RoleLevelLead, "Lead"},
		{"Unknown level", RoleLevelUnknown, "Unknown"},
		{"Invalid level", 999, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoleLevelName(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectRoleGetLevelName(t *testing.T) {
	tests := []struct {
		name  string
		level uint
		want  string
	}{
		{"Junior level", RoleLevelJunior, "Junior"},
		{"Mid level", RoleLevelMid, "Mid"},
		{"Senior level", RoleLevelSenior, "Senior"},
		{"Lead level", RoleLevelLead, "Lead"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := ProjectRole{Level: tt.level}
			assert.Equal(t, tt.want, pr.GetLevelName())
		})
	}
}

func TestProjectRoleValidateLevel(t *testing.T) {
	tests := []struct {
		name      string
		level     uint
		wantError error
	}{
		{"Valid: Junior level", RoleLevelJunior, nil},
		{"Valid: Mid level", RoleLevelMid, nil},
		{"Valid: Senior level", RoleLevelSenior, nil},
		{"Valid: Lead level", RoleLevelLead, nil},
		{"Valid: Manager level", RoleLevelManager, nil},
		{"Valid: Director level", RoleLevelDirector, nil},
		{"Valid: VP level", RoleLevelVP, nil},
		{"Valid: C-Level level", RoleLevelCLevel, nil},
		{"Invalid: Unknown level", RoleLevelUnknown, ErrProjectRoleInvalidLevel},
		{"Invalid: High value", 999, ErrProjectRoleInvalidLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := ProjectRole{Level: tt.level}
			err := pr.validateLevel()
			if tt.wantError != nil {
				assert.Equal(t, tt.wantError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectRoleValidate(t *testing.T) {
	tests := []struct {
		name      string
		pr        ProjectRole
		wantError error
	}{
		{
			name: "Valid project role",
			pr: ProjectRole{
				ProjectID: 1,
				Name:      "Backend Developer",
				Level:     RoleLevelMid,
				Headcount: 2,
			},
			wantError: nil,
		},
		{
			name: "Valid with minimum headcount",
			pr: ProjectRole{
				ProjectID: 1,
				Name:      "QA",
				Level:     RoleLevelJunior,
				Headcount: 0,
			},
			wantError: nil,
		},
		{
			name: "Valid with trimmed name",
			pr: ProjectRole{
				ProjectID: 1,
				Name:      "  Frontend Developer  ",
				Level:     RoleLevelSenior,
				Headcount: 1,
			},
			wantError: nil,
		},
		{
			name: "Missing name",
			pr: ProjectRole{
				ProjectID: 1,
				Name:      "",
				Level:     RoleLevelMid,
				Headcount: 1,
			},
			wantError: ErrProjectRoleNameRequired,
		},
		{
			name: "Whitespace-only name",
			pr: ProjectRole{
				ProjectID: 1,
				Name:      "   ",
				Level:     RoleLevelMid,
				Headcount: 1,
			},
			wantError: ErrProjectRoleNameRequired,
		},
		{
			name: "Missing project ID",
			pr: ProjectRole{
				ProjectID: 0,
				Name:      "Developer",
				Level:     RoleLevelMid,
				Headcount: 1,
			},
			wantError: ErrProjectRoleInvalidProjectID,
		},
		{
			name: "Invalid level - unknown",
			pr: ProjectRole{
				ProjectID: 1,
				Name:      "Developer",
				Level:     RoleLevelUnknown,
				Headcount: 1,
			},
			wantError: ErrProjectRoleInvalidLevel,
		},
		{
			name: "Invalid level - out of range",
			pr: ProjectRole{
				ProjectID: 1,
				Name:      "Developer",
				Level:     999,
				Headcount: 1,
			},
			wantError: ErrProjectRoleInvalidLevel,
		},
		{
			name: "Invalid headcount - negative",
			pr: ProjectRole{
				ProjectID: 1,
				Name:      "Developer",
				Level:     RoleLevelMid,
				Headcount: -1,
			},
			wantError: ErrProjectRoleInvalidHeadcount,
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

func setupProjectRoleTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all required tables
	err = db.AutoMigrate(&Client{}, &Project{}, &ProjectRole{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func createTestClientForProjectRole(t *testing.T, db *gorm.DB) *Client {
	client := &Client{
		Name:   "Test Client",
		Email:  "test@client.com",
		Status: ClientStatusActive,
	}
	result := db.Create(client)
	assert.NoError(t, result.Error)
	return client
}

func createTestProjectForProjectRole(t *testing.T, db *gorm.DB, clientID uint) *Project {
	project := &Project{
		Name:     "Test Project",
		ClientID: clientID,
		Status:   ProjectStatusActive,
	}
	result := db.Create(project)
	assert.NoError(t, result.Error)
	return project
}

func TestProjectRoleBeforeCreate(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	client := createTestClientForProjectRole(t, db)
	project := createTestProjectForProjectRole(t, db, client.ID)

	t.Run("Valid project role creation", func(t *testing.T) {
		pr := ProjectRole{
			ProjectID: project.ID,
			Name:      "Backend Developer",
			Level:     RoleLevelMid,
			Headcount: 2,
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.NotZero(t, pr.ID)
		assert.NotZero(t, pr.CreatedAt)
		assert.NotZero(t, pr.UpdatedAt)
	})

	t.Run("Default level set to mid when invalid", func(t *testing.T) {
		pr := ProjectRole{
			ProjectID: project.ID,
			Name:      "Designer",
			Level:     999, // Invalid level
			Headcount: 1,
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.Equal(t, uint(RoleLevelMid), pr.Level)
	})

	t.Run("Default level set to mid when zero", func(t *testing.T) {
		pr := ProjectRole{
			ProjectID: project.ID,
			Name:      "DevOps",
			// Level not set (will be 0)
			Headcount: 1,
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.Equal(t, uint(RoleLevelMid), pr.Level)
	})

	t.Run("Default headcount set to 1 when zero", func(t *testing.T) {
		pr := ProjectRole{
			ProjectID: project.ID,
			Name:      "QA Engineer",
			Level:     RoleLevelJunior,
			// Headcount not set (will be 0)
		}
		result := db.Create(&pr)
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, pr.Headcount)
	})

	t.Run("Validation fails - missing project ID", func(t *testing.T) {
		pr := ProjectRole{
			ProjectID: 0,
			Name:      "Developer",
			Level:     RoleLevelMid,
			Headcount: 1,
		}
		result := db.Create(&pr)
		assert.Error(t, result.Error)
		assert.Contains(t, result.Error.Error(), ErrProjectRoleInvalidProjectID.Error())
	})

	t.Run("Validation fails - missing name", func(t *testing.T) {
		pr := ProjectRole{
			ProjectID: project.ID,
			Name:      "",
			Level:     RoleLevelMid,
			Headcount: 1,
		}
		result := db.Create(&pr)
		assert.Error(t, result.Error)
		assert.Contains(t, result.Error.Error(), ErrProjectRoleNameRequired.Error())
	})

	t.Run("Validation fails - negative headcount", func(t *testing.T) {
		pr := ProjectRole{
			ProjectID: project.ID,
			Name:      "Tester",
			Level:     RoleLevelMid,
			Headcount: -5,
		}
		result := db.Create(&pr)
		assert.Error(t, result.Error)
		assert.Contains(t, result.Error.Error(), ErrProjectRoleInvalidHeadcount.Error())
	})
}

func TestProjectRoleBeforeUpdate(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	client := createTestClientForProjectRole(t, db)
	project := createTestProjectForProjectRole(t, db, client.ID)

	// Create a valid project role first
	pr := ProjectRole{
		ProjectID: project.ID,
		Name:      "Developer",
		Level:     RoleLevelMid,
		Headcount: 2,
	}
	result := db.Create(&pr)
	assert.NoError(t, result.Error)

	tests := []struct {
		name       string
		updateFunc func(*ProjectRole)
		wantError  error
	}{
		{
			name: "Valid update - change name",
			updateFunc: func(p *ProjectRole) {
				p.Name = "Senior Developer"
			},
			wantError: nil,
		},
		{
			name: "Valid update - change level",
			updateFunc: func(p *ProjectRole) {
				p.Level = RoleLevelSenior
			},
			wantError: nil,
		},
		{
			name: "Valid update - change headcount",
			updateFunc: func(p *ProjectRole) {
				p.Headcount = 5
			},
			wantError: nil,
		},
		{
			name: "Update fails - empty name",
			updateFunc: func(p *ProjectRole) {
				p.Name = ""
			},
			wantError: ErrProjectRoleNameRequired,
		},
		{
			name: "Update fails - invalid level",
			updateFunc: func(p *ProjectRole) {
				p.Level = 999
			},
			wantError: ErrProjectRoleInvalidLevel,
		},
		{
			name: "Update fails - negative headcount",
			updateFunc: func(p *ProjectRole) {
				p.Headcount = -1
			},
			wantError: ErrProjectRoleInvalidHeadcount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reload project role from database
			var testPR ProjectRole
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

func TestProjectRoleCRUDOperations(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	client := createTestClientForProjectRole(t, db)
	project := createTestProjectForProjectRole(t, db, client.ID)

	// Create
	pr := ProjectRole{
		ProjectID: project.ID,
		Name:      "Backend Developer",
		Level:     RoleLevelSenior,
		Headcount: 3,
	}
	result := db.Create(&pr)
	assert.NoError(t, result.Error)
	assert.NotZero(t, pr.ID)

	// Read
	var retrieved ProjectRole
	result = db.First(&retrieved, pr.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, pr.ProjectID, retrieved.ProjectID)
	assert.Equal(t, pr.Name, retrieved.Name)
	assert.Equal(t, pr.Level, retrieved.Level)
	assert.Equal(t, pr.Headcount, retrieved.Headcount)

	// Update
	retrieved.Name = "Lead Backend Developer"
	retrieved.Level = RoleLevelLead
	retrieved.Headcount = 1
	result = db.Save(&retrieved)
	assert.NoError(t, result.Error)

	// Verify update
	var updated ProjectRole
	result = db.First(&updated, pr.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Lead Backend Developer", updated.Name)
	assert.Equal(t, uint(RoleLevelLead), updated.Level)
	assert.Equal(t, 1, updated.Headcount)

	// Delete
	result = db.Delete(&updated)
	assert.NoError(t, result.Error)

	// Verify deletion
	var deleted ProjectRole
	result = db.First(&deleted, pr.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestProjectRoleUniqueConstraint(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	client := createTestClientForProjectRole(t, db)
	project := createTestProjectForProjectRole(t, db, client.ID)

	// Create first project role
	pr1 := ProjectRole{
		ProjectID: project.ID,
		Name:      "Developer",
		Level:     RoleLevelMid,
		Headcount: 2,
	}
	result := db.Create(&pr1)
	assert.NoError(t, result.Error)

	t.Run("Same name different level is allowed", func(t *testing.T) {
		pr2 := ProjectRole{
			ProjectID: project.ID,
			Name:      "Developer",
			Level:     RoleLevelSenior, // Different level
			Headcount: 1,
		}
		result := db.Create(&pr2)
		assert.NoError(t, result.Error)
	})

	t.Run("Same level different name is allowed", func(t *testing.T) {
		pr3 := ProjectRole{
			ProjectID: project.ID,
			Name:      "QA Engineer", // Different name
			Level:     RoleLevelMid,
			Headcount: 1,
		}
		result := db.Create(&pr3)
		assert.NoError(t, result.Error)
	})

	t.Run("Different project same name and level is allowed", func(t *testing.T) {
		project2 := &Project{Name: "Project 2", ClientID: client.ID, Status: ProjectStatusActive}
		db.Create(project2)

		pr4 := ProjectRole{
			ProjectID: project2.ID, // Different project
			Name:      "Developer",
			Level:     RoleLevelMid,
			Headcount: 1,
		}
		result := db.Create(&pr4)
		assert.NoError(t, result.Error)
	})
}

func TestProjectRoleQueryOperations(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	client := createTestClientForProjectRole(t, db)
	project1 := createTestProjectForProjectRole(t, db, client.ID)
	project2 := &Project{Name: "Project 2", ClientID: client.ID, Status: ProjectStatusActive}
	db.Create(project2)

	// Create multiple project roles
	prs := []ProjectRole{
		{ProjectID: project1.ID, Name: "Backend Developer", Level: RoleLevelJunior, Headcount: 2},
		{ProjectID: project1.ID, Name: "Backend Developer", Level: RoleLevelMid, Headcount: 3},
		{ProjectID: project1.ID, Name: "Backend Developer", Level: RoleLevelSenior, Headcount: 1},
		{ProjectID: project1.ID, Name: "QA Engineer", Level: RoleLevelMid, Headcount: 2},
		{ProjectID: project2.ID, Name: "Frontend Developer", Level: RoleLevelMid, Headcount: 2},
	}

	for i := range prs {
		db.Create(&prs[i])
	}

	// Query all project roles
	var allPRs []ProjectRole
	result := db.Find(&allPRs)
	assert.NoError(t, result.Error)
	assert.Len(t, allPRs, 5)

	// Query by project ID
	var projectPRs []ProjectRole
	result = db.Where("project_id = ?", project1.ID).Find(&projectPRs)
	assert.NoError(t, result.Error)
	assert.Len(t, projectPRs, 4)

	// Query by name
	var devPRs []ProjectRole
	result = db.Where("name = ?", "Backend Developer").Find(&devPRs)
	assert.NoError(t, result.Error)
	assert.Len(t, devPRs, 3)

	// Query by level
	var midPRs []ProjectRole
	result = db.Where("level = ?", RoleLevelMid).Find(&midPRs)
	assert.NoError(t, result.Error)
	assert.Len(t, midPRs, 3)

	// Query by project and name
	var proj1DevPRs []ProjectRole
	result = db.Where("project_id = ? AND name = ?", project1.ID, "Backend Developer").Find(&proj1DevPRs)
	assert.NoError(t, result.Error)
	assert.Len(t, proj1DevPRs, 3)

	// Query by project, name, and level (should be unique)
	var uniquePR ProjectRole
	result = db.Where("project_id = ? AND name = ? AND level = ?", project1.ID, "Backend Developer", RoleLevelSenior).First(&uniquePR)
	assert.NoError(t, result.Error)
	assert.Equal(t, 1, uniquePR.Headcount)
}

func TestProjectRoleRelationships(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	client := createTestClientForProjectRole(t, db)
	project := createTestProjectForProjectRole(t, db, client.ID)

	// Create project role
	pr := ProjectRole{
		ProjectID: project.ID,
		Name:      "Developer",
		Level:     RoleLevelMid,
		Headcount: 2,
	}
	result := db.Create(&pr)
	assert.NoError(t, result.Error)

	// Load project role with project relationship
	var loadedPR ProjectRole
	result = db.Preload("Project").First(&loadedPR, pr.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, loadedPR.Project)
	assert.Equal(t, project.ID, loadedPR.Project.ID)
	assert.Equal(t, project.Name, loadedPR.Project.Name)
}

func TestProjectRoleAllowedSortFields(t *testing.T) {
	expectedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"project_id": "project_id",
		"level":      "level",
		"headcount":  "headcount",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	assert.Equal(t, expectedFields, ProjectRoleAllowedSortField)
}

func BenchmarkProjectRoleValidate(b *testing.B) {
	pr := ProjectRole{
		ProjectID: 1,
		Name:      "Developer",
		Level:     RoleLevelMid,
		Headcount: 2,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pr.Validate()
	}
}

func BenchmarkProjectRoleCreate(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Client{}, &Project{}, &ProjectRole{}) //nolint:errcheck

	client := &Client{Name: "Test", Email: "test@test.com", Status: ClientStatusActive}
	db.Create(client)
	project := &Project{Name: "Test", ClientID: client.ID, Status: ProjectStatusActive}
	db.Create(project)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pr := ProjectRole{
			ProjectID: project.ID,
			Name:      "Developer",
			Level:     RoleLevelMid,
			Headcount: 2,
		}
		db.Create(&pr)
	}
}
