package repositories

import (
	"context"
	"testing"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHumanResourceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the HumanResource table
	err = db.AutoMigrate(&entities.HumanResource{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestNewHumanResourceRepository(t *testing.T) {
	db := setupHumanResourceTestDB(t)
	repo := NewHRRepository(db)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestHumanResourceRepository_Create(t *testing.T) {
	db := setupHumanResourceTestDB(t)
	repo := NewHRRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		humanResource *entities.HumanResource
		wantError     bool
		checkErr      func(*testing.T, error)
	}{
		{
			name: "Valid HumanResource creation",
			humanResource: &entities.HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: entities.HumanResourceStatusActive,
			},
			wantError: false,
		},
		{
			name: "HumanResource with all fields",
			humanResource: &entities.HumanResource{
				Name:   "Jane Smith",
				Title:  "Lead Developer",
				Level:  "Lead",
				Status: entities.HumanResourceStatusActive,
			},
			wantError: false,
		},
		{
			name: "Invalid HumanResource - empty name",
			humanResource: &entities.HumanResource{
				Name:   "",
				Title:  "Senior Developer",
				Level:  "Senior",
				Status: entities.HumanResourceStatusActive,
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "human resource name is required")
			},
		},
		{
			name: "Invalid HumanResource - empty title",
			humanResource: &entities.HumanResource{
				Name:   "John Doe",
				Title:  "",
				Level:  "Senior",
				Status: entities.HumanResourceStatusActive,
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "human resource title is required")
			},
		},
		{
			name: "Invalid HumanResource - empty level",
			humanResource: &entities.HumanResource{
				Name:   "John Doe",
				Title:  "Senior Developer",
				Level:  "",
				Status: entities.HumanResourceStatusActive,
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "human resource level is required")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(ctx, tt.humanResource)

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

				// Verify the returned HumanResource is the same instance
				assert.Equal(t, tt.humanResource, result)

				// Verify input HumanResource was also updated with generated fields
				assert.NotZero(t, tt.humanResource.ID, "Original HumanResource should have ID populated")
				assert.NotZero(t, tt.humanResource.CreatedAt, "Original HumanResource should have CreatedAt populated")
				assert.NotZero(t, tt.humanResource.UpdatedAt, "Original HumanResource should have UpdatedAt populated")
			}
		})
	}
}

func TestHumanResourceRepository_GetOne(t *testing.T) {
	db := setupHumanResourceTestDB(t)
	repo := NewHRRepository(db)
	ctx := context.Background()

	// Create a test HumanResource
	humanResource := &entities.HumanResource{
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: entities.HumanResourceStatusActive,
	}
	created, err := repo.Create(ctx, humanResource)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		id        uint
		wantError bool
		checkErr  func(*testing.T, error)
	}{
		{
			name:      "Get existing HumanResource",
			id:        created.ID,
			wantError: false,
		},
		{
			name:      "Get non-existent HumanResource",
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
				assert.Equal(t, "John Doe", result.Name)
				assert.Equal(t, "Senior Developer", result.Title)
			}
		})
	}
}

