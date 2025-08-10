package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GokulP48/go_gin_learning/config"
	"github.com/GokulP48/go_gin_learning/internal/db"
	"github.com/GokulP48/go_gin_learning/internal/logger"
	"github.com/GokulP48/go_gin_learning/internal/router"
)

func main() {

	var (
		create    = flag.String("create", "", "Create a new migration file")
		up        = flag.Bool("up", false, "Run all up migrations")
		down      = flag.Bool("down", false, "Roll back one migration")
		version   = flag.Bool("version", false, "Show current migration version")
		migrateTo = flag.Uint("migrate-to", 0, "Migrate to specific version")
	)
	flag.Parse()

	// Handle create migration (doesn't need DB connection)
	if *create != "" {
		if err := db.CreateMigration(*create); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		return
	}

	// Load Config
	config.LoadConfig("config/config.yaml")

	// Initialize logger
	logger.InitLogger(logger.NewZapLogger(os.Stdout, config.AppConfig.Logger.Level))

	// Initialize DB connection
	db.InitDBConnection()

	// Run migrations automatically on startup
	logger.Infof("Running database migrations on startup...")
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	switch {
	case *up:
		if err := db.RunMigrations(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

	case *down:
		if err := db.RollbackMigration(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}

	case *version:
		version, dirty, err := db.GetMigrationVersion()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}

		status := "clean"
		if dirty {
			status = "dirty"
		}
		fmt.Printf("Current version: %d (status: %s)\n", version, status)

	case *migrateTo > 0:
		if err := db.MigrateTo(*migrateTo); err != nil {
			log.Fatalf("Migration to version %d failed: %v", *migrateTo, err)
		}

	default:
		flag.Usage()
		fmt.Println("\nExamples:")
		fmt.Println("  go run cmd/migrate/main.go -create create_users_table")
		fmt.Println("  go run cmd/migrate/main.go -up")
		fmt.Println("  go run cmd/migrate/main.go -down")
		fmt.Println("  go run cmd/migrate/main.go -version")
		fmt.Println("  go run cmd/migrate/main.go -migrate-to 5")
	}

	// Start Server
	router.InitRouter()
}
