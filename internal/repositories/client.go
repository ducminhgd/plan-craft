package repositories

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"gorm.io/gorm"
)

// ClientRepository defines the interface for client data operations
type ClientRepository interface {
	// Get retrieves a single client by ID
	Get(ctx context.Context, id uint) (*entities.Client, error)

	// GetList retrieves a paginated list of clients with filtering and sorting
	GetList(ctx context.Context, query *entities.ClientQuery) (int64, []*entities.Client, error)

	// Create inserts a new client
	Create(ctx context.Context, client *entities.Client) (*entities.Client, error)

	// Update modifies an existing client
	Update(ctx context.Context, client *entities.Client) (int64, error)
}

// clientRepository is the concrete implementation of ClientRepository
type clientRepository struct {
	db *gorm.DB
}

// NewClientRepository creates a new instance of ClientRepository
func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

// Get retrieves a single client by ID
func (r *clientRepository) Get(ctx context.Context, id uint) (*entities.Client, error) {
	var client entities.Client
	err := r.db.WithContext(ctx).First(&client, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// GetList retrieves a paginated list of clients with filtering and sorting
func (r *clientRepository) GetList(ctx context.Context, query *entities.ClientQuery) (int64, []*entities.Client, error) {
	var clients []*entities.Client
	var total int64

	// Start with base query
	db := r.db.WithContext(ctx).Model(&entities.Client{})

	// Apply filters
	if len(query.ID_In) > 0 {
		db = db.Where("id IN ?", query.ID_In)
	}

	if query.Name != "" {
		db = db.Where("name = ?", query.Name)
	}

	if query.Name_Like != "" {
		db = db.Where("name LIKE ?", "%"+query.Name_Like+"%")
	}

	if query.Email != "" {
		db = db.Where("email = ?", query.Email)
	}

	if query.Email_Like != "" {
		db = db.Where("email LIKE ?", "%"+query.Email_Like+"%")
	}

	if query.Phone != "" {
		db = db.Where("phone = ?", query.Phone)
	}

	if query.Phone_Like != "" {
		db = db.Where("phone LIKE ?", "%"+query.Phone_Like+"%")
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.CreatedAt_Gte != nil {
		db = db.Where("created_at >= ?", *query.CreatedAt_Gte)
	}

	if query.CreatedAt_Lte != nil {
		db = db.Where("created_at <= ?", *query.CreatedAt_Lte)
	}

	if query.UpdatedAt_Gte != nil {
		db = db.Where("updated_at >= ?", *query.UpdatedAt_Gte)
	}

	if query.UpdatedAt_Lte != nil {
		db = db.Where("updated_at <= ?", *query.UpdatedAt_Lte)
	}

	// Count total records matching the filters
	if err := db.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// Apply sorting
	allowedSortFields := query.AllowedSortFields()
	if query.Sort != nil {
		db = query.Sort.Apply(db, allowedSortFields)
	}

	// Apply pagination
	if query.Pagination != nil {
		db = query.Pagination.Apply(db)
	}

	// Execute query
	if err := db.Find(&clients).Error; err != nil {
		return 0, nil, err
	}

	return total, clients, nil
}

// Create inserts a new client
func (r *clientRepository) Create(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	// Track if we need to set IsActive to false after creation
	needsInactiveUpdate := !client.IsActive

	if err := r.db.WithContext(ctx).Create(client).Error; err != nil {
		return nil, err
	}

	// If IsActive should be false, we need to explicitly update it
	// because GORM's default:true tag creates a DB default constraint with DEFAULT true (1)
	// and Go's zero value (false) is not being inserted by GORM
	if needsInactiveUpdate {
		// Use Model().Update to avoid raw SQL
		if err := r.db.WithContext(ctx).Model(client).Update("is_active", false).Error; err != nil {
			return nil, err
		}
		// Reload the client to get the updated value
		if err := r.db.WithContext(ctx).First(client, client.ID).Error; err != nil {
			return nil, err
		}
	}

	return client, nil
}

// Update modifies an existing client
func (r *clientRepository) Update(ctx context.Context, client *entities.Client) (int64, error) {
	result := r.db.WithContext(ctx).Save(client)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
