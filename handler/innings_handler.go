package handler

import (
	"crickxi-backend/database/dbHelper"
	"crickxi-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOverDetails(ctx *gin.Context) {
	inningsID := ctx.Param("inningsID")

	overDetails, err := dbHelper.OverDetails(inningsID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, overDetails)
}
