package app

import (
	"github.com/southern-martin/user-api/controller/ping"
	"github.com/southern-martin/user-api/controller/user"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/user", user.Create)
	router.GET("/user/:user_id", user.Get)
	router.PUT("/user/:user_id", user.Update)
	router.PATCH("/user/:user_id", user.Update)
	router.DELETE("/user/:user_id", user.Delete)
	router.GET("/internal/user/search", user.Search)
	router.POST("/user/login", user.Login)
}
