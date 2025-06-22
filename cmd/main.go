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

	// Start Server
	router.InitRouter()
}
