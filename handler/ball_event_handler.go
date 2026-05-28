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
	}

}

func ProcessBallEventHelper(delivery models.Delivery, liveMatchData models.LiveMatchDetails, matchID string, inningEnded *bool, matchEnded *bool) error {
	err := database.Tx(func(tx *sqlx.Tx) error {

		err := dbHelper.CreateBallEvent(tx, delivery)
		if err != nil {
			return err
		}

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
				}
				if *delivery.ExtraType == "no_ball" {
					extraNoBall = delivery.RunsExtra
					noBall += 1
					nextFreeHit = true
				}
			}
		}

		//update batting scoreCard
		err = dbHelper.UpdateBatterStats(tx, delivery, legalBall, fours, sixes)
		if err != nil {
			return err
		}

		//Update dismissed player
		wicket := 0
		inningsWicket := 0
		if delivery.IsWicket && delivery.WicketPlayerID != nil {
			var dismissalBy *string
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
				}
			}

			err = dbHelper.UpdateDismissedBatter(tx, delivery, dismissalBy)
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
				totalMovementRuns := delivery.RunsBatter + delivery.RunsExtra
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
						// if strike rotated,
						// new batter becomes non striker
						if strikerID != originalStrikerID {
							nonStrikerID = delivery.NextBatterID
						} else {
							strikerID = *delivery.NextBatterID
						}
					} else if nonStrikerID != nil {
						// last batter standing
						strikerID = *nonStrikerID
						nonStrikerID = nil
					}
				}

				// original non striker out
				if originalNonStrikerID != nil && *delivery.WicketPlayerID == *originalNonStrikerID {
					if delivery.NextBatterID != nil {
						// if strike rotated,
						// new batter becomes striker
						if strikerID == *originalNonStrikerID {
							strikerID = *delivery.NextBatterID
						} else {
							nonStrikerID = delivery.NextBatterID
						}
					} else {
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
			if liveMatchData.PreviousInningsScore != nil {
				//handle score chased
				currentScore := liveMatchData.CurrentScore + totalRuns
				target := *liveMatchData.PreviousInningsScore
				if currentScore > target {
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
	})
	if err != nil {
		return err
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
