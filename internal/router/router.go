package router

import (
	"fmt"

	"github.com/GokulP48/go_gin_learning/config"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()

	router.Run(fmt.Sprintf(":%s", config.AppConfig.Server.Port))
}
