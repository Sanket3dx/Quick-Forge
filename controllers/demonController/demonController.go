package demonController

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllhandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get request resived at getAllhandler",
	})
}
