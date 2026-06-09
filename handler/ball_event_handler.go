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

func BallEvent(ctx *gin.Context) {
	var req models.BallEventReq

	matchID := ctx.GetString("match_id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	liveMatchData, err := dbHelper.GetLiveMatchDetails(matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusNotFound, err, "invalid match id")
			return
		}

		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	if liveMatchData.EndTime != nil {
		utils.ErrorResponse(ctx, http.StatusConflict, errors.New("match already completed"), "match already completed")
		return
	}

	if liveMatchData.IsCompleted {
		utils.ErrorResponse(ctx, http.StatusConflict, errors.New("innings ended. start next innings"), "innings ended. start next innings")
		return
	}

	var delivery models.Delivery

	PrepareDeliveryHelper(&delivery, liveMatchData, req)
	if delivery.IsWicket {
		if delivery.WicketPlayerID == nil || delivery.WicketType == nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid wicket data"), "invalid wicket data")
			return
		}
	}

	// striker, non-striker and next batter not out validation
	if err = ValidateBattersNotOutHelper(delivery); err != nil {
		utils.ErrorResponse(ctx, http.StatusConflict, err, err.Error())
		return
	}

	matchEnded := false
	inningEnded := false
	txErr := ProcessBallEventHelper(delivery, liveMatchData, matchID, &inningEnded, &matchEnded)

	if txErr != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, txErr, "failed to add ball")
		return
	}

	if inningEnded {
		ctx.JSON(http.StatusOK, gin.H{"message": "innings ended. start next innings"})
		return
	}
	if matchEnded {
		ctx.JSON(http.StatusOK, gin.H{"message": "match ended"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "ball added successfully",
	})
}

func ValidateBattersNotOutHelper(delivery models.Delivery) error {
	if delivery.NextBatterID != nil {
		if *delivery.NextBatterID == delivery.BowlerID {
			return errors.New("next batter cannot be the current bowler")
		}

		isNextBatterOut, err := dbHelper.IsPlayerOut(delivery.InningsID, *delivery.NextBatterID)
		if err != nil {
			return err
		}
		if isNextBatterOut {
			return errors.New("next batter is already out")
		}
	}

	isStrikerOut, err := dbHelper.IsPlayerOut(delivery.InningsID, delivery.StrikerID)
	if err != nil {
		return err
	}
	if isStrikerOut {
		return errors.New("striker is already out")
	}

	if delivery.NonStrikerID != nil {
		isNonStrikerOut, err := dbHelper.IsPlayerOut(delivery.InningsID, *delivery.NonStrikerID)
		if err != nil {
			return err
		}
		if isNonStrikerOut {
			return errors.New("non-striker is already out")
		}
	}

	return nil
}

func PrepareDeliveryHelper(delivery *models.Delivery, liveMatchData models.LiveMatchDetails, req models.BallEventReq) {
	delivery.InningsID = liveMatchData.CurrentInningID
	delivery.StrikerID = liveMatchData.StrikerID
	delivery.NonStrikerID = liveMatchData.NonStrikerID
	delivery.BowlerID = liveMatchData.CurrentBowlerID
	delivery.LegalBalls = liveMatchData.LegalBalls

	delivery.BallSequence = liveMatchData.CurrentBallSequence + 1
	delivery.OverNumber = liveMatchData.LegalBalls / 6
	delivery.BallInOver = (liveMatchData.LegalBalls % 6) + 1
	delivery.IsFreeHit = liveMatchData.IsFreeHit

	delivery.IsLegalDelivery = true
	if req.ExtraType != nil {
		delivery.ExtraType = req.ExtraType
		switch *req.ExtraType {
		case "wide", "no_ball":
			delivery.BallInOver = liveMatchData.LegalBalls % 6
			delivery.IsLegalDelivery = false
		}
	}

	delivery.RunsBatter = req.Runs
	delivery.RunsExtra = req.ExtraRuns

	if req.IsWicket != nil && *req.IsWicket {
		delivery.IsWicket = true
		delivery.WicketType = req.WicketType
		delivery.WicketPlayerID = req.WicketPlayerID
		delivery.FielderID = req.FielderID
		delivery.NextBatterID = req.NextBatterID

		if req.WicketType != nil && (*req.WicketType == "retired_hurt" || *req.WicketType == "retired_out") {
			delivery.IsLegalDelivery = false
			delivery.BallInOver = liveMatchData.LegalBalls % 6
		}
	}

}

