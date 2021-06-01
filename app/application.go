package app

import (
	"github.com/gin-gonic/gin"
	"github.com/southern-martin/util-go/logger/logger"

	//"github.com/southern-martin/util-go/logger"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()

	logger.Info("about to start the application...")
	router.Run(":8082")
}
