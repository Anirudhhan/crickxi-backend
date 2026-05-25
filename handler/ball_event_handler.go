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

	matchID := ctx.Param("matchID")
	if matchID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("missing match id"), "missing match id")
		return
	}
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

	var delivery models.Delivery

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

	if delivery.IsWicket {

		if delivery.NextBatterID == nil || delivery.WicketPlayerID == nil || delivery.WicketType == nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid wicket data"), "invalid wicket data")
			return
		}
	}

	isStrikerOut, err := dbHelper.IsPlayerOut(delivery.InningsID, delivery.StrikerID)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	if isStrikerOut {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("striker is already out"), "striker is already out")
		return
	}

	isNonStrikerOut, err := dbHelper.IsPlayerOut(delivery.InningsID, delivery.NonStrikerID)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	if isNonStrikerOut {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("non striker is already out"), "non striker is already out")
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {

		err := dbHelper.CreateBallEvent(tx, delivery)
		if err != nil {
			return err
		}

		err = dbHelper.UpdateBattingScoreCard(tx, delivery)
		if err != nil {
			return err
		}

		err = dbHelper.UpdateBowlingScoreCard(tx, delivery)
		if err != nil {
			return err
		}

		err = dbHelper.UpdateInnings(tx, delivery)
		if err != nil {
			return err
		}

		err = dbHelper.UpdateLiveMatch(tx, delivery, matchID)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, txErr, "failed to add ball")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "ball added successfully",
	})
}
