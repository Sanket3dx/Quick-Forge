package demonController

import (
	"fmt"
	"net/http"
	mongooDbModel "quick_forge/database/mongooDb/models"
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
	dbType, exists := ctx.Get("dbType")

	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbType not found in context")
		return
	}
	var allData []map[string]interface{}
	var err error
	switch dbType {
	case "mysql":
		allData, err = mysql_models.GetAllData(route.DBTableName)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "error fetching data from DB")
			return
		}
	case "mongo":
		allData, err = mongooDbModel.GetAllData(route.DBTableName)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "error fetching data from DB")
			return
		}
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

	dbType, exists := ctx.Get("dbType")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbType not found in context")
		return
	}

	var data map[string]interface{}
	var err error

	switch dbType {
	case "mysql":
		data, err = mysql_models.GetData(route, ctx.Param("arg"))
	case "mongo":
		data, err = mongooDbModel.GetData(route, ctx.Param("arg"))
	default:
		utils.HandleError(ctx, http.StatusInternalServerError, "Unsupported database type")
		return
	}

	if err != nil {
		fmt.Println(err)
		utils.HandleError(ctx, http.StatusInternalServerError, "Error fetching data from DB")
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

	dbType, exists := ctx.Get("dbType")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbType not found in context")
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
	var insertedID interface{}
	var err error

	switch dbType {
	case "mysql":
		insertedID, err = mysql_models.InsertData(route.DBTableName, data)
	case "mongo":
		insertedID, err = mongooDbModel.InsertData(route, data)
	default:
		utils.HandleError(ctx, http.StatusInternalServerError, "Unsupported database type")
		return
	}

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

	dbType, exists := ctx.Get("dbType")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbType not found in context")
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

	var err error

	switch dbType {
	case "mysql":
		err = mysql_models.UpdateData(route, data, ctx.Param("arg"))
	case "mongo":
		err = mongooDbModel.UpdateData(route, data, ctx.Param("arg"))
	default:
		utils.HandleError(ctx, http.StatusInternalServerError, "Unsupported database type")
		return
	}

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Record update failed")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Record updated successfully",
	})
}

func Deletehandler(ctx *gin.Context) {
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

	dbType, exists := ctx.Get("dbType")
	if !exists {
		utils.HandleError(ctx, http.StatusInternalServerError, "dbType not found in context")
		return
	}

	var err error

	switch dbType {
	case "mysql":
		err = mysql_models.DeleteData(route, ctx.Param("arg"))
	case "mongo":
		err = mongooDbModel.DeleteData(route, ctx.Param("arg"))
	default:
		utils.HandleError(ctx, http.StatusInternalServerError, "Unsupported database type")
		return
	}

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Record deletion failed")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Record deleted successfully",
	})
}