func TestHumanResourceRepository_GetMany(t *testing.T) {
	db := setupHumanResourceTestDB(t)
	repo := NewHRRepository(db)
	ctx := context.Background()

	// Create test HRs
	humanResources := []*entities.HumanResource{
		{
			Name:   "John Doe",
			Title:  "Senior Developer",
			Level:  "Senior",
			Status: entities.HumanResourceStatusActive,
		},
		{
			Name:   "Jane Smith",
			Title:  "Lead Developer",
			Level:  "Lead",
			Status: entities.HumanResourceStatusInactive,
		},
		{
			Name:   "Bob Johnson",
			Title:  "Junior Developer",
			Level:  "Junior",
			Status: entities.HumanResourceStatusActive,
		},
		{
			Name:   "Alice Williams",
			Title:  "Senior Engineer",
			Level:  "Senior",
			Status: entities.HumanResourceStatusActive,
		},
	}

	var createdIDs []uint
	for _, h := range humanResources {
		created, err := repo.Create(ctx, h)
		assert.NoError(t, err)
		createdIDs = append(createdIDs, created.ID)
	}

	tests := []struct {
		name          string
		qParams       *entities.HumanResourceQueryParams
		wantCount     int64
		wantResultLen int
		checkFunc     func(*testing.T, []*entities.HumanResource)
	}{
		{
			name:          "Get all HRs (nil params)",
			qParams:       nil,
			wantCount:     4,
			wantResultLen: 4,
		},
		{
			name: "Filter by status - Active",
			qParams: &entities.HumanResourceQueryParams{
				Status: entities.HumanResourceStatusActive,
			},
			wantCount:     3,
			wantResultLen: 3,
			checkFunc: func(t *testing.T, results []*entities.HumanResource) {
				for _, h := range results {
					assert.Equal(t, uint(entities.HumanResourceStatusActive), h.Status)
				}
			},
		},
		{
			name: "Filter by status - Inactive",
			qParams: &entities.HumanResourceQueryParams{
				Status: entities.HumanResourceStatusInactive,
			},
			wantCount:     1,
			wantResultLen: 1,
			checkFunc: func(t *testing.T, results []*entities.HumanResource) {
				assert.Equal(t, "Jane Smith", results[0].Name)
			},
		},
		{
			name: "Filter by ID_In",
			qParams: &entities.HumanResourceQueryParams{
				ID_In: []uint{createdIDs[0], createdIDs[2]},
			},
			wantCount:     2,
			wantResultLen: 2,
		},
		{
			name: "Filter by Name",
			qParams: &entities.HumanResourceQueryParams{
				Name: "John Doe",
			},
			wantCount:     1,
			wantResultLen: 1,
			checkFunc: func(t *testing.T, results []*entities.HumanResource) {
				assert.Equal(t, "John Doe", results[0].Name)
			},
		},
		{
			name: "Filter by Name_Like",
			qParams: &entities.HumanResourceQueryParams{
				Name_Like: "John",
			},
			wantCount:     2,
			wantResultLen: 2,
		},
		{
			name: "Filter by Title",
			qParams: &entities.HumanResourceQueryParams{
				Title: "Senior Developer",
			},
			wantCount:     1,
			wantResultLen: 1,
		},
		{
			name: "Filter by Title_Like",
			qParams: &entities.HumanResourceQueryParams{
				Title_Like: "Developer",
			},
			wantCount:     3,
			wantResultLen: 3,
		},
		{
			name: "Filter by Level",
			qParams: &entities.HumanResourceQueryParams{
				Level: "Senior",
			},
			wantCount:     2,
			wantResultLen: 2,
		},
		{
			name: "Filter by Level_Like",
			qParams: &entities.HumanResourceQueryParams{
				Level_Like: "Senior",
			},
			wantCount:     2,
			wantResultLen: 2,
		},
		{
			name: "Filter by Status_In",
			qParams: &entities.HumanResourceQueryParams{
				Status_In: []uint{entities.HumanResourceStatusActive},
			},
			wantCount:     3,
			wantResultLen: 3,
		},
		{
			name: "With pagination",
			qParams: &entities.HumanResourceQueryParams{
				QueryParams: &entities.QueryParams{
					Pagination: entities.NewPagination(1, 2),
				},
			},
			wantCount:     4,
			wantResultLen: 2,
		},
		{
			name: "With sorting by name",
			qParams: &entities.HumanResourceQueryParams{
				QueryParams: &entities.QueryParams{
					Sorts: []*entities.Sort{
						entities.NewSort("name", entities.SortOrderAsc),
					},
				},
			},
			wantCount:     4,
			wantResultLen: 4,
			checkFunc: func(t *testing.T, results []*entities.HumanResource) {
				assert.Equal(t, "Alice Williams", results[0].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, count, err := repo.GetMany(ctx, tt.qParams)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCount, count)
			assert.Len(t, results, tt.wantResultLen)

			if tt.checkFunc != nil {
				tt.checkFunc(t, results)
			}
		})
	}
}

func TestHumanResourceRepository_Update(t *testing.T) {
	db := setupHumanResourceTestDB(t)
	repo := NewHRRepository(db)
	ctx := context.Background()

	// Create a test HumanResource
	humanResource := &entities.HumanResource{
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: entities.HumanResourceStatusActive,
	}
	created, err := repo.Create(ctx, humanResource)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		updateFunc   func(*entities.HumanResource)
		wantError    bool
		checkErr     func(*testing.T, error)
		wantAffected int64
	}{
		{
			name: "Valid update",
			updateFunc: func(h *entities.HumanResource) {
				h.Name = "Jane Smith"
				h.Title = "Lead Developer"
			},
			wantError:    false,
			wantAffected: 1,
		},
		{
			name: "Update status to inactive",
			updateFunc: func(h *entities.HumanResource) {
				h.Status = entities.HumanResourceStatusInactive
			},
			wantError:    false,
			wantAffected: 1,
		},
		{
			name: "Invalid update - empty name",
			updateFunc: func(h *entities.HumanResource) {
				h.Name = ""
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "human resource name is required")
			},
		},
		{
			name: "Invalid update - empty title",
			updateFunc: func(h *entities.HumanResource) {
				h.Title = ""
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "human resource title is required")
			},
		},
		{
			name: "Invalid update - empty level",
			updateFunc: func(h *entities.HumanResource) {
				h.Level = ""
			},
			wantError: true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "human resource level is required")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get a fresh copy from database
			hrToUpdate, err := repo.GetOne(ctx, created.ID)
			assert.NoError(t, err)

			// Apply update
			tt.updateFunc(hrToUpdate)

			// Attempt update
			affected, err := repo.Update(ctx, hrToUpdate)

			if tt.wantError {
				assert.Error(t, err)
				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAffected, affected)

				// Verify update in database
				updated, err := repo.GetOne(ctx, created.ID)
				assert.NoError(t, err)
				assert.Equal(t, hrToUpdate.Name, updated.Name)
				assert.Equal(t, hrToUpdate.Title, updated.Title)
				assert.Equal(t, hrToUpdate.Level, updated.Level)
				assert.Equal(t, hrToUpdate.Status, updated.Status)
			}
		})
	}
}

