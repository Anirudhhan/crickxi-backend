package routes

import (
	"crickxi-backend/handler"
	"crickxi-backend/middleware"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes() *gin.Engine {
	router := gin.Default()
	var ipOrigins = os.Getenv("IP_ORIGIN")

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			ipOrigins,
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	v1 := router.Group("/v1")

	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "server is running",
			})
		})

		v1.POST("/register", handler.RegisterUser)
		v1.POST("/login", handler.LoginUser)
		v1.POST("/reset-password", handler.ResetPassword)
		v1.GET("/refresh", handler.RefreshToken)

		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())
		auth.PUT("/logout", handler.Logout)

		{
			profile := auth.Group("/players")
			profile.GET("/:userID", handler.GetPlayerStats)
		}
	}

	return router
}
