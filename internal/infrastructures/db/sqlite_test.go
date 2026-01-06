package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ducminhgd/plan-craft/config"
)

func TestInitializeDatabase(t *testing.T) {
	// Setup: Create temporary directory for test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Override config for testing
	originalCfg := config.Cfg
	config.Cfg = config.Config{
		DB: config.DBConfig{
			DSN:         dbPath,
			JournalMode: "WAL",
			Synchronous: "NORMAL",
			ForeignKeys: "ON",
			BusyTimeout: "5000",
			CacheSize:   "-64000",
			TempStore:   "MEMORY",
			AutoVacuum:  "INCREMENTAL",
		},
		LogLevel: "ERROR",
	}
	defer func() {
		config.Cfg = originalCfg
	}()

	// Test initialization
	err := InitializeDatabase()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Verify DB is not nil
	if DB == nil {
		t.Fatal("Expected DB to be initialized, got nil")
	}

	// Verify database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Expected database file to exist at %s", dbPath)
	}

	// Cleanup
	if err := CloseDatabase(); err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	originalCfg := config.Cfg
	config.Cfg = config.Config{
		DB: config.DBConfig{
			DSN:         dbPath,
			JournalMode: "WAL",
			Synchronous: "NORMAL",
			ForeignKeys: "ON",
			BusyTimeout: "5000",
			CacheSize:   "-64000",
			TempStore:   "MEMORY",
			AutoVacuum:  "INCREMENTAL",
		},
		LogLevel: "ERROR",
	}
	defer func() {
		config.Cfg = originalCfg
	}()

	// Initialize database
	if err := InitializeDatabase(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer CloseDatabase()

	// Test health check
	err := HealthCheck()
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
}

func TestHealthCheck_NotInitialized(t *testing.T) {
	// Save current DB state
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set DB to nil to simulate uninitialized state
	DB = nil

	// Test health check on uninitialized database
	err := HealthCheck()
	if err == nil {
		t.Error("Expected error for uninitialized database, got nil")
	}

	expectedMsg := "database not initialized"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestCloseDatabase(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	originalCfg := config.Cfg
	config.Cfg = config.Config{
		DB: config.DBConfig{
			DSN:         dbPath,
			JournalMode: "WAL",
			Synchronous: "NORMAL",
			ForeignKeys: "ON",
			BusyTimeout: "5000",
			CacheSize:   "-64000",
			TempStore:   "MEMORY",
			AutoVacuum:  "INCREMENTAL",
		},
		LogLevel: "ERROR",
	}
	defer func() {
		config.Cfg = originalCfg
	}()

	// Initialize database
	if err := InitializeDatabase(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Test close
	err := CloseDatabase()
	if err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

func TestCloseDatabase_NotInitialized(t *testing.T) {
	// Save current DB state
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set DB to nil
	DB = nil

	// Test close on uninitialized database (should not error)
	err := CloseDatabase()
	if err != nil {
		t.Errorf("Expected no error for closing uninitialized database, got: %v", err)
	}
}

func TestBuildSQLiteDSN(t *testing.T) {
	dbConfig := config.DBConfig{
		DSN:         "/tmp/test.db",
		JournalMode: "WAL",
		Synchronous: "NORMAL",
		ForeignKeys: "ON",
		BusyTimeout: "5000",
		CacheSize:   "-64000",
		TempStore:   "MEMORY",
		AutoVacuum:  "INCREMENTAL",
	}

	dsn := buildSQLiteDSN(dbConfig)

	// Check that DSN starts with the path
	dbPath := dbConfig.DSN
	if dsn[:len(dbPath)] != dbPath {
		t.Errorf("Expected DSN to start with %s, got %s", dbPath, dsn)
	}

	// Check for required parameters
	requiredParams := []string{
		"_journal_mode=WAL",
		"_synchronous=NORMAL",
		"_foreign_keys=ON",
		"_busy_timeout=5000",
		"_cache_size=-64000",
		"_temp_store=MEMORY",
		"_auto_vacuum=INCREMENTAL",
	}

	for _, param := range requiredParams {
		if !contains(dsn, param) {
			t.Errorf("Expected DSN to contain %s, got %s", param, dsn)
		}
	}
}

func TestBuildSQLiteDSN_EmptyParams(t *testing.T) {
	// Test with empty parameters - should only include DSN path
	dbConfig := config.DBConfig{
		DSN: "/tmp/test.db",
		// All other fields empty
	}

	dsn := buildSQLiteDSN(dbConfig)

	// Should only have the path, no parameters
	if dsn != "/tmp/test.db" {
		t.Errorf("Expected DSN to be '/tmp/test.db', got '%s'", dsn)
	}
}

func TestBuildSQLiteDSN_PartialParams(t *testing.T) {
	// Test with only some parameters set
	dbConfig := config.DBConfig{
		DSN:         "/tmp/test.db",
		JournalMode: "WAL",
		ForeignKeys: "ON",
		// Other fields empty
	}

	dsn := buildSQLiteDSN(dbConfig)

	// Check that DSN starts with the path
	if dsn[:len("/tmp/test.db")] != "/tmp/test.db" {
		t.Errorf("Expected DSN to start with '/tmp/test.db', got '%s'", dsn)
	}

	// Check that set parameters are included
	if !contains(dsn, "_journal_mode=WAL") {
		t.Errorf("Expected DSN to contain '_journal_mode=WAL', got '%s'", dsn)
	}
	if !contains(dsn, "_foreign_keys=ON") {
		t.Errorf("Expected DSN to contain '_foreign_keys=ON', got '%s'", dsn)
	}

	// Check that unset parameters are NOT included
	if contains(dsn, "_synchronous=") {
		t.Errorf("Expected DSN to NOT contain '_synchronous=', got '%s'", dsn)
	}
	if contains(dsn, "_busy_timeout=") {
		t.Errorf("Expected DSN to NOT contain '_busy_timeout=', got '%s'", dsn)
	}
}

func TestBuildSQLiteDSN_EmptyDSN(t *testing.T) {
	// Test with empty DSN - should use default
	dbConfig := config.DBConfig{
		JournalMode: "WAL",
	}

	dsn := buildSQLiteDSN(dbConfig)

	// Should use default path
	if dsn[:len("data/plancraft.db")] != "data/plancraft.db" {
		t.Errorf("Expected DSN to start with 'data/plancraft.db', got '%s'", dsn)
	}

	// Should include the parameter
	if !contains(dsn, "_journal_mode=WAL") {
		t.Errorf("Expected DSN to contain '_journal_mode=WAL', got '%s'", dsn)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