func TestHumanResourceRepository_Update_NonExistent(t *testing.T) {
	db := setupHumanResourceTestDB(t)
	repo := NewHRRepository(db)
	ctx := context.Background()

	// Try to update non-existent HumanResource
	humanResource := &entities.HumanResource{
		ID:     99999,
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: entities.HumanResourceStatusActive,
	}

	affected, err := repo.Update(ctx, humanResource)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), affected)
}

func TestHumanResourceRepository_Delete(t *testing.T) {
	db := setupHumanResourceTestDB(t)
	repo := NewHRRepository(db)
	ctx := context.Background()

	// Create a test HumanResource
	humanResource := &entities.HumanResource{
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: entities.HumanResourceStatusActive,
	}
	created, err := repo.Create(ctx, humanResource)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		id        uint
		wantError bool
		checkErr  func(*testing.T, error)
	}{
		{
			name:      "Delete existing HumanResource",
			id:        created.ID,
			wantError: false,
		},
		{
			name:      "Delete non-existent HumanResource",
			id:        99999,
			wantError: false, // GORM doesn't error on delete of non-existent record
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.id)

			if tt.wantError {
				assert.Error(t, err)
				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
			} else {
				assert.NoError(t, err)

				// Verify deletion (only for the first test case)
				if tt.id == created.ID {
					result, err := repo.GetOne(ctx, tt.id)
					assert.Error(t, err)
					assert.Nil(t, result)
					assert.Equal(t, entities.ErrRecordNotFound, err)
				}
			}
		})
	}
}

func TestHumanResourceRepository_CRUDFlow(t *testing.T) {
	db := setupHumanResourceTestDB(t)
	repo := NewHRRepository(db)
	ctx := context.Background()

	// Create
	humanResource := &entities.HumanResource{
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: entities.HumanResourceStatusActive,
	}
	created, err := repo.Create(ctx, humanResource)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// Read
	retrieved, err := repo.GetOne(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.Name, retrieved.Name)

	// Update
	retrieved.Name = "Jane Smith"
	affected, err := repo.Update(ctx, retrieved)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	// Verify update
	updated, err := repo.GetOne(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Smith", updated.Name)

	// Delete
	err = repo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	// Verify deletion
	deleted, err := repo.GetOne(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, deleted)
	assert.Equal(t, entities.ErrRecordNotFound, err)
}

func BenchmarkHRRepository_Create(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) //nolint:errcheck
	_ = db.AutoMigrate(&entities.HumanResource{})               //nolint:errcheck
	repo := NewHRRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		humanResource := &entities.HumanResource{
			Name:   "John Doe",
			Title:  "Senior Developer",
			Level:  "Senior",
			Status: entities.HumanResourceStatusActive,
		}
		_, _ = repo.Create(ctx, humanResource) //nolint:errcheck
	}
}

func BenchmarkHRRepository_GetOne(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) //nolint:errcheck
	_ = db.AutoMigrate(&entities.HumanResource{})               //nolint:errcheck
	repo := NewHRRepository(db)
	ctx := context.Background()

	// Create test data
	humanResource := &entities.HumanResource{
		Name:   "John Doe",
		Title:  "Senior Developer",
		Level:  "Senior",
		Status: entities.HumanResourceStatusActive,
	}
	created, _ := repo.Create(ctx, humanResource) //nolint:errcheck

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetOne(ctx, created.ID) //nolint:errcheck
	}
}

func BenchmarkHRRepository_GetMany(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) //nolint:errcheck
	_ = db.AutoMigrate(&entities.HumanResource{})               //nolint:errcheck
	repo := NewHRRepository(db)
	ctx := context.Background()

	// Create test data
	for i := 0; i < 100; i++ {
		humanResource := &entities.HumanResource{
			Name:   "John Doe",
			Title:  "Senior Developer",
			Level:  "Senior",
			Status: entities.HumanResourceStatusActive,
		}
		_, _ = repo.Create(ctx, humanResource) //nolint:errcheck
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = repo.GetMany(ctx, nil) //nolint:errcheck
	}
}
