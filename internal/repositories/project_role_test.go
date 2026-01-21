package repositories

import (
	"context"
	"testing"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProjectRoleTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all required tables
	err = db.AutoMigrate(&entities.Client{}, &entities.Project{}, &entities.ProjectRole{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func createTestClientForProjectRoleRepo(t *testing.T, db *gorm.DB) *entities.Client {
	client := &entities.Client{
		Name:   "Test Client",
		Email:  "test@client.com",
		Status: entities.ClientStatusActive,
	}
	result := db.Create(client)
	assert.NoError(t, result.Error)
	return client
}

func createTestProjectForProjectRoleRepo(t *testing.T, db *gorm.DB, clientID uint) *entities.Project {
	project := &entities.Project{
		Name:     "Test Project",
		ClientID: clientID,
		Status:   entities.ProjectStatusActive,
	}
	result := db.Create(project)
	assert.NoError(t, result.Error)
	return project
}

func TestNewProjectRoleRepository(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	repo := NewProjectRoleRepository(db)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestProjectRoleRepository_Create(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	client := createTestClientForProjectRoleRepo(t, db)
	project := createTestProjectForProjectRoleRepo(t, db, client.ID)

	tests := []struct {
		name        string
		projectRole *entities.ProjectRole
		wantError   bool
		checkErr    func(*testing.T, error)
	}{
		{
			name: "Valid project role creation",
			projectRole: &entities.ProjectRole{
				ProjectID: project.ID,
				Name:      "Backend Developer",
				Level:     entities.RoleLevelMid,
				Headcount: 2,
			},
			wantError: false,
		},
		{
			name: "Project role with all levels",
			projectRole: &entities.ProjectRole{
				ProjectID: project.ID,
				Name:      "Frontend Developer",
				Level:     entities.RoleLevelSenior,
				Headcount: 1,
			},
			wantError: false,
		},
		{
			name: "Invalid project role - empty name",
			projectRole: &entities.ProjectRole{
				ProjectID: project.ID,
				Name:      "",
				Level:     entities.RoleLevelMid,
				Headcount: 1,
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "project role name is required")
			},
		},
		{
			name: "Invalid project role - missing project ID",
			projectRole: &entities.ProjectRole{
				ProjectID: 0,
				Name:      "Developer",
				Level:     entities.RoleLevelMid,
				Headcount: 1,
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "project role must belong to a project")
			},
		},
		{
			name: "Invalid project role - negative headcount",
			projectRole: &entities.ProjectRole{
				ProjectID: project.ID,
				Name:      "QA",
				Level:     entities.RoleLevelMid,
				Headcount: -1,
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "headcount must be non-negative")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(ctx, tt.projectRole)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Verify database-generated fields are populated
				assert.NotZero(t, result.ID, "ID should be auto-generated")
				assert.NotZero(t, result.CreatedAt, "CreatedAt should be auto-generated")
				assert.NotZero(t, result.UpdatedAt, "UpdatedAt should be auto-generated")

				// Verify the returned project role is the same instance
				assert.Equal(t, tt.projectRole, result)

				// Verify input project role was also updated with generated fields
				assert.NotZero(t, tt.projectRole.ID, "Original project role should have ID populated")
				assert.NotZero(t, tt.projectRole.CreatedAt, "Original project role should have CreatedAt populated")
				assert.NotZero(t, tt.projectRole.UpdatedAt, "Original project role should have UpdatedAt populated")
			}
		})
	}
}

func TestProjectRoleRepository_GetOne(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	client := createTestClientForProjectRoleRepo(t, db)
	project := createTestProjectForProjectRoleRepo(t, db, client.ID)

	// Create a test project role
	projectRole := &entities.ProjectRole{
		ProjectID: project.ID,
		Name:      "Backend Developer",
		Level:     entities.RoleLevelMid,
		Headcount: 2,
	}
	created, err := repo.Create(ctx, projectRole)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		id        uint
		wantError bool
		checkErr  func(*testing.T, error)
	}{
		{
			name:      "Get existing project role",
			id:        created.ID,
			wantError: false,
		},
		{
			name:      "Get non-existent project role",
			id:        99999,
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, entities.ErrRecordNotFound, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetOne(ctx, tt.id)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
				assert.Equal(t, "Backend Developer", result.Name)
				assert.Equal(t, uint(entities.RoleLevelMid), result.Level)
				assert.Equal(t, 2, result.Headcount)
			}
		})
	}
}

