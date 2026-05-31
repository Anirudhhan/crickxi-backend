package handler

import (
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBattingLeaderboard(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	leaderboard, err := dbHelper.GetBattingLeaderboard(page, limit)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, leaderboard)
}

func GetBowlingLeaderboard(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	leaderboard, err := dbHelper.GetBowlingLeaderboard(page, limit)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to get bowling leaderboard")
		return
	}

	ctx.JSON(http.StatusOK, leaderboard)
}

func GetFieldingLeaderboard(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	leaderboard, err := dbHelper.GetFieldingLeaderboard(page, limit)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to get fielding leaderboard")
		return
	}

	ctx.JSON(http.StatusOK, leaderboard)
}
