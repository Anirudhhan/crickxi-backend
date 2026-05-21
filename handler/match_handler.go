package handler

import (
	"crickxi-backend/database"
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/models"
	"crickxi-backend/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CreateMatch(ctx *gin.Context) {
	var createMatchReq models.CreateMatchRequest
	hostID := ctx.GetString("user_id")

	if err := ctx.ShouldBindJSON(&createMatchReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	if hostID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("missing player id"), "unauthorized")
		return
	}

	var matchData models.MatchData
	err := database.Tx(func(tx *sqlx.Tx) error {
		var txErr error
		matchData, txErr = dbHelper.CreateMatch(tx, createMatchReq, hostID)
		return txErr
	})

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to create match")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "match created successfully",
		"data":    matchData})
}

func GetMatches(ctx *gin.Context) {
	matches, err := dbHelper.GetMatches()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to get matches")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"matches": matches,
	})
}
