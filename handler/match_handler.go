package handler

import (
	"crickxi-backend/database"
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/models"
	"crickxi-backend/utils"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

	if createMatchReq.Overs <= 0 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("over should be greater than 0"), "over should be greater than 0")
		return
	}

	teamA := strings.ToLower(strings.TrimSpace(createMatchReq.TeamAName))
	teamB := strings.ToLower(strings.TrimSpace(createMatchReq.TeamBName))

	if teamA == teamB {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("both teams can't have the same name"), "both teams can't have the same name")
		return
	}

	err := ValidateBattersHelper(createMatchReq.StrikerID, createMatchReq.NonStrikerID, createMatchReq.CurrentBowlerID)
	if err != nil {
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

		if !IsPlayerInTeam(battingPlayers, createMatchReq.StrikerID) {
			return errors.New("striker must be in the batting team")
		}

		if createMatchReq.NonStrikerID != nil && !IsPlayerInTeam(battingPlayers, *createMatchReq.NonStrikerID) {
			return errors.New("non-striker must be in the batting team")
		}

		if !IsPlayerInTeam(bowlingPlayers, createMatchReq.CurrentBowlerID) {
			return errors.New("bowler must be in the bowling team")
		}

		inningID, err := dbHelper.StartInnings(tx, matchData.MatchID, battingTeamID, bowlingTeamID, 1, "normal")
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
	search := ctx.DefaultQuery("search", "")
	status := ctx.DefaultQuery("status", "")

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

	matches, err := dbHelper.GetMatches(search, status, page, limit)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to get matches")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"status":  status,
		"page":    page,
		"total":   len(matches),
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

		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, matchCard)
}

func StartNextInnings(ctx *gin.Context) {
	var nextInningsReq models.StartNextInningsReq
	matchID := ctx.GetString("match_id")

	if err := ctx.ShouldBindJSON(&nextInningsReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	err := ValidateBattersHelper(nextInningsReq.StrikerID, nextInningsReq.NonStrikerID, nextInningsReq.BowlerID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	matchData, err := dbHelper.GetLiveMatchDetails(matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusNotFound, err, "invalid match id")
			return
		}

		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	if matchData.CurrentInningNo != 1 {
		utils.ErrorResponse(ctx, http.StatusConflict, errors.New("second innings already started"), "second innings already started")
		return
	}

	if matchData.EndTime != nil {
		utils.ErrorResponse(ctx, http.StatusConflict, errors.New("match already completed"), "match already completed")
		return
	}

	if !matchData.IsCompleted {
		utils.ErrorResponse(ctx, http.StatusConflict, errors.New("first innings not completed"), "first innings not completed")
		return
	}
	//Transactions
	{
		txErr := database.Tx(func(tx *sqlx.Tx) error {

			// swap teams
			battingTeamID := matchData.BowlingTeamID

			bowlingTeamID := matchData.BattingTeamID

			inningID, err := dbHelper.StartInnings(tx, matchID, battingTeamID, bowlingTeamID, 2, "normal")

			if err != nil {
				return err
			}

			battingPlayers, err := dbHelper.GetPlayersByTeamID(battingTeamID)

			if err != nil {
				return err
			}

			bowlingPlayers, err := dbHelper.GetPlayersByTeamID(bowlingTeamID)

			if err != nil {
				return err
			}

			if !IsPlayerInTeam(battingPlayers, nextInningsReq.StrikerID) {
				return errors.New("striker must be in the batting team")
			}

			if nextInningsReq.NonStrikerID != nil && !IsPlayerInTeam(battingPlayers, *nextInningsReq.NonStrikerID) {
				return errors.New("non-striker must be in the batting team")
			}

			if !IsPlayerInTeam(bowlingPlayers, nextInningsReq.BowlerID) {
				return errors.New("bowler must be in the bowling team")
			}

			err = dbHelper.CreateBattingScorecards(tx, inningID, battingPlayers)

			if err != nil {
				return err
			}

			err = dbHelper.CreateBowlingScorecards(tx, inningID, bowlingPlayers)

			if err != nil {
				return err
			}

			err = dbHelper.ResetLiveMatchForNextInnings(tx, matchID, inningID, nextInningsReq)

			if err != nil {
				return err
			}

			err = dbHelper.UpdateMatchInningNo(tx, matchID, 2)

			if err != nil {
				return err
			}

			return nil
		})

		if txErr != nil {

			utils.ErrorResponse(ctx, http.StatusInternalServerError, txErr, "failed to start second innings")
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "second innings started",
	})

}

func IsPlayerInTeam(players []models.Player, playerID string) bool {
	for _, p := range players {
		if p.PlayerID == playerID {
			return true
		}
	}
	return false
}

func ValidateBattersHelper(strikerID string, nonStrikerID *string, bowlerID string) error {
	if strikerID == "" {
		return errors.New("striker is required")
	}

	if strikerID == bowlerID {
		return errors.New("bowler cannot be striker")
	}
	if nonStrikerID != nil {
		if strikerID == *nonStrikerID {
			return errors.New("striker and non striker cannot be same")
		}

		if *nonStrikerID == bowlerID {
			return errors.New("bowler cannot be non striker")
		}
	}
	return nil
}
