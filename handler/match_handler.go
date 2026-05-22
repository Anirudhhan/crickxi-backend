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
	if hostID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("missing user id"), "unauthorized")
		return
	}

	if err := ctx.ShouldBindJSON(&createMatchReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	var matchData models.MatchData

	txErr := database.Tx(func(tx *sqlx.Tx) error {

		teamAID, err := dbHelper.CreateTeam(tx, createMatchReq.TeamAName, hostID)
		if err != nil {
			return err
		}

		matchData.TeamAID = teamAID

		teamBID, err := dbHelper.CreateTeam(tx, createMatchReq.TeamBName, hostID)
		if err != nil {
			return err
		}

		matchData.TeamBID = teamBID

		err = dbHelper.AddPlayersToTeam(tx, teamAID, createMatchReq.TeamAPlayers)
		if err != nil {
			return err
		}

		err = dbHelper.AddPlayersToTeam(tx, teamBID, createMatchReq.TeamBPlayers)
		if err != nil {
			return err
		}

		var tossWinnerTeamID string
		if createMatchReq.TossWinner == "A" {
			tossWinnerTeamID = teamAID
		} else {
			tossWinnerTeamID = teamBID
		}

		matchID, err := dbHelper.CreateMatch(tx, createMatchReq, hostID, tossWinnerTeamID, teamAID, teamBID)
		if err != nil {
			return err
		}

		matchData.MatchID = matchID
		return nil
	})

	if txErr != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, txErr, "failed to create match")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "match created successfully",
		"data":    matchData,
	})
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
