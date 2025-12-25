package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate the Client table
	err = db.AutoMigrate(&entities.Client{})
	assert.NoError(t, err)

	return db
}

func seedClients(t *testing.T, db *gorm.DB) []*entities.Client {
	clientsData := []struct {
		client       *entities.Client
		shouldBeInactive bool
	}{
		{
			client: &entities.Client{
				Name:     "Acme Corp",
				Email:    "contact@acme.com",
				Phone:    "123-456-7890",
				Address:  "123 Main St",
				IsActive: true,
				Notes:    "Premium client",
			},
			shouldBeInactive: false,
		},
		{
			client: &entities.Client{
				Name:     "Tech Solutions",
				Email:    "info@techsolutions.com",
				Phone:    "987-654-3210",
				Address:  "456 Tech Ave",
				IsActive: true,
				Notes:    "Regular client",
			},
			shouldBeInactive: false,
		},
		{
			client: &entities.Client{
				Name:     "Inactive Client",
				Email:    "old@client.com",
				Phone:    "555-555-5555",
				Address:  "789 Old St",
				IsActive: false,
				Notes:    "No longer active",
			},
			shouldBeInactive: true,
		},
	}

	var clients []*entities.Client
	for _, data := range clientsData {
		client := data.client

		// Create the client - GORM will use the default for IsActive
		err := db.Create(client).Error
		assert.NoError(t, err)

		// After Create, GORM updates the client with DB defaults, so check our flag
		// For inactive clients, explicitly update the IsActive field using integer 0
		// This is needed because GORM's default tag creates a DB constraint with DEFAULT true (1)
		// and Go's zero value (false) is not being inserted
		if data.shouldBeInactive {
			result := db.Exec("UPDATE clients SET is_active = 0 WHERE id = ?", client.ID)
			assert.NoError(t, result.Error)

			// Reload the client to reflect the change in our struct
			db.First(client, client.ID)
		}

		clients = append(clients, client)
	}

	return clients
}

func TestNewClientRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClientRepository(db)
	assert.NotNil(t, repo)
}

func TestClientRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClientRepository(db)

	tests := []struct {
		name    string
		client  *entities.Client
		wantErr bool
	}{
		{
			name: "Valid client",
			client: &entities.Client{
				Name:     "New Client",
				Email:    "new@client.com",
				Phone:    "111-222-3333",
				Address:  "New Address",
				IsActive: true,
			},
			wantErr: false,
		},
		{
			name: "Client without email",
			client: &entities.Client{
				Name:     "Client No Email",
				Phone:    "222-333-4444",
				IsActive: true,
			},
			wantErr: false,
		},
		{
			name: "Client with empty name (should fail validation)",
			client: &entities.Client{
				Name:  "",
				Email: "test@test.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := repo.Create(ctx, tt.client)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotZero(t, result.ID)
				assert.Equal(t, tt.client.Name, result.Name)
				assert.Equal(t, tt.client.Email, result.Email)
			}
		})
	}
}

func TestClientRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClientRepository(db)
	seededClients := seedClients(t, db)

	tests := []struct {
		name       string
		id         uint
		wantClient *entities.Client
		wantNil    bool
		wantErr    bool
	}{
		{
			name:       "Existing client",
			id:         seededClients[0].ID,
			wantClient: seededClients[0],
			wantNil:    false,
			wantErr:    false,
		},
		{
			name:       "Non-existent client",
			id:         9999,
			wantClient: nil,
			wantNil:    true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := repo.Get(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantClient.ID, result.ID)
				assert.Equal(t, tt.wantClient.Name, result.Name)
				assert.Equal(t, tt.wantClient.Email, result.Email)
			}
		})
	}
}

func TestClientRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClientRepository(db)
	seededClients := seedClients(t, db)

	tests := []struct {
		name            string
		client          *entities.Client
		wantRowsUpdated int64
		wantErr         bool
	}{
		{
			name: "Update existing client",
			client: &entities.Client{
				ID:       seededClients[0].ID,
				Name:     "Updated Acme Corp",
				Email:    "updated@acme.com",
				Phone:    "999-888-7777",
				Address:  "New Address",
				IsActive: true,
			},
			wantRowsUpdated: 1,
			wantErr:         false,
		},
		{
			name: "Update with empty name (should fail validation)",
			client: &entities.Client{
				ID:    seededClients[1].ID,
				Name:  "",
				Email: "test@test.com",
			},
			wantRowsUpdated: 0,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rowsAffected, err := repo.Update(ctx, tt.client)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRowsUpdated, rowsAffected)

				// Verify the update
				updated, getErr := repo.Get(ctx, tt.client.ID)
				assert.NoError(t, getErr)
				assert.Equal(t, tt.client.Name, updated.Name)
				assert.Equal(t, tt.client.Email, updated.Email)
			}
		})
	}
}

