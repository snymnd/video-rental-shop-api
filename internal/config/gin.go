package config

import "github.com/gin-gonic/gin"

func NewGin() *gin.Engine {
	router := gin.New()
	router.ContextWithFallback = true

	return router
}