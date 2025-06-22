package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GokulP48/go_gin_learning/config"
	"github.com/GokulP48/go_gin_learning/internal/logger"
	"github.com/GokulP48/go_gin_learning/internal/router"
)

func main() {

	// Load Config
	config.LoadConfig("config/config.yaml")

	// Initialize logger
	logger.InitLogger(logger.NewZapLogger(os.Stdout, config.LogLevel()))

	// Initialize DB connection
	// db.InitDBConnection()

	r := &router.Router{}
	handler := r.InitRouter()

	port := config.ServerPort()
	server := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	log.Printf("🚀 Server is running at %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