func TestClientRepository_GetList(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClientRepository(db)
	seedClients(t, db)

	tests := []struct {
		name       string
		query      *entities.ClientQuery
		wantTotal  int64
		wantCount  int
		wantErr    bool
		checkFirst func(*testing.T, *entities.Client)
	}{
		{
			name: "Get all clients without filters",
			query: &entities.ClientQuery{
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
				},
			},
			wantTotal: 3,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "Filter by active status",
			query: &entities.ClientQuery{
				IsActive: func() *bool { b := true; return &b }(),
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
				},
			},
			wantTotal: 2,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "Filter by inactive status",
			query: &entities.ClientQuery{
				IsActive: func() *bool { b := false; return &b }(),
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
				},
			},
			wantTotal: 1,
			wantCount: 1,
			wantErr:   false,
			checkFirst: func(t *testing.T, c *entities.Client) {
				assert.Equal(t, "Inactive Client", c.Name)
				assert.False(t, c.IsActive)
			},
		},
		{
			name: "Filter by name like",
			query: &entities.ClientQuery{
				Name_Like: "Tech",
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
				},
			},
			wantTotal: 1,
			wantCount: 1,
			wantErr:   false,
			checkFirst: func(t *testing.T, c *entities.Client) {
				assert.Equal(t, "Tech Solutions", c.Name)
			},
		},
		{
			name: "Filter by email like",
			query: &entities.ClientQuery{
				Email_Like: "acme",
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
				},
			},
			wantTotal: 1,
			wantCount: 1,
			wantErr:   false,
			checkFirst: func(t *testing.T, c *entities.Client) {
				assert.Contains(t, c.Email, "acme")
			},
		},
		{
			name: "Filter by exact name",
			query: &entities.ClientQuery{
				Name: "Acme Corp",
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
				},
			},
			wantTotal: 1,
			wantCount: 1,
			wantErr:   false,
			checkFirst: func(t *testing.T, c *entities.Client) {
				assert.Equal(t, "Acme Corp", c.Name)
			},
		},
		{
			name: "Sort by name ascending",
			query: &entities.ClientQuery{
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
					Sort:       entities.NewSort("name", entities.SortOrderAsc),
				},
			},
			wantTotal: 3,
			wantCount: 3,
			wantErr:   false,
			checkFirst: func(t *testing.T, c *entities.Client) {
				assert.Equal(t, "Acme Corp", c.Name)
			},
		},
		{
			name: "Sort by name descending",
			query: &entities.ClientQuery{
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
					Sort:       entities.NewSort("name", entities.SortOrderDesc),
				},
			},
			wantTotal: 3,
			wantCount: 3,
			wantErr:   false,
			checkFirst: func(t *testing.T, c *entities.Client) {
				assert.Equal(t, "Tech Solutions", c.Name)
			},
		},
		{
			name: "Pagination: page 1, size 2",
			query: &entities.ClientQuery{
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 2),
					Sort:       entities.NewSort("id", entities.SortOrderAsc),
				},
			},
			wantTotal: 3,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "Pagination: page 2, size 2",
			query: &entities.ClientQuery{
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(2, 2),
					Sort:       entities.NewSort("id", entities.SortOrderAsc),
				},
			},
			wantTotal: 3,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "Filter by created_at range",
			query: &entities.ClientQuery{
				CreatedAt_Gte: func() *time.Time { t := time.Now().Add(-24 * time.Hour); return &t }(),
				CreatedAt_Lte: func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
				},
			},
			wantTotal: 3,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "Filter by ID_In",
			query: &entities.ClientQuery{
				ID_In: []uint{1, 2},
				QueryParams: entities.QueryParams{
					Pagination: entities.NewPagination(1, 10),
				},
			},
			wantTotal: 2,
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			total, clients, err := repo.GetList(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTotal, total)
				assert.Equal(t, tt.wantCount, len(clients))

				if tt.checkFirst != nil && len(clients) > 0 {
					tt.checkFirst(t, clients[0])
				}
			}
		})
	}
}

func TestClientRepository_GetList_CombinedFilters(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClientRepository(db)
	seedClients(t, db)

	// Test combining multiple filters
	query := &entities.ClientQuery{
		IsActive:   func() *bool { b := true; return &b }(),
		Name_Like:  "Corp",
		Email_Like: "acme",
		QueryParams: entities.QueryParams{
			Pagination: entities.NewPagination(1, 10),
			Sort:       entities.NewSort("name", entities.SortOrderAsc),
		},
	}

	ctx := context.Background()
	total, clients, err := repo.GetList(ctx, query)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(clients))
	assert.Equal(t, "Acme Corp", clients[0].Name)
	assert.True(t, clients[0].IsActive)
}