func TestProjectRoleRepository_GetMany(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	client := createTestClientForProjectRoleRepo(t, db)
	project1 := createTestProjectForProjectRoleRepo(t, db, client.ID)
	project2 := &entities.Project{Name: "Project 2", ClientID: client.ID, Status: entities.ProjectStatusActive}
	db.Create(project2)

	// Create test project roles
	projectRoles := []*entities.ProjectRole{
		{ProjectID: project1.ID, Name: "Backend Developer", Level: entities.RoleLevelJunior, Headcount: 2},
		{ProjectID: project1.ID, Name: "Backend Developer", Level: entities.RoleLevelMid, Headcount: 3},
		{ProjectID: project1.ID, Name: "Backend Developer", Level: entities.RoleLevelSenior, Headcount: 1},
		{ProjectID: project1.ID, Name: "QA Engineer", Level: entities.RoleLevelMid, Headcount: 2},
		{ProjectID: project2.ID, Name: "Frontend Developer", Level: entities.RoleLevelMid, Headcount: 2},
	}

	var createdIDs []uint
	for _, pr := range projectRoles {
		created, err := repo.Create(ctx, pr)
		assert.NoError(t, err)
		createdIDs = append(createdIDs, created.ID)
	}

	tests := []struct {
		name          string
		qParams       *entities.ProjectRoleQueryParams
		wantCount     int64
		wantResultLen int
		checkFunc     func(*testing.T, []*entities.ProjectRole)
	}{
		{
			name:          "Get all project roles (nil params)",
			qParams:       nil,
			wantCount:     5,
			wantResultLen: 5,
		},
		{
			name: "Filter by project ID",
			qParams: &entities.ProjectRoleQueryParams{
				ProjectID: project1.ID,
			},
			wantCount:     4,
			wantResultLen: 4,
			checkFunc: func(t *testing.T, results []*entities.ProjectRole) {
				for _, pr := range results {
					assert.Equal(t, project1.ID, pr.ProjectID)
				}
			},
		},
		{
			name: "Filter by name",
			qParams: &entities.ProjectRoleQueryParams{
				Name: "Backend Developer",
			},
			wantCount:     3,
			wantResultLen: 3,
			checkFunc: func(t *testing.T, results []*entities.ProjectRole) {
				for _, pr := range results {
					assert.Equal(t, "Backend Developer", pr.Name)
				}
			},
		},
		{
			name: "Filter by Name_Like",
			qParams: &entities.ProjectRoleQueryParams{
				Name_Like: "Developer",
			},
			wantCount:     4,
			wantResultLen: 4,
		},
		{
			name: "Filter by level",
			qParams: &entities.ProjectRoleQueryParams{
				Level: entities.RoleLevelMid,
			},
			wantCount:     3,
			wantResultLen: 3,
			checkFunc: func(t *testing.T, results []*entities.ProjectRole) {
				for _, pr := range results {
					assert.Equal(t, uint(entities.RoleLevelMid), pr.Level)
				}
			},
		},
		{
			name: "Filter by Level_In",
			qParams: &entities.ProjectRoleQueryParams{
				Level_In: []uint{entities.RoleLevelJunior, entities.RoleLevelSenior},
			},
			wantCount:     2,
			wantResultLen: 2,
		},
		{
			name: "Filter by ID_In",
			qParams: &entities.ProjectRoleQueryParams{
				ID_In: []uint{createdIDs[0], createdIDs[2]},
			},
			wantCount:     2,
			wantResultLen: 2,
		},
		{
			name: "Filter by ProjectID_In",
			qParams: &entities.ProjectRoleQueryParams{
				ProjectID_In: []uint{project1.ID, project2.ID},
			},
			wantCount:     5,
			wantResultLen: 5,
		},
		{
			name: "Filter by Headcount_Gte",
			qParams: func() *entities.ProjectRoleQueryParams {
				headcount := 2
				return &entities.ProjectRoleQueryParams{
					Headcount_Gte: &headcount,
				}
			}(),
			wantCount:     4,
			wantResultLen: 4,
		},
		{
			name: "Filter by Headcount_Lte",
			qParams: func() *entities.ProjectRoleQueryParams {
				headcount := 2
				return &entities.ProjectRoleQueryParams{
					Headcount_Lte: &headcount,
				}
			}(),
			wantCount:     4,
			wantResultLen: 4,
		},
		{
			name: "Pagination - page 1",
			qParams: func() *entities.ProjectRoleQueryParams {
				qp := &entities.ProjectRoleQueryParams{
					QueryParams: &entities.QueryParams{
						Pagination: entities.NewPagination(1, 2),
					},
				}
				return qp
			}(),
			wantCount:     5,
			wantResultLen: 2,
		},
		{
			name: "Pagination - page 2",
			qParams: func() *entities.ProjectRoleQueryParams {
				qp := &entities.ProjectRoleQueryParams{
					QueryParams: &entities.QueryParams{
						Pagination: entities.NewPagination(2, 2),
					},
				}
				return qp
			}(),
			wantCount:     5,
			wantResultLen: 2,
		},
		{
			name: "Sort by name ascending",
			qParams: func() *entities.ProjectRoleQueryParams {
				qp := &entities.ProjectRoleQueryParams{
					QueryParams: &entities.QueryParams{
						Sorts: []*entities.Sort{
							entities.NewSort("name", entities.SortOrderAsc),
						},
					},
				}
				return qp
			}(),
			wantCount:     5,
			wantResultLen: 5,
			checkFunc: func(t *testing.T, results []*entities.ProjectRole) {
				assert.Equal(t, "Backend Developer", results[0].Name)
			},
		},
		{
			name: "Sort by name descending",
			qParams: func() *entities.ProjectRoleQueryParams {
				qp := &entities.ProjectRoleQueryParams{
					QueryParams: &entities.QueryParams{
						Sorts: []*entities.Sort{
							entities.NewSort("name", entities.SortOrderDesc),
						},
					},
				}
				return qp
			}(),
			wantCount:     5,
			wantResultLen: 5,
			checkFunc: func(t *testing.T, results []*entities.ProjectRole) {
				assert.Equal(t, "QA Engineer", results[0].Name)
			},
		},
		{
			name: "Combined filters - project and level",
			qParams: &entities.ProjectRoleQueryParams{
				ProjectID: project1.ID,
				Level:     entities.RoleLevelMid,
			},
			wantCount:     2,
			wantResultLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, count, err := repo.GetMany(ctx, tt.qParams)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCount, count, "Count mismatch")
			assert.Equal(t, tt.wantResultLen, len(results), "Result length mismatch")

			if tt.checkFunc != nil {
				tt.checkFunc(t, results)
			}
		})
	}
}

