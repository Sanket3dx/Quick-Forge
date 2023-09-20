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
		utils.HandleError(ctx, http.StatusInternalServerError, "Route info not found in context")
		return
	}
	route, ok := routeInfo.(utils.Route)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid routeInfo type")
		return
	}

	allData, err := mysql_models.GetAllData(route.DBTableName)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "error fetching data from DB")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"data":    allData,
		"records": len(allData),
	})
}

func Gethandler(ctx *gin.Context) {
	routeInfo, exists := ctx.Get("routeInfo")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Route info not found in context")
		return
	}
	route, ok := routeInfo.(utils.Route)

	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid routeInfo type")
		return
	}
	data, err := mysql_models.GetData(route, ctx.Param("arg"))
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "error fetching data from DB")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  data,
	})
}

func Posthandler(ctx *gin.Context) {
	// Retrieve route information from context
	routeInfo, exists := ctx.Get("routeInfo")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Route info not found in context")
		return
	}

	route, ok := routeInfo.(utils.Route)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid routeInfo type")
		return
	}

	validatedRequestBody, exists := ctx.Get("validatedRequestBody")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "Validated request body not found in context")
		return
	}

	data, ok := validatedRequestBody.(map[string]interface{})
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "Invalid validatedRequestBody type")
		return
	}

	// Insert data into the database
	insertedID, err := mysql_models.InsertData(route.DBTableName, data)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Record insertion failed")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":     false,
		"record_id": insertedID,
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
