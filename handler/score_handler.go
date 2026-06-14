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

func GetScorecardByMatchIDAndInnings(ctx *gin.Context) {
	matchID := ctx.Param("matchID")
	inningsOrderStr := ctx.Param("inningsOrder")

	inningsOrder, err := strconv.Atoi(inningsOrderStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest,
			errors.New("innings order must be a number"), "innings order must be a number")
		return
	}

	if matchID == "" || inningsOrder <= 0 {
		utils.ErrorResponse(ctx, http.StatusBadRequest,
			errors.New("valid match ID and innings order are required"), "valid match ID and innings order are required")
		return
	}
	
	var matchScoreCard models.MatchScoreCard
	matchScoreCard.InningsOrder = inningsOrder

	inningsDetails, err := dbHelper.GetInningsDetails(matchID, matchScoreCard.InningsOrder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusNotFound, err, "invalid match id")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}
	matchScoreCard.BattingTeamID = inningsDetails.BattingTeamID
	matchScoreCard.BattingTeamName = inningsDetails.BattingTeamName

	matchScoreCard.BowlingTeamID = inningsDetails.BowlingTeamID
	matchScoreCard.BowlingTeamName = inningsDetails.BowlingTeamName

	matchScoreCard.TotalRuns = inningsDetails.TotalRuns
	matchScoreCard.Extras = inningsDetails.Extras

	battingScoreCard, err := dbHelper.GetBattingScorecardByMatchIDAndInnings(nil, matchID, matchScoreCard.InningsOrder)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	bowlingScoreCard, err := dbHelper.GetBowlingScorecardByMatchIDAndInnings(nil, matchID, matchScoreCard.InningsOrder)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	matchScoreCard.BattingScoreCard = battingScoreCard
	matchScoreCard.BowlingScoreCard = bowlingScoreCard

	ctx.JSON(http.StatusOK, matchScoreCard)
}
