package handler

import (
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/models"
	"crickxi-backend/utils"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetScorecardByMatchIDAndInning(ctx *gin.Context) {
	matchID := ctx.Param("matchID")
	inningOrderStr := ctx.Param("inningOrder")

	inningOrder, err := strconv.Atoi(inningOrderStr)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest,
			errors.New("inning order must be a number"), "inning order must be a number")
		return
	}

	if matchID == "" || (inningOrder != 1 && inningOrder != 2) {
		utils.ErrorResponse(ctx, http.StatusBadRequest,
			errors.New("valid match ID and inning order are required"), "valid match ID and inning order are required")
		return
	}

	var matchScoreCard models.MatchScoreCard
	matchScoreCard.InningOrder = inningOrder

	inningDetails, err := dbHelper.GetInningDetails(matchID, matchScoreCard.InningOrder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusNotFound, err, "invalid match id")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}
	matchScoreCard.BattingTeamID = inningDetails.BattingTeamID
	matchScoreCard.BattingTeamName = inningDetails.BattingTeamName

	matchScoreCard.BowlingTeamID = inningDetails.BowlingTeamID
	matchScoreCard.BowlingTeamName = inningDetails.BowlingTeamName

	battingScoreCard, err := dbHelper.GetBattingScorecardByMatchIDAndInning(matchID, matchScoreCard.InningOrder)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	bowlingScoreCard, err := dbHelper.GetBowlingScorecardByMatchIDAndInning(matchID, matchScoreCard.InningOrder)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	matchScoreCard.BattingScoreCard = battingScoreCard
	matchScoreCard.BowlingScoreCard = bowlingScoreCard

	ctx.JSON(http.StatusOK, matchScoreCard)
}
