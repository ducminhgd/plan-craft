package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/ducminhgd/plan-craft/internal/requires"
	"gorm.io/gorm"
)

// Example model for demonstration
type Project struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(50);default:'not_started'"`
}

func main() {
	// Initialize database
	if err := requires.InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer requires.CloseDatabase()

	// Check database health
	if err := requires.HealthCheck(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	slog.Info("Database initialized and healthy")

	// Run examples
	runBasicCRUD()
	runQueryExamples()
	runTransactionExample()
	runBatchOperations()
}

// Example 1: Basic CRUD operations
func runBasicCRUD() {
	fmt.Println("\n=== Example 1: Basic CRUD Operations ===")

	// Auto-migrate the schema
	requires.DB.AutoMigrate(&Project{})

	// Create
	project := &Project{
		Name:        "E-Commerce Platform",
		Description: "Build a modern e-commerce platform",
		Status:      "in_progress",
	}

	result := requires.DB.Create(project)
	if result.Error != nil {
		log.Printf("Failed to create project: %v", result.Error)
		return
	}
	fmt.Printf("✓ Created project with ID: %d\n", project.ID)

	// Read
	var retrieved Project
	requires.DB.First(&retrieved, project.ID)
	fmt.Printf("✓ Retrieved project: %s\n", retrieved.Name)

	// Update
	requires.DB.Model(&retrieved).Update("Status", "completed")
	fmt.Printf("✓ Updated project status to: completed\n")

	// Delete (soft delete if using gorm.Model)
	requires.DB.Delete(&retrieved)
	fmt.Printf("✓ Deleted project\n")
}

// Example 2: Query operations
func runQueryExamples() {
	fmt.Println("\n=== Example 2: Query Operations ===")

	// Create sample data
	projects := []Project{
		{Name: "Website Redesign", Status: "in_progress"},
		{Name: "Mobile App", Status: "not_started"},
		{Name: "API Development", Status: "completed"},
	}

	requires.DB.Create(&projects)

	// Find all
	var allProjects []Project
	requires.DB.Find(&allProjects)
	fmt.Printf("✓ Found %d projects\n", len(allProjects))

	// Find with condition
	var activeProjects []Project
	requires.DB.Where("status = ?", "in_progress").Find(&activeProjects)
	fmt.Printf("✓ Found %d active projects\n", len(activeProjects))

	// Find with multiple conditions
	var filtered []Project
	requires.DB.Where("status IN ?", []string{"in_progress", "not_started"}).Find(&filtered)
	fmt.Printf("✓ Found %d projects (in_progress or not_started)\n", len(filtered))

	// Count
	var count int64
	requires.DB.Model(&Project{}).Where("status = ?", "completed").Count(&count)
	fmt.Printf("✓ Completed projects count: %d\n", count)
}

// Example 3: Transaction
func runTransactionExample() {
	fmt.Println("\n=== Example 3: Transaction ===")

	err := requires.DB.Transaction(func(tx *gorm.DB) error {
		// Create multiple projects in a transaction
		projects := []Project{
			{Name: "Project A", Status: "not_started"},
			{Name: "Project B", Status: "not_started"},
		}

		for _, p := range projects {
			if err := tx.Create(&p).Error; err != nil {
				return err // Rollback on error
			}
		}

		fmt.Printf("✓ Created %d projects in transaction\n", len(projects))
		return nil // Commit
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	}
}

// Example 4: Batch operations
func runBatchOperations() {
	fmt.Println("\n=== Example 4: Batch Operations ===")

	// Create batch
	projects := make([]Project, 100)
	for i := 0; i < 100; i++ {
		projects[i] = Project{
			Name:   fmt.Sprintf("Batch Project %d", i+1),
			Status: "not_started",
		}
	}

	// Insert in batches of 10
	result := requires.DB.CreateInBatches(projects, 10)
	if result.Error != nil {
		log.Printf("Batch insert failed: %v", result.Error)
		return
	}
	fmt.Printf("✓ Inserted %d projects in batches\n", len(projects))

	// Batch update
	requires.DB.Model(&Project{}).
		Where("status = ?", "not_started").
		Update("status", "in_progress")
	fmt.Printf("✓ Updated all not_started projects to in_progress\n")
}

// Example 5: Using context
func runContextExample() {
	fmt.Println("\n=== Example 5: Using Context ===")

	ctx := context.Background()

	var projects []Project
	result := requires.DB.WithContext(ctx).
		Where("status = ?", "in_progress").
		Find(&projects)

	if result.Error != nil {
		log.Printf("Query failed: %v", result.Error)
		return
	}

	fmt.Printf("✓ Found %d projects with context\n", len(projects))
}

