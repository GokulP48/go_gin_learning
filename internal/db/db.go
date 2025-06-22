package db

import (
	"fmt"

	"github.com/GokulP48/go_gin_learning/config"
	"github.com/GokulP48/go_gin_learning/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDBConnection() {

	dbConfig := config.AppConfig.DB

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("InitDBConnection: failed to connect to database: " + err.Error())
	}
	logger.Infof("InitDBConnection: Initialize the DB connection successfully")
}

func DBHealthCheck() map[string]string {
	sqlDB, err := DB.DB()
	if err != nil {
		return map[string]string{"status": "error", "message": "failed to get DB instance"}
	}

	err = sqlDB.Ping()
	if err != nil {
		return map[string]string{"status": "down", "message": err.Error()}
	}

	return map[string]string{"status": "up", "message": "DB connection is healthy"}

}
