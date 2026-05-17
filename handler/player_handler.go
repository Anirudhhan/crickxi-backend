package handler

import (
	"crickxi-backend/database"
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/models"
	"crickxi-backend/utils"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

func UpdateProfile(ctx *gin.Context) {
	var req models.UpdateProfileRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	userID := ctx.GetString("user_id")
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("missing user id"), "unauthorized")
		return
	}

	err := database.Tx(func(tx *sqlx.Tx) error {
		return dbHelper.UpdateUserProfile(tx, userID, req.Name, req.BattingStyle, req.BowlingStyle)
	})

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to update profile")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}
