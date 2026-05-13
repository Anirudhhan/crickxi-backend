package routes

import (
	"crickxi-backend/handler"
	"crickxi-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")

	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "server is running",
			})
		})
		v1.POST("/register", handler.RegisterUser)
		v1.POST("/login", handler.LoginUser)
		v1.GET("/refresh", handler.RefreshToken)

		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())
		auth.PUT("/logout", handler.Logout)
	}

	return router
}
