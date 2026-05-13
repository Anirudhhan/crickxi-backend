package handler

import (
	"crickxi-backend/database"
	"crickxi-backend/database/dbhelper"
	"crickxi-backend/models"
	"crickxi-backend/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	DuplicatePGCode = "23505"
)

func RegisterUser(ctx *gin.Context) {
	var registerUserReq models.RegisterUser

	if err := ctx.ShouldBindJSON(&registerUserReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(registerUserReq.Password)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, err.Error())
		return
	}

	var userID string
	var playerID string
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, err = dbhelper.RegisterUser(tx, registerUserReq.Name, registerUserReq.Phone, hashedPassword)
		if err != nil {
			return err
		}

		playerID, err = dbhelper.RegisterPlayerStats(tx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		var pgErr *pq.Error

		if errors.As(txErr, &pgErr) {
			if pgErr.Code == DuplicatePGCode {
				utils.ErrorResponse(ctx, http.StatusConflict, txErr, "user already exists")
				return
			}
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, txErr, txErr.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "user registered successfully",
		"userID":   userID,
		"playerID": playerID,
	})
}
