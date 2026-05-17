package handler

import (
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/utils"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayerProfile(ctx *gin.Context) {
	userID := ctx.Param("playerStatsID")

	playerStats, err := dbHelper.GetPlayerProfileByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err, "invalid player id")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, playerStats)
}
