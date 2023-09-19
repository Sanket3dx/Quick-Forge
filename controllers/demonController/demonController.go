package demonController

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllhandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "GET ALL request recived at getAllhandler",
	})
}

func Gethandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "GET request recived at Gethandler",
	})
}

func Posthandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "POST request recived at Posthandler",
	})
}

func Puthandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "PUT request recived at Puthandler",
	})
}

func Deletehandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "DELETE request recived at Deletehandler",
	})
}
