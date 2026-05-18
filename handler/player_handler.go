package handler

import (
	"crickxi-backend/database"
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/models"
	"crickxi-backend/utils"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func SearchPlayer(ctx *gin.Context) {
	search := strings.TrimSpace(
		ctx.Query("q"),
	)

	if search == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("missing search query"), "search query required")
		return
	}

	players, err := dbHelper.SearchPlayers(search)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to search players")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"players": players})
}

func CreateGuestPlayer(ctx *gin.Context) {
	var req models.CreateGuestPlayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	var userID string
	var playerID string

	err := database.Tx(func(tx *sqlx.Tx) error {
		var txErr error
		userID, playerID, txErr =
			dbHelper.CreateGuestPlayer(tx, req.Name, req.Phone)
		return txErr
	})

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == DuplicatePGCode {
				utils.ErrorResponse(ctx, http.StatusConflict, err, "player already exists")
				return
			}
		}

		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to create player")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "player created successfully",
		"userID":   userID,
		"playerID": playerID,
	})
}
