package main

import (
	"os"

	"github.com/GokulP48/go_gin_learning/config"
	"github.com/GokulP48/go_gin_learning/internal/db"
	"github.com/GokulP48/go_gin_learning/internal/logger"
	"github.com/GokulP48/go_gin_learning/internal/router"
)

func main() {

	// Load Config
	config.LoadConfig("config/config.yaml")

	// Initialize logger
	logger.InitLogger(logger.NewZapLogger(os.Stdout, config.AppConfig.Logger.Level))

	// Initialize DB connection
	db.InitDBConnection()

	// Run migrations automatically on startup
	logger.Infof("Running database migrations on startup...")
	if err := db.RunMigrations(); err != nil {
		logger.Fatalf("Migration failed: %v", err)
	}

	// Start Server
	router.InitRouter()
}
