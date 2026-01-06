package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ducminhgd/plan-craft/internal"
	"github.com/ducminhgd/plan-craft/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ClientRepository is the repository for client entities
type ClientRepository struct {
	db *gorm.DB
}

// NewClientRepository creates a new client repository
func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

// Create creates a new client and returns it with database-generated fields populated
func (r *ClientRepository) Create(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	err := r.db.WithContext(ctx).Create(client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "client", "method", "Create", "error", err)
			return nil, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "client", "method", "Create", "error", err)
			return nil, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			internal.Logger.Error("duplicated key", "repository", "client", "method", "Create", "error", err)
			return nil, entities.ErrDuplicatedKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "client", "method", "Create", "error", err)
			return nil, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "client", "method", "Create", "error", err)
			return nil, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to create client", "repository", "client", "method", "Create", "error", err)
		return nil, err
	}
	return client, nil
}

// GetOne gets a client by ID
func (r *ClientRepository) GetOne(ctx context.Context, id uint) (*entities.Client, error) {
	var client entities.Client
	err := r.db.WithContext(ctx).Model(&entities.Client{}).First(&client, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "client", "method", "Get", "error", err)
			return nil, entities.ErrRecordNotFound
		}
		internal.Logger.Error("failed to get client", "repository", "client", "method", "Get", "error", err)
		return nil, err
	}
	return &client, err
}

// GetMany gets multiple clients by query parameters
func (r *ClientRepository) GetMany(ctx context.Context, qParams *entities.ClientQueryParams) ([]*entities.Client, int64, error) {
	var (
		clients []*entities.Client
		count   int64 = 0
	)
	q := r.db.WithContext(ctx).Model(&entities.Client{})

	if qParams == nil {
		qParams = &entities.ClientQueryParams{}
	}

	if len(qParams.ID_In) > 0 {
		q = q.Where("id IN @ID_In", sql.Named("ID_In", qParams.ID_In))
	}
	if qParams.Name != "" {
		q = q.Where("name = @Name", sql.Named("Name", qParams.Name))
	}
	if qParams.Name_Like != "" {
		q = q.Where("name LIKE @Name_Like", sql.Named("Name_Like", "%"+qParams.Name_Like+"%"))
	}
	if qParams.Email != "" {
		q = q.Where("email = @Email", sql.Named("Email", qParams.Email))
	}
	if qParams.Email_Like != "" {
		q = q.Where("email LIKE @Email_Like", sql.Named("Email_Like", "%"+qParams.Email_Like+"%"))
	}
	if qParams.Phone != "" {
		q = q.Where("phone = @Phone", sql.Named("Phone", qParams.Phone))
	}
	if qParams.Phone_Like != "" {
		q = q.Where("phone LIKE @Phone_Like", sql.Named("Phone_Like", "%"+qParams.Phone_Like+"%"))
	}
	if qParams.Address_Like != "" {
		q = q.Where("address LIKE @Address_Like", sql.Named("Address_Like", "%"+qParams.Address_Like+"%"))
	}
	if qParams.ContactPerson_Like != "" {
		q = q.Where("contact_person LIKE @ContactPerson_Like", sql.Named("ContactPerson_Like", "%"+qParams.ContactPerson_Like+"%"))
	}
	if qParams.Notes_Like != "" {
		q = q.Where("notes LIKE @Notes_Like", sql.Named("Notes_Like", "%"+qParams.Notes_Like+"%"))
	}
	if qParams.Status != entities.ClientStatusUnknown {
		q = q.Where("status = @Status", sql.Named("Status", qParams.Status))
	}
	if len(qParams.Status_In) > 0 {
		q = q.Where("status IN ?", qParams.Status_In)
	}
	if qParams.CreatedAt_Gte != nil {
		q = q.Where("created_at >= @CreatedAt_Gte", sql.Named("CreatedAt_Gte", qParams.CreatedAt_Gte))
	}
	if qParams.CreatedAt_Lte != nil {
		q = q.Where("created_at <= @CreatedAt_Lte", sql.Named("CreatedAt_Lte", qParams.CreatedAt_Lte))
	}
	if qParams.UpdatedAt_Gte != nil {
		q = q.Where("updated_at >= @UpdatedAt_Gte", sql.Named("UpdatedAt_Gte", qParams.UpdatedAt_Gte))
	}
	if qParams.UpdatedAt_Lte != nil {
		q = q.Where("updated_at <= @UpdatedAt_Lte", sql.Named("UpdatedAt_Lte", qParams.UpdatedAt_Lte))
	}

	q = q.Session(&gorm.Session{})
	result := q.Count(&count)
	if result.Error != nil {
		internal.Logger.Error("failed to count clients", "repository", "client", "method", "GetMany", "error", result.Error)
		return nil, 0, result.Error
	}

	// Apply sorting params
	if qParams.QueryParams != nil {
		if qParams.Sorts != nil {
			for _, sort := range qParams.Sorts {
				q = sort.Apply(q, entities.ClientAllowedSortField)
			}
		}
		if qParams.Pagination != nil {
			q = qParams.Pagination.Apply(q)
		}
	}

	// Execute query
	result = q.Find(&clients)
	if result.Error != nil {
		internal.Logger.Error("failed to get clients", "repository", "client", "method", "GetMany", "error", result.Error)
		return nil, count, result.Error
	}
	return clients, count, nil
}

// Update updates a client and returns it with updated database fields
func (r *ClientRepository) Update(ctx context.Context, client *entities.Client) (int64, error) {
	result := r.db.WithContext(ctx).Model(client).Clauses(clause.Returning{}).Where("id = ?", client.ID).Select("*").Updates(&client)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "client", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrInvalidData) {
			internal.Logger.Error("invalid data", "repository", "client", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrInvalidData
		}
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			internal.Logger.Error("unsupported relation", "repository", "client", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrUnsupportedRelation
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "client", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrForeignKeyViolated
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			internal.Logger.Error("check constraint violated", "repository", "client", "method", "Update", "error", err)
			return result.RowsAffected, entities.ErrCheckConstraintViolated
		}
		internal.Logger.Error("failed to update client", "repository", "client", "method", "Update", "error", err)
		return result.RowsAffected, err
	}
	return result.RowsAffected, nil
}

// Delete deletes a client by ID
func (r *ClientRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entities.Client{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Logger.Error("record not found", "repository", "client", "method", "Delete", "error", err)
			return entities.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			internal.Logger.Error("foreign key violated", "repository", "client", "method", "Delete", "error", err)
			return entities.ErrForeignKeyViolated
		}
		internal.Logger.Error("failed to delete client", "repository", "client", "method", "Delete", "error", err)
		return err
	}
	return nil
}
