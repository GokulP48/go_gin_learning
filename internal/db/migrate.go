package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GokulP48/go_gin_learning/internal/logger"
	"gorm.io/gorm"
)

// MigrationRecord represents a migration entry in the database
type MigrationRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Migration string `gorm:"unique;not null"`
	Batch     int    `gorm:"not null"`
}

// RunMigrations executes all pending migrations from the migrations folder
func RunMigrations() error {
	// Create migrations table if it doesn't exist
	if err := DB.AutoMigrate(&MigrationRecord{}); err != nil {
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
	var executedMigrations []MigrationRecord
	DB.Find(&executedMigrations)

	executedMap := make(map[string]bool)
	for _, migration := range executedMigrations {
		executedMap[migration.Migration] = true
	}

	// Get next batch number
	var lastBatch int
	DB.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch)
	nextBatch := lastBatch + 1

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
			migrationRecord := MigrationRecord{
				Migration: migrationFile,
				Batch:     nextBatch,
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

// GetMigrationStatus returns the current migration status
func GetMigrationStatus() ([]MigrationRecord, error) {
	var migrations []MigrationRecord
	err := DB.Order("batch ASC, id ASC").Find(&migrations).Error
	return migrations, err
}
