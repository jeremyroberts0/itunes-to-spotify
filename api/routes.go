package api

import (
	"github.com/gin-gonic/gin"
)

func createError(message string, err error) map[string]string {
	return map[string]string{
		"message": message,
		"error":   err.Error(),
	}
}

func GetRouter() *gin.Engine {
	router := gin.Default()

	Auth(router)
	MatchItunesPlaylistToSpotify(router)
	ApplyStatsRoutes(router)

	return router
}