func TestProjectRoleRepository_Update(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	client := createTestClientForProjectRoleRepo(t, db)
	project := createTestProjectForProjectRoleRepo(t, db, client.ID)

	// Create a test project role
	projectRole := &entities.ProjectRole{
		ProjectID: project.ID,
		Name:      "Backend Developer",
		Level:     entities.RoleLevelMid,
		Headcount: 2,
	}
	created, err := repo.Create(ctx, projectRole)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		updateFunc func(*entities.ProjectRole)
		wantError  bool
		checkFunc  func(*testing.T, *entities.ProjectRole)
	}{
		{
			name: "Valid update - name",
			updateFunc: func(pr *entities.ProjectRole) {
				pr.Name = "Senior Backend Developer"
			},
			wantError: false,
			checkFunc: func(t *testing.T, pr *entities.ProjectRole) {
				assert.Equal(t, "Senior Backend Developer", pr.Name)
			},
		},
		{
			name: "Valid update - level",
			updateFunc: func(pr *entities.ProjectRole) {
				pr.Level = entities.RoleLevelSenior
			},
			wantError: false,
			checkFunc: func(t *testing.T, pr *entities.ProjectRole) {
				assert.Equal(t, uint(entities.RoleLevelSenior), pr.Level)
			},
		},
		{
			name: "Valid update - headcount",
			updateFunc: func(pr *entities.ProjectRole) {
				pr.Headcount = 5
			},
			wantError: false,
			checkFunc: func(t *testing.T, pr *entities.ProjectRole) {
				assert.Equal(t, 5, pr.Headcount)
			},
		},
		{
			name: "Invalid update - empty name",
			updateFunc: func(pr *entities.ProjectRole) {
				pr.Name = ""
			},
			wantError: true,
			checkFunc: nil,
		},
		{
			name: "Invalid update - invalid level",
			updateFunc: func(pr *entities.ProjectRole) {
				pr.Level = 999
			},
			wantError: true,
			checkFunc: nil,
		},
		{
			name: "Invalid update - negative headcount",
			updateFunc: func(pr *entities.ProjectRole) {
				pr.Headcount = -1
			},
			wantError: true,
			checkFunc: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get fresh copy of the project role
			toUpdate, err := repo.GetOne(ctx, created.ID)
			assert.NoError(t, err)

			originalUpdatedAt := toUpdate.UpdatedAt

			// Apply update
			tt.updateFunc(toUpdate)

			// Perform update
			rowsAffected, err := repo.Update(ctx, toUpdate)

			if tt.wantError {
				assert.Error(t, err)
				// For validation errors, no rows should be affected
				assert.Equal(t, int64(0), rowsAffected)
			} else {
				assert.NoError(t, err)
				// Exactly one row should be affected
				assert.Equal(t, int64(1), rowsAffected)

				// Verify UpdatedAt changed
				assert.NotEqual(t, originalUpdatedAt, toUpdate.UpdatedAt, "UpdatedAt should change")

				// Run additional checks if provided
				if tt.checkFunc != nil {
					tt.checkFunc(t, toUpdate)
				}

				// Verify changes were persisted to database
				retrieved, err := repo.GetOne(ctx, created.ID)
				assert.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, retrieved)
				}
			}
		})
	}
}

