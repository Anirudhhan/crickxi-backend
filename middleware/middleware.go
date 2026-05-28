package middleware

import (
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const bearerPrefix = "Bearer "

func AuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("missing authorization header"), "unauthorized")
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, bearerPrefix) {

			utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("invalid authorization format"), "unauthorized")
			ctx.Abort()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := utils.ValidateAccessToken(accessToken)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err, "invalid token")
			ctx.Abort()
			return
		}

		sessionID, ok := claims["sid"].(string)
		if !ok {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("invalid session id"), "invalid token")
			ctx.Abort()
			return
		}

		claimUserID, ok := claims["uid"].(string)
		if !ok {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("invalid user id"), "invalid token")
			ctx.Abort()
			return
		}

		sessionUserDetails, err := dbHelper.GetUserAndPlayerIDByActiveSession(sessionID)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err, "invalid session")
			ctx.Abort()
			return
		}

		if sessionUserDetails.UserID != claimUserID {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("session mismatch"), "invalid session")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", sessionUserDetails.UserID)
		ctx.Set("player_id", sessionUserDetails.PlayerID)
		ctx.Set("session_id", sessionID)

		ctx.Next()
	}
}

func HostMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userID := ctx.GetString("user_id")
		matchID := ctx.Param("matchID")

		if matchID == "" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("missing match id"), "missing match id")
			ctx.Abort()
			return
		}

		isValid, err := dbHelper.ValidateHostOrScorer(matchID, userID)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
			ctx.Abort()
			return
		}

		if !isValid {
			utils.ErrorResponse(ctx, http.StatusForbidden, errors.New("invalid access"), "invalid access")
			ctx.Abort()
			return
		}

		ctx.Set("match_id", matchID)

		ctx.Next()
	}
}
