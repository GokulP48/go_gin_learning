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
	logger.Infof(".............: %s", config.AppConfig.DB.Host)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DB.Host, config.AppConfig.DB.User, config.AppConfig.DB.Password, config.AppConfig.DB.Name, config.AppConfig.DB.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("InitDBConnection: failed to connect to database: " + err.Error())
	}
	logger.Infof("InitDBConnection: Initialize the DB connection successfully")
}