func TestProjectRoleRepository_Delete(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	client := createTestClientForProjectRoleRepo(t, db)
	project := createTestProjectForProjectRoleRepo(t, db, client.ID)

	tests := []struct {
		name      string
		setupFunc func() uint
		wantError bool
	}{
		{
			name: "Delete existing project role",
			setupFunc: func() uint {
				projectRole := &entities.ProjectRole{
					ProjectID: project.ID,
					Name:      "To Be Deleted",
					Level:     entities.RoleLevelMid,
					Headcount: 1,
				}
				created, err := repo.Create(ctx, projectRole)
				if err != nil {
					t.Fatalf("Failed to create project role: %v", err)
				}
				return created.ID
			},
			wantError: false,
		},
		{
			name: "Delete non-existent project role",
			setupFunc: func() uint {
				return 99999
			},
			wantError: false, // GORM doesn't error on deleting non-existent records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.setupFunc()

			err := repo.Delete(ctx, id)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify project role was deleted
				_, err := repo.GetOne(ctx, id)
				assert.Error(t, err)
				assert.Equal(t, entities.ErrRecordNotFound, err)
			}
		})
	}
}

func TestProjectRoleRepository_GetByProjectNameAndLevel(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	client := createTestClientForProjectRoleRepo(t, db)
	project := createTestProjectForProjectRoleRepo(t, db, client.ID)

	// Create a test project role
	projectRole := &entities.ProjectRole{
		ProjectID: project.ID,
		Name:      "Backend Developer",
		Level:     entities.RoleLevelMid,
		Headcount: 2,
	}
	_, err := repo.Create(ctx, projectRole)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		projectID uint
		roleName  string
		level     uint
		wantError bool
		checkErr  func(*testing.T, error)
	}{
		{
			name:      "Get existing project role by unique key",
			projectID: project.ID,
			roleName:  "Backend Developer",
			level:     entities.RoleLevelMid,
			wantError: false,
		},
		{
			name:      "Get non-existent - wrong project ID",
			projectID: 99999,
			roleName:  "Backend Developer",
			level:     entities.RoleLevelMid,
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, entities.ErrRecordNotFound, err)
			},
		},
		{
			name:      "Get non-existent - wrong name",
			projectID: project.ID,
			roleName:  "Frontend Developer",
			level:     entities.RoleLevelMid,
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, entities.ErrRecordNotFound, err)
			},
		},
		{
			name:      "Get non-existent - wrong level",
			projectID: project.ID,
			roleName:  "Backend Developer",
			level:     entities.RoleLevelSenior,
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, entities.ErrRecordNotFound, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByProjectNameAndLevel(ctx, tt.projectID, tt.roleName, tt.level)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.projectID, result.ProjectID)
				assert.Equal(t, tt.roleName, result.Name)
				assert.Equal(t, tt.level, result.Level)
			}
		})
	}
}

