package handler

import (
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOverDetails(ctx *gin.Context) {
	matchID := ctx.Param("matchID")
	inningsOrderStr := ctx.Param("inningsOrder")

	var inningsOrder int
	if inningsOrderStr == "1" {
		inningsOrder = 1
	} else if inningsOrderStr == "2" {
		inningsOrder = 2
	} else {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid innings order"), "invalid innings order")
		return
	}

	overDetails, err := dbHelper.OverDetails(matchID, inningsOrder)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, overDetails)
}
