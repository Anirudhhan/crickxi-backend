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

	if createMatchReq.StrikerID == createMatchReq.CurrentBowlerID || createMatchReq.NonStrikerID == createMatchReq.CurrentBowlerID {
		utils.ErrorResponse(ctx, http.StatusBadRequest,
			errors.New("bowler cannot be striker or non striker same time"),
			"bowler cannot be striker or non striker same time")
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
		var tossLostTeamID string
		var tossWinnerPlayers []models.Player
		var tossLostPlayers []models.Player
		if createMatchReq.TossWinner == "A" {
			tossWinnerTeamID = teamAID
			tossWinnerPlayers = createMatchReq.TeamAPlayers
			tossLostTeamID = teamBID
			tossLostPlayers = createMatchReq.TeamBPlayers
		} else {
			tossWinnerTeamID = teamBID
			tossWinnerPlayers = createMatchReq.TeamBPlayers
			tossLostTeamID = teamAID
			tossLostPlayers = createMatchReq.TeamAPlayers

		}

		matchID, err := dbHelper.CreateMatch(tx, createMatchReq, hostID, tossWinnerTeamID, teamAID, teamBID)
		if err != nil {
			return err
		}
		matchData.MatchID = matchID

		var battingTeamID string
		var bowlingTeamID string
		var battingPlayers []models.Player
		var bowlingPlayers []models.Player
		if createMatchReq.TossDecision == "bat" {
			battingTeamID = tossWinnerTeamID
			battingPlayers = tossWinnerPlayers
			bowlingTeamID = tossLostTeamID
			bowlingPlayers = tossLostPlayers
		} else {
			battingTeamID = tossLostTeamID
			battingPlayers = tossLostPlayers
			bowlingTeamID = tossWinnerTeamID
			bowlingPlayers = tossWinnerPlayers
		}

		inningID, err := dbHelper.StartInning(tx, matchData.MatchID, battingTeamID, bowlingTeamID, 1, "normal")
		if err != nil {
			return err
		}
		matchData.CurrentInningID = inningID

		err = dbHelper.CreateBattingScorecards(tx, matchData.CurrentInningID, battingPlayers)
		if err != nil {
			return err
		}

		err = dbHelper.CreateBowlingScorecards(tx, inningID, bowlingPlayers)
		if err != nil {
			return err
		}

		err = dbHelper.StartLiveMatch(tx, matchData.MatchID, matchData.CurrentInningID, createMatchReq)
		if err != nil {
			return err
		}

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

func GetMatchByID(ctx *gin.Context) {
	matchID := ctx.Param("matchID")
	if matchID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("match id is required"), "match id is required")
		return
	}

	matchCard, err := dbHelper.GetMatchByID(matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusNotFound, err, "invalid player id")
			return
		}

		utils.ErrorResponse(ctx, http.StatusNotFound, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, matchCard)
}
