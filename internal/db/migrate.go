package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GokulP48/go_gin_learning/internal/logger"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"gorm.io/gorm"
)

// RunMigrations executes all pending migrations from the migrations folder
func RunMigrations() error {
	// Create migrations table if it doesn't exist
	if err := DB.AutoMigrate(&Migrations{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get all migration files
	migrationFiles, err := getMigrationFiles("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	if len(migrationFiles) == 0 {
		logger.Infof("No migration files found")
		return nil
	}

	// Get already executed migrations
	var executedMigrations []Migrations
	DB.Find(&executedMigrations)

	executedMap := make(map[string]bool)
	for _, migration := range executedMigrations {
		executedMap[migration.Name] = true
	}

	// Execute pending migrations
	pendingMigrations := []string{}
	for _, file := range migrationFiles {
		if !executedMap[file] {
			pendingMigrations = append(pendingMigrations, file)
		}
	}

	if len(pendingMigrations) == 0 {
		logger.Infof("No pending migrations to run")
		return nil
	}

	logger.Infof("Found %d pending migrations to run", len(pendingMigrations))

	// Run migrations in transaction
	return DB.Transaction(func(tx *gorm.DB) error {
		for _, migrationFile := range pendingMigrations {
			logger.Infof("Running migration: %s", migrationFile)

			// Read migration file
			content, err := os.ReadFile(filepath.Join("migrations", migrationFile))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %v", migrationFile, err)
			}

			// Execute the SQL
			if err := tx.Exec(string(content)).Error; err != nil {
				return fmt.Errorf("failed to execute migration %s: %v", migrationFile, err)
			}

			// Record the migration
			migrationRecord := Migrations{
				Name:      migrationFile,
				Timestamp: time.Now().UnixMilli(),
			}
			if err := tx.Create(&migrationRecord).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %v", migrationFile, err)
			}

			logger.Infof("Successfully executed migration: %s", migrationFile)
		}
		return nil
	})
}

// getMigrationFiles returns sorted list of migration files
func getMigrationFiles(migrationDir string) ([]string, error) {
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Sort files to ensure consistent execution order
	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

func CreateMigration(name string) error {
	// Create migrations directory if it doesn't exist
	migrationsDir := "migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %v", err)
	}

	// Get timestamp for migration
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Create up and down migration files
	upFile := fmt.Sprintf("%s_%s.up.sql", timestamp, name)
	downFile := fmt.Sprintf("%s_%s.down.sql", timestamp, name)

	upPath := filepath.Join(migrationsDir, upFile)
	downPath := filepath.Join(migrationsDir, downFile)

	// Create up migration file
	upTemplate := GenerateUpMigrationTemplate(name)
	if err := os.WriteFile(upPath, []byte(upTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %v", err)
	}

	// Create down migration file
	downTemplate := GenerateDownMigrationTemplate(name)
	if err := os.WriteFile(downPath, []byte(downTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %v", err)
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upPath)
	fmt.Printf("  %s\n", downPath)

	return nil
}

func GenerateUpMigrationTemplate(name string) string {
	words := strings.ReplaceAll(name, "_", " ")
	words = strings.Title(words)

	return fmt.Sprintf(`-- Migration: %s (UP)
-- Created: %s
-- Description: Add your migration description here

-- Example templates:

-- Create table:
-- CREATE TABLE example_table (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     email VARCHAR(255) UNIQUE,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- Add column:
-- ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- Create index:
-- CREATE INDEX idx_users_email ON users(email);

-- Your migration SQL here:

`, words, time.Now().Format("2006-01-02 15:04:05"))
}

func GenerateDownMigrationTemplate(name string) string {
	words := strings.ReplaceAll(name, "_", " ")
	words = strings.Title(words)

	return fmt.Sprintf(`-- Migration: %s (DOWN)
-- Created: %s
-- Description: Rollback for the up migration

-- Example rollback templates:

-- Drop table:
-- DROP TABLE IF EXISTS example_table;

-- Drop column:
-- ALTER TABLE users DROP COLUMN IF EXISTS phone;

-- Drop index:
-- DROP INDEX IF EXISTS idx_users_email;

-- Your rollback SQL here:

`, words, time.Now().Format("2006-01-02 15:04:05"))
}

// MigrateTo migrates to a specific version
func MigrateTo(version uint) error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	logger.Infof("Migrating to version %d...", version)
	err = m.Migrate(version)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration to version %d failed: %v", version, err)
	}

	if err == migrate.ErrNoChange {
		logger.Infof("Already at version %d", version)
	} else {
		logger.Infof("Successfully migrated to version %d", version)
	}

	return nil
}

// RollbackMigration rolls back one migration
func RollbackMigration() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	logger.Infof("Rolling back one migration...")
	err = m.Steps(-1) // Roll back one migration
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rollback failed: %v", err)
	}

	if err == migrate.ErrNoChange {
		logger.Infof("No migrations to rollback")
	} else {
		logger.Infof("Rollback completed successfully")
	}

	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion() (uint, bool, error) {
	sqlDB, err := DB.DB()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	return version, dirty, err
}
