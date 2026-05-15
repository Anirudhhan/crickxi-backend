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
		userID, err = dbHelper.RegisterUser(tx, registerUserReq.Name, registerUserReq.Phone, hashedPassword)
		if err != nil {
			return err
		}

		playerID, err = dbHelper.RegisterPlayerStats(tx, userID)
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

	userDetails, err := dbHelper.GetLoginDetailsByPhone(loginUserReq.Phone)
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

	sessionID, err := dbHelper.CreateUserSession(userDetails.UserID)
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

func Logout(ctx *gin.Context) {
	sessionID := ctx.GetString("session_id")

	err := dbHelper.ArchiveUserSession(sessionID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to logout user")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user logged out successfully",
	})
}

func RefreshToken(ctx *gin.Context) {
	sessionID := ctx.GetHeader("sessionID")

	if sessionID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("missing session"), "unauthorized")
		return
	}

	userID, err := dbHelper.GetUserIDByActiveSession(sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err, "invalid session")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, err.Error())
		return
	}

	accessToken, err := utils.GenerateAccessToken(
		userID,
		sessionID,
	)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to generate token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}

func ResetPassword(ctx *gin.Context) {
	var userReq struct {
		PhoneNo     string `db:"phone_no" json:"phone" binding:"required"`
		NewPassword string `db:"password" json:"password" binding:"required,min=8"`
		OTP         string `json:"otp" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	if userReq.OTP != "8080" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("invalid otp"), "invalid otp")
		return
	}

	hashedPassword, err := utils.HashPassword(userReq.NewPassword)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, err := dbHelper.UpdateUserPassword(tx, userReq.PhoneNo, hashedPassword)
		if err != nil {
			return err
		}

		err = dbHelper.ArchiveUserSessions(tx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {

		if errors.Is(txErr, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusNotFound, txErr, "user not found")
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, txErr, "failed to update password")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "password updated successfully",
	})

}
