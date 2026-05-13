package handler

import (
	"crickxi-backend/database"
	"crickxi-backend/database/dbhelper"
	"crickxi-backend/models"
	"crickxi-backend/utils"
	"database/sql"
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

func LoginUser(ctx *gin.Context) {
	var loginUserReq models.LoginUser

	if err := ctx.ShouldBindJSON(&loginUserReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	userDetails, err := dbhelper.GetLoginDetailsByPhone(loginUserReq.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusForbidden, err, "invalid credentials")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, err.Error())
		return
	}

	if err := utils.CheckPasswordHash(loginUserReq.Password, userDetails.HashPassword); err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err, "invalid credentials")
		return
	}

	sessionID, err := dbhelper.CreateUserSession(userDetails.UserID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, err.Error())
		return
	}

	accessToken, err := utils.GenerateAccessToken(userDetails.UserID, sessionID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "user logged in successfully",
		"sessionID":   sessionID,
		"accessToken": accessToken,
	})
}