func ProcessBallEventHelper(delivery models.Delivery, liveMatchData models.LiveMatchDetails, matchID string, inningEnded *bool, matchEnded *bool) error {
	err := database.Tx(func(tx *sqlx.Tx) error {

		err := dbHelper.CreateBallEvent(tx, delivery)
		if err != nil {
			return err
		}

		return ApplyBallStats(tx, delivery, liveMatchData, matchID, inningEnded, matchEnded)
	})
	if err != nil {
		return err
	}
	return nil
}

func ApplyBallStats(tx *sqlx.Tx, delivery models.Delivery, liveMatchData models.LiveMatchDetails, matchID string, inningEnded *bool, matchEnded *bool) error {
	legalBall := 0
	if delivery.IsLegalDelivery {
		legalBall = 1
	}

	fours := 0
	if delivery.RunsBatter == 4 {
		fours = 1
	}

	sixes := 0
	if delivery.RunsBatter == 6 {
		sixes = 1
	}

	totalRuns := delivery.RunsBatter + delivery.RunsExtra
	wides := 0
	noBall := 0
	extraWide := 0
	extraNoBall := 0
	nextFreeHit := false

	if !delivery.IsLegalDelivery {
		if delivery.ExtraType != nil {
			if *delivery.ExtraType == "wide" {
				extraWide = delivery.RunsExtra
				wides += 1
				if delivery.IsFreeHit {
					nextFreeHit = true
				}
			}
			if *delivery.ExtraType == "no_ball" {
				extraNoBall = delivery.RunsExtra
				noBall += 1
				nextFreeHit = true
			}
		}
	}

	//update batting scoreCard
	err := dbHelper.UpdateBatterStats(tx, delivery, legalBall, fours, sixes)
	if err != nil {
		return err
	}

	//Update dismissed player
	wicket := 0
	inningsWicket := 0
	if delivery.IsWicket && delivery.WicketPlayerID != nil {
		var dismissalBy *string
		isOut := true
		if delivery.WicketType != nil {
			switch *delivery.WicketType {
			case "bowled", "caught", "lbw", "stumped", "hit_wicket":
				wicket = 1
				inningsWicket = 1
				dismissalBy = &delivery.BowlerID
			case "run_out":
				inningsWicket = 1
				if delivery.FielderID != nil {
					dismissalBy = delivery.FielderID
				}
			case "retired_hurt":
				wicket = 0
				inningsWicket = 0
				isOut = false
			case "retired_out":
				wicket = 0
				inningsWicket = 1
			}
		}

		err = dbHelper.UpdateDismissedBatter(tx, delivery, dismissalBy, isOut)
		if err != nil {
			return err
		}
	}

	//update bowling scorecard
	err = dbHelper.UpdateBowlingScorecard(tx, delivery, legalBall, totalRuns, wides, noBall, wicket)
	if err != nil {
		return err
	}

	//Update innings
	err = dbHelper.UpdateInnings(tx, delivery, totalRuns, inningsWicket, legalBall, extraWide, extraNoBall)
	if err != nil {
		return err
	}

	//update live match
	{
		// original positions before ball
		originalStrikerID := delivery.StrikerID
		originalNonStrikerID := delivery.NonStrikerID

		// current positions after ball movement
		strikerID := delivery.StrikerID
		nonStrikerID := delivery.NonStrikerID

		// strike rotation
		if nonStrikerID != nil {
			runningExtras := delivery.RunsExtra
			if !delivery.IsLegalDelivery && delivery.ExtraType != nil {
				if *delivery.ExtraType == "wide" || *delivery.ExtraType == "no_ball" {
					runningExtras -= 1
				}
			}

			totalMovementRuns := delivery.RunsBatter + runningExtras
			if totalMovementRuns%2 == 1 {
				temp := strikerID
				strikerID = *nonStrikerID
				nonStrikerID = &temp
			}
		}

		// wicket handling
		if delivery.IsWicket && delivery.WicketPlayerID != nil {
			// original striker out
			if *delivery.WicketPlayerID == originalStrikerID {
				if delivery.NextBatterID != nil {
					// if strike rotated (now at non-striker end),
					// new batter becomes non-striker
					if strikerID != originalStrikerID {
						nonStrikerID = delivery.NextBatterID
					} else {
						// if strike didn't rotate (still at striker end),
						// new batter becomes striker
						strikerID = *delivery.NextBatterID
					}
				} else {
					// LAST WICKET/NO NEXT BATTER
					// If striker is out and no one is coming in,
					// the person who was NOT out (the non-striker) must become the striker
					if originalNonStrikerID != nil {
						strikerID = *originalNonStrikerID
						nonStrikerID = nil
					}
				}
			}

			// original non-striker out
			if originalNonStrikerID != nil && *delivery.WicketPlayerID == *originalNonStrikerID {
				if delivery.NextBatterID != nil {
					// if strike rotated (now at striker end),
					// new batter becomes striker
					if strikerID == *originalNonStrikerID {
						strikerID = *delivery.NextBatterID
					} else {
						// if strike didn't rotate (still at non-striker end),
						// new batter becomes non-striker
						nonStrikerID = delivery.NextBatterID
					}
				} else {
					// LAST WICKET/NO NEXT BATTER
					// Non-striker is out, striker stays as striker
					nonStrikerID = nil
				}
			}
		}

		// over completed
		newLegalBalls := delivery.LegalBalls + legalBall
		if legalBall == 1 && newLegalBalls%6 == 0 {
			if nonStrikerID != nil {
				temp := strikerID
				strikerID = *nonStrikerID
				nonStrikerID = &temp
			}
		}
		err = dbHelper.UpdateLiveMatch(tx, delivery, matchID, totalRuns, inningsWicket, legalBall, strikerID, nonStrikerID, nextFreeHit)
		if err != nil {
			return err
		}
	}

	// Match Completion handling
	{
		liveMatchData.CurrentScore += totalRuns
		if liveMatchData.PreviousInningsScore != nil {
			//handle score chased
			target := *liveMatchData.PreviousInningsScore
			if liveMatchData.CurrentScore > target {
				err = HandleInningsOrMatchCompletion(tx, liveMatchData, matchID, inningEnded, matchEnded)
				if err != nil {
					return err
				}

				return nil
			}
		}
		if delivery.IsWicket {
			//handle wicket ended
			currentWickets := liveMatchData.Wickets + inningsWicket
			if currentWickets >= liveMatchData.BattingPlayerCount {
				err = HandleInningsOrMatchCompletion(tx, liveMatchData, matchID, inningEnded, matchEnded)
				if err != nil {
					return err
				}

				return nil
			}
		}

		//handle over/inning ended
		currentLegalBalls := liveMatchData.LegalBalls
		if delivery.IsLegalDelivery {
			currentLegalBalls++
		}
		if currentLegalBalls >= liveMatchData.OversPerSide*6 {
			err = HandleInningsOrMatchCompletion(tx, liveMatchData, matchID, inningEnded, matchEnded)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func HandleInningsOrMatchCompletion(tx *sqlx.Tx, liveMatchData models.LiveMatchDetails, matchID string, inningEnded *bool, matchEnded *bool) error {

	err := dbHelper.CompleteInnings(tx, liveMatchData.CurrentInningID)
	if err != nil {
		return err
	}

	// second innings, match completed
	if liveMatchData.CurrentInningNo == 2 {
		var winnerTeamID *string
		if liveMatchData.PreviousInningsScore != nil {
			if liveMatchData.CurrentScore > *liveMatchData.PreviousInningsScore {
				winnerTeamID = &liveMatchData.BattingTeamID

			} else if liveMatchData.CurrentScore < *liveMatchData.PreviousInningsScore {
				winnerTeamID = &liveMatchData.BowlingTeamID
			}
		}

		err = dbHelper.CompleteMatch(tx, matchID, winnerTeamID)
		if err != nil {
			return err
		}
		err = UpdatePlayerCareerStats(tx, matchID, winnerTeamID)
		if err != nil {
			return err
		}
	}

	if liveMatchData.CurrentInningNo == 1 {
		*inningEnded = true
	} else {
		*matchEnded = true
	}

	return nil
}

func ChangeBowler(ctx *gin.Context) {
	var req struct {
		BowlerID string `json:"bowlerID"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	matchID := ctx.GetString("match_id")
	validBowler, err := dbHelper.ValidateBowlerID(matchID, req.BowlerID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	if !validBowler {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid bowler id"), "invalid bowler")
		return
	}

	err = dbHelper.ChangeBowler(matchID, req.BowlerID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "bowler changed successfully"})
}

func UpdatePlayerCareerStats(tx *sqlx.Tx, matchID string, winnerTeamID *string) error {

	for inningOrder := 1; inningOrder <= 2; inningOrder++ {

		battingScorecards, err := dbHelper.GetBattingScorecardByMatchIDAndInnings(tx, matchID, inningOrder)
		if err != nil {
			return err
		}

		matchPlayed := 0
		if inningOrder == 2 {
			matchPlayed = 1
		}

		for _, batting := range battingScorecards {

			matchesWon := 0
			matchesLost := 0

			if winnerTeamID != nil {
				if batting.TeamID == *winnerTeamID {
					matchesWon = 1
				} else {
					matchesLost = 1
				}
			}

			ducks := 0
			goldenDucks := 0

			if batting.IsOut && batting.Runs == 0 {
				ducks = 1

				if batting.Balls == 1 {
					goldenDucks = 1
				}
			}

			fifties := 0
			if batting.Runs >= 50 && batting.Runs < 100 {
				fifties = 1
			}

			hundreds := 0
			if batting.Runs >= 100 {
				hundreds = 1
			}

			inningsBatted := 0
			if batting.Balls > 0 || batting.IsOut || batting.DismissalType != nil {
				inningsBatted = 1
			}

			notOuts := 0
			if inningsBatted == 1 && !batting.IsOut {
				notOuts = 1
			}

			stats := models.UpdatePlayerStats{
				Runs:          &batting.Runs,
				BallsFaced:    &batting.Balls,
				InningsBatted: &inningsBatted,
				NotOuts:       &notOuts,
				Fours:         &batting.Fours,
				Sixes:         &batting.Sixes,

				HighestScore: &batting.Runs,

				Ducks:       &ducks,
				GoldenDucks: &goldenDucks,
				Fifties:     &fifties,
				Hundreds:    &hundreds,

				MatchesPlayed: &matchPlayed,
				MatchesWon:    &matchesWon,
				MatchesLost:   &matchesLost,
			}

			err = dbHelper.UpdatePlayerStats(
				tx,
				batting.PlayerID,
				stats,
			)
			if err != nil {
				return err
			}
		}

		bowlingScorecards, err := dbHelper.GetBowlingScorecardByMatchIDAndInnings(
			tx,
			matchID,
			inningOrder,
		)
		if err != nil {
			return err
		}

		for _, bowling := range bowlingScorecards {

			inningsBowled := 0
			if bowling.LegalBalls > 0 || bowling.Wides > 0 || bowling.NoBalls > 0 {
				inningsBowled = 1
			}

			stats := models.UpdatePlayerStats{
				Wickets:      &bowling.Wickets,
				BallsBowled:  &bowling.LegalBalls,
				RunsConceded: &bowling.RunsGiven,
				MaidenOvers:  &bowling.Maidens,
				Wides:        &bowling.Wides,
				NoBalls:      &bowling.NoBalls,

				BestWickets: &bowling.Wickets,
				BestRuns:    &bowling.RunsGiven,

				MatchesPlayed: &matchPlayed,
				InningsBowled: &inningsBowled,
			}

			err = dbHelper.UpdatePlayerStats(
				tx,
				bowling.PlayerID,
				stats,
			)
			if err != nil {
				return err
			}
		}
	}

	fieldingStats, err := dbHelper.GetFieldingStatsByMatchID(tx, matchID)
	if err != nil {
		return err
	}

	for _, fielding := range fieldingStats {

		stats := models.UpdatePlayerStats{
			Catches:   &fielding.Catches,
			RunOuts:   &fielding.RunOuts,
			Stumpings: &fielding.Stumpings,
		}

		err = dbHelper.UpdatePlayerStats(
			tx,
			fielding.PlayerID,
			stats,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func UndoBall(ctx *gin.Context) {
	matchID := ctx.GetString("match_id")

	liveMatchData, err := dbHelper.GetLiveMatchDetails(matchID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to get match details")
		return
	}

	// if match is already completed, do not allow undo
	matchCard, err := dbHelper.GetMatchByID(matchID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to get match card")
		return
	}
	if matchCard.MatchStatus == "completed" {
		utils.ErrorResponse(ctx, http.StatusConflict, errors.New("cannot undo after match is completed"), "match already completed")
		return
	}

	lastBall, err := dbHelper.GetLastBall(matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusNotFound, err, "no balls to undo")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to get last ball")
		return
	}

	// allow undoing balls from the CURRENT inning (shouldnt go back to prev inning)
	if lastBall.InningsID != liveMatchData.CurrentInningID {
		utils.ErrorResponse(ctx, http.StatusConflict, errors.New("cannot undo balls from a previous inning"), "previous inning locked")
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err = dbHelper.ArchiveBall(tx, lastBall.InningsID, lastBall.BallSequence)
		if err != nil {
			return err
		}

		err = dbHelper.ResetInningsStats(tx, lastBall.InningsID)
		if err != nil {
			return err
		}
		err = dbHelper.ResetScorecards(tx, lastBall.InningsID)
		if err != nil {
			return err
		}

		activeBalls, err := dbHelper.GetAllActiveBalls(tx, lastBall.InningsID)
		if err != nil {
			return err
		}

		var strikerID, bowlerID string
		var nonStrikerID *string

		if len(activeBalls) > 0 {
			strikerID, nonStrikerID, bowlerID, err = dbHelper.GetInitialInningsState(tx, lastBall.InningsID)
			if err != nil {
				return err
			}
		} else {
			strikerID = lastBall.StrikerID
			nonStrikerID = lastBall.NonStrikerID
			bowlerID = lastBall.BowlerID
		}

		err = dbHelper.ResetLiveMatchStats(tx, matchID, strikerID, nonStrikerID, bowlerID)
		if err != nil {
			return err
		}

		recalcMatchData := liveMatchData
		recalcMatchData.CurrentScore = 0
		recalcMatchData.Wickets = 0
		recalcMatchData.LegalBalls = 0
		recalcMatchData.IsCompleted = false
		recalcMatchData.IsFreeHit = false

		for _, ball := range activeBalls {
			ball.LegalBalls = recalcMatchData.LegalBalls

			var inningEnded, matchEnded bool
			err = ApplyBallStats(tx, ball, recalcMatchData, matchID, &inningEnded, &matchEnded)
			if err != nil {
				return err
			}

			totalRuns := ball.RunsBatter + ball.RunsExtra
			recalcMatchData.CurrentScore += totalRuns
			if ball.IsLegalDelivery {
				recalcMatchData.LegalBalls++
			}
			if ball.IsWicket {
				recalcMatchData.Wickets++
			}
			if !ball.IsLegalDelivery && ball.ExtraType != nil && *ball.ExtraType == "no_ball" {
				recalcMatchData.IsFreeHit = true
			} else if ball.IsLegalDelivery {
				recalcMatchData.IsFreeHit = false
			}
		}

		return nil
	})

	if txErr != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, txErr, "failed to undo ball")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ball undone successfully"})
}
