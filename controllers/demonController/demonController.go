package demonController

import (
	"net/http"
	mysql_models "quick_forge/database/mysql/models"
	"quick_forge/utils"

	"github.com/gin-gonic/gin"
)

func GetAllhandler(ctx *gin.Context) {
	routeInfo, exists := ctx.Get("routeInfo")

	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Route info not found in context",
		})
		ctx.Abort()
		return
	}
	route, ok := routeInfo.(utils.Route)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid routeInfo type",
		})
		ctx.Abort()
		return
	}

	allData, err := mysql_models.GetAllData(route.DBTableName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error fetching data from DB",
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"data":    allData,
		"records": len(allData),
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