func TestProjectRoleRepository_CRUD_Integration(t *testing.T) {
	db := setupProjectRoleTestDB(t)
	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	client := createTestClientForProjectRoleRepo(t, db)
	project := createTestProjectForProjectRoleRepo(t, db, client.ID)

	// Create
	projectRole := &entities.ProjectRole{
		ProjectID: project.ID,
		Name:      "Backend Developer",
		Level:     entities.RoleLevelMid,
		Headcount: 2,
	}
	created, err := repo.Create(ctx, projectRole)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	createdID := created.ID

	// Read
	retrieved, err := repo.GetOne(ctx, createdID)
	assert.NoError(t, err)
	assert.Equal(t, "Backend Developer", retrieved.Name)
	assert.Equal(t, uint(entities.RoleLevelMid), retrieved.Level)
	assert.Equal(t, 2, retrieved.Headcount)

	// Update
	retrieved.Name = "Senior Backend Developer"
	retrieved.Level = entities.RoleLevelSenior
	retrieved.Headcount = 1
	rowsAffected, err := repo.Update(ctx, retrieved)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Verify update persisted
	retrieved2, err := repo.GetOne(ctx, createdID)
	assert.NoError(t, err)
	assert.Equal(t, "Senior Backend Developer", retrieved2.Name)
	assert.Equal(t, uint(entities.RoleLevelSenior), retrieved2.Level)
	assert.Equal(t, 1, retrieved2.Headcount)

	// Delete
	err = repo.Delete(ctx, createdID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetOne(ctx, createdID)
	assert.Error(t, err)
	assert.Equal(t, entities.ErrRecordNotFound, err)
}

func BenchmarkProjectRoleRepository_Create(b *testing.B) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&entities.Client{}, &entities.Project{}, &entities.ProjectRole{}); err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	client := &entities.Client{Name: "Test", Email: "test@test.com", Status: entities.ClientStatusActive}
	db.Create(client)
	project := &entities.Project{Name: "Test", ClientID: client.ID, Status: entities.ProjectStatusActive}
	db.Create(project)

	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		projectRole := &entities.ProjectRole{
			ProjectID: project.ID,
			Name:      "Developer",
			Level:     entities.RoleLevelMid,
			Headcount: 2,
		}
		_, _ = repo.Create(ctx, projectRole)
	}
}

func BenchmarkProjectRoleRepository_Get(b *testing.B) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&entities.Client{}, &entities.Project{}, &entities.ProjectRole{}); err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	client := &entities.Client{Name: "Test", Email: "test@test.com", Status: entities.ClientStatusActive}
	db.Create(client)
	project := &entities.Project{Name: "Test", ClientID: client.ID, Status: entities.ProjectStatusActive}
	db.Create(project)

	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	projectRole := &entities.ProjectRole{
		ProjectID: project.ID,
		Name:      "Developer",
		Level:     entities.RoleLevelMid,
		Headcount: 2,
	}
	created, err := repo.Create(ctx, projectRole)
	if err != nil {
		b.Fatalf("Failed to create project role: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetOne(ctx, created.ID)
	}
}

func BenchmarkProjectRoleRepository_Update(b *testing.B) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&entities.Client{}, &entities.Project{}, &entities.ProjectRole{}); err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	client := &entities.Client{Name: "Test", Email: "test@test.com", Status: entities.ClientStatusActive}
	db.Create(client)
	project := &entities.Project{Name: "Test", ClientID: client.ID, Status: entities.ProjectStatusActive}
	db.Create(project)

	repo := NewProjectRoleRepository(db)
	ctx := context.Background()

	projectRole := &entities.ProjectRole{
		ProjectID: project.ID,
		Name:      "Developer",
		Level:     entities.RoleLevelMid,
		Headcount: 2,
	}
	created, err := repo.Create(ctx, projectRole)
	if err != nil {
		b.Fatalf("Failed to create project role: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		created.Headcount = i % 10
		_, _ = repo.Update(ctx, created)
	}
}
