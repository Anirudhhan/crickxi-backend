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
	var ipOrigins2 = os.Getenv("IP_ORIGIN2")

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			ipOrigins,
			ipOrigins2,
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	v1 := router.Group("/v1")

	// public routes
	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "server is running",
			})
		})

		v1.POST("/register", handler.RegisterUser)
		v1.POST("/login", handler.LoginUser)
		v1.POST("/request-password-reset", handler.RequestPasswordReset)
		v1.POST("/reset-password", handler.ResetPassword)
		v1.GET("/refresh", handler.RefreshToken)

		v1.GET("/profile/:playerStatsID", handler.GetPlayerProfile)

		v1.GET("/matches", handler.GetMatches)
		v1.GET("/match/:matchID", handler.GetMatchByID)
		v1.GET("/match/:matchID/scorecard/:inningOrder", handler.GetScorecardByMatchIDAndInning)
	}

	// protected routes
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware())

	{
		auth.PUT("/logout", handler.Logout)

		auth.GET("/profile/me", handler.GetMyProfile)

		auth.PUT("/profile/me", handler.UpdateProfile)

		auth.POST("/player", handler.CreateGuestPlayer)

		auth.GET("/players/search", handler.SearchPlayer)

		auth.POST("/match", handler.CreateMatch)

		auth.POST("/matches/:matchID/balls", handler.BallEvent)

		auth.POST("/matches/:matchID/change-bowler", handler.ChangeBowler)

		auth.POST("/matches/:matchID/start-next-innings", handler.StartNextInnings)
	}

	return router
}
